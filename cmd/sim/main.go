// Command sim reproduce, contra la Postgres local + el AIClient real de Azure, la
// conversación del caso de María trabajado con Mercedes, para evaluar a ojo el
// comportamiento del prompt de Alizia. NO es parte del binario de producción.
//
// Uso (DATABASE_URL DEBE apuntar a la DB LOCAL, nunca a Railway):
//
//	DATABASE_URL="postgres://postgres:postgres@localhost:5481/alizia_inclusion?sslmode=disable" \
//	  ENV=local AI_AGENTIC_ENABLED=true go run ./cmd/sim
//
// Las claves AZURE_OPENAI_* se toman del entorno (las exporta el script de corrida).
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/app/database"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	air "github.com/educabot/alizia-inclusion-be/src/repositories/ai"
	catalogr "github.com/educabot/alizia-inclusion-be/src/repositories/catalog"
	inclusionr "github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
	mgmtr "github.com/educabot/alizia-inclusion-be/src/repositories/management"
)

func main() {
	cfg := config.Load()

	// Salvaguarda: nunca correr la simulación contra una DB remota.
	if dsn := os.Getenv("DATABASE_URL"); dsn == "" ||
		(!strings.Contains(dsn, "localhost") && !strings.Contains(dsn, "127.0.0.1")) {
		fmt.Println("ABORTANDO: DATABASE_URL debe apuntar a la DB local (localhost). DSN:", dsn)
		os.Exit(1)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		fmt.Println("db connect:", err)
		os.Exit(1)
	}

	var aiClient providers.AIClient
	if cfg.AzureOpenAIKey != "" && cfg.AzureOpenAIEndpoint != "" {
		aiClient = air.NewAzureClient(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIKey, cfg.AzureOpenAIModel, cfg.AzureOpenAIAPIVersion)
	} else {
		fmt.Println("ABORTANDO: faltan AZURE_OPENAI_API_KEY / AZURE_OPENAI_ENDPOINT")
		os.Exit(1)
	}

	deps := inclusionuc.AssistClassroomDeps{
		AI:            aiClient,
		Students:      inclusionr.NewStudentRepo(db),
		Profiles:      inclusionr.NewStudentProfileRepo(db),
		Classrooms:    mgmtr.NewClassroomRepo(db),
		Devices:       catalogr.NewDeviceRepo(db),
		Conversations: nil, // sin persistencia: la historia la pasamos a mano
		Summaries:     inclusionr.NewConversationSummaryRepo(db),
		Adaptations:   inclusionr.NewAdaptationRepo(db),
		Content:       inclusionr.NewPedagogicalContentRepo(db),
		Embedder:      air.NewAzureEmbedder(cfg.AzureEmbeddingEndpoint, cfg.AzureEmbeddingAPIKey, cfg.AzureEmbeddingDeployment, cfg.AzureEmbeddingAPIVersion, cfg.EmbeddingDim),
		RAG:           inclusionr.NewRAGSearchRepo(db),
		Usage:         inclusionr.NewAIUsageRepo(db),
		Agentic:       cfg.AIAgenticEnabled,
	}
	assist := inclusionuc.NewAssistClassroom(deps)

	orgID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	var userID int64 = 7 // docente1@demo.edu

	// Mensajes del profe, tal cual el caso trabajado con Mercedes (1ra reunión: apertura
	// + edad/momento/tipo; 2da: profundizar organización -> propuesta -> cierre).
	turns := []string{
		"María es muy inquieta no deja de moverse durante las tareas y consignas, le propuse de todo y no logro ayudarla a que se concentre y realice sus trabajos con atención.",
		"8 años. Todas. Activa.",
		"Dale, sigamos afinando.",
		"Organización",
		"Cuando llega al aula: al llegar no se sienta, no saca los materiales, no sabe dónde dejar la mochila y le tengo que decir lo que tiene que hacer.",
		"Ambas",
		"Bueno, voy a empezar con el organizador, después me ayudás a elegir el alumno. ¡Gracias!",
	}

	ctx := context.Background()
	var history []providers.ChatMessage

	for i, msg := range turns {
		fmt.Printf("\n\n========================= TURNO %d =========================\n", i+1)
		fmt.Printf(">> PROFE: %s\n\n", msg)

		resp, err := assist.Execute(ctx, inclusionuc.AssistClassroomRequest{
			OrgID:   orgID,
			UserID:  userID,
			Message: msg,
			Mode:    "assist",
			History: history,
		})
		if err != nil {
			fmt.Println("!! ERROR:", err)
			os.Exit(1)
		}

		fmt.Printf("<< ALIZIA:\n%s\n", resp.Response)
		if resp.Adaptation != nil {
			fmt.Printf("\n[adaptación estructurada: %q tipo=%s pasos=%d]\n", resp.Adaptation.Title, resp.Adaptation.Type, len(resp.Adaptation.Steps))
		}
		if len(resp.ReferencedContent) > 0 {
			fmt.Printf("[referenced_content: %d (NO debería haber cita)]\n", len(resp.ReferencedContent))
		}

		history = append(history,
			providers.ChatMessage{Role: "user", Content: msg},
			providers.ChatMessage{Role: "assistant", Content: resp.Response},
		)
	}
	fmt.Println("\n\n=== fin de la simulación ===")
}
