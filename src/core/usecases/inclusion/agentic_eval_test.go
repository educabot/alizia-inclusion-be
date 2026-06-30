//go:build eval

// Eval de comportamiento del modelo real (Azure / gpt-5.4): NO mockea el AI.
//
// A diferencia de agentic_test.go (que mockea ChatWithTools y testea el loop), este
// suite verifica que, ante inputs que SÍ O SÍ deberían disparar una búsqueda, el
// modelo efectivamente decide llamar la tool correcta. La verdad se lee de
// AssistClassroomResponse.SourcesUsed.Tools (el mismo trace que el log chat.sources_used).
//
// No corre en el CI bloqueante: requiere el tag `eval` y credenciales de Azure, y es
// no-determinístico por naturaleza (por eso pass-rate y reintentos). Pensado para
// on-demand / nightly.
//
//	make eval                 # corre todo el suite (carga .env)
//	go test -tags eval -run TestAgenticEval ./src/core/usecases/inclusion/ -v
//
// Variables de entorno:
//
//	AZURE_OPENAI_API_KEY / _ENDPOINT / _MODEL / _API_VERSION  (se cargan del .env si faltan)
//	EVAL_RUNS       cuántas veces se corre cada caso (default 3)
//	EVAL_THRESHOLD  pass-rate mínimo para aprobar un caso, 0..1 (default 0.5)
package inclusion

import (
	"bufio"
	"context"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	airepo "github.com/educabot/alizia-inclusion-be/src/repositories/ai"
)

// ---------------------------------------------------------------------------
// Stubs de contexto. Solo necesitamos poblar el system prompt (alumnos + valija);
// los providers de las tools quedan nil: el dispatcher es nil-safe y, aunque la tool
// devuelva "unavailable", el trace IGUAL registra que el modelo la invocó, que es lo
// único que el eval afirma. Ver runAgenticChat: trace se acumula antes del error.
// ---------------------------------------------------------------------------

type evalStudents struct{ list []entities.Student }

func (s *evalStudents) List(_ context.Context, _ uuid.UUID) ([]entities.Student, error) {
	return s.list, nil
}
func (s *evalStudents) GetStudent(_ context.Context, _ uuid.UUID, _ int64) (*entities.Student, error) {
	return nil, nil
}
func (s *evalStudents) ListByClassroom(_ context.Context, _ uuid.UUID, _ int64) ([]entities.Student, error) {
	return s.list, nil
}
func (s *evalStudents) Create(_ context.Context, _ *entities.Student) error  { return nil }
func (s *evalStudents) Update(_ context.Context, _ *entities.Student) error  { return nil }
func (s *evalStudents) Delete(_ context.Context, _ uuid.UUID, _ int64) error { return nil }

type evalDevices struct{ list []entities.Device }

func (d *evalDevices) ListDevices(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
	return d.list, nil
}
func (d *evalDevices) GetDevice(_ context.Context, _ uuid.UUID, _ int64) (*entities.Device, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Casos: input que debería forzar tool calls + qué tool(s) esperamos ver.
// ---------------------------------------------------------------------------

func ptrInt64(v int64) *int64 { return &v }

type evalCase struct {
	name      string
	students  []entities.Student
	devices   []entities.Device
	studentID *int64
	message   string
	// wantAnyTool: el caso "llamó tool" si Tools contiene AL MENOS una de estas.
	wantAnyTool []string
}

func evalCases() []evalCase {
	maria := entities.Student{ID: 7, Name: "María", ClassroomID: 1}
	devices := []entities.Device{
		{ID: 1, Name: "Cojín dinámico", RampID: 1},
		{ID: 31, Name: "Banda elástica para sillas", RampID: 1},
	}
	return []evalCase{
		{
			name:        "historial_alumno",
			students:    []entities.Student{maria},
			devices:     devices,
			studentID:   ptrInt64(7),
			message:     "¿Qué veníamos trabajando con María en las sesiones anteriores? Quiero retomar desde ahí.",
			wantAnyTool: []string{"get_student_history"},
		},
		{
			name:        "adaptaciones_previas",
			students:    []entities.Student{maria},
			devices:     devices,
			studentID:   ptrInt64(7),
			message:     "¿Qué adaptaciones le habíamos armado antes a María? No me quiero repetir.",
			wantAnyTool: []string{"get_past_adaptations"},
		},
		{
			name:        "busqueda_contenido_rag",
			students:    []entities.Student{maria},
			devices:     devices,
			message:     "Necesito actividades concretas de matemática para un nene de 8 con discalculia. ¿Qué recursos pedagógicos tenés?",
			wantAnyTool: []string{"search_content", "search_content_hibrido", "get_content"},
		},
		{
			// Caso "propone sin haber buscado": pedido de recomendación basado en una
			// barrera observable. Aunque el modelo intente cerrar con un paso a paso sin
			// buscar, la red de seguridad (requireSearchBeforeProposal en AssistClassroom)
			// lo obliga a llamar search_content_hibrido antes de proponer.
			name:        "recomendacion_fundamentada_rag",
			students:    []entities.Student{maria},
			devices:     devices,
			message:     "Tengo una alumna de 8 que se mueve todo el tiempo y no logra sostener la tarea. ¿Cómo la ayudo a engancharse con la actividad? Dame algo concreto.",
			wantAnyTool: []string{"search_content", "search_content_hibrido", "get_content"},
		},
	}
}

func TestAgenticEval_ToolCalls(t *testing.T) {
	ai := newEvalAIClient(t)
	runs := envInt("EVAL_RUNS", 3)
	threshold := envFloat("EVAL_THRESHOLD", 0.5)
	orgID := uuid.New()

	for _, tc := range evalCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			deps := AssistClassroomDeps{
				AI:       ai,
				Students: &evalStudents{list: tc.students},
				Devices:  &evalDevices{list: tc.devices},
				Agentic:  true,
				// Resto de providers nil a propósito (ver comentario de los stubs).
			}
			uc := NewAssistClassroom(deps)

			hits := 0
			for i := 0; i < runs; i++ {
				resp, err := uc.Execute(context.Background(), AssistClassroomRequest{
					OrgID:     orgID,
					UserID:    1,
					Message:   tc.message,
					StudentID: tc.studentID,
				})
				if err != nil {
					t.Fatalf("run %d: Execute falló: %v", i+1, err)
				}
				called := tc.wantAnyTool
				ok := containsAny(resp.SourcesUsed.Tools, called)
				if ok {
					hits++
				}
				t.Logf("run %d/%d: tools=%v esperaba_una_de=%v -> %v",
					i+1, runs, resp.SourcesUsed.Tools, called, passLabel(ok))
			}

			rate := float64(hits) / float64(runs)
			t.Logf("pass-rate %s: %d/%d = %.0f%% (umbral %.0f%%)",
				tc.name, hits, runs, rate*100, threshold*100)
			if rate < threshold {
				t.Errorf("%s: el modelo llamó la tool esperada en %d/%d corridas (%.0f%% < umbral %.0f%%). Tools esperadas: %v",
					tc.name, hits, runs, rate*100, threshold*100, tc.wantAnyTool)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func containsAny(got, want []string) bool {
	set := map[string]bool{}
	for _, g := range got {
		set[g] = true
	}
	for _, w := range want {
		if set[w] {
			return true
		}
	}
	return false
}

func passLabel(ok bool) string {
	if ok {
		return "OK"
	}
	return "MISS"
}

// newEvalAIClient construye el cliente Azure real desde el entorno. Carga el .env del
// repo si las vars no están seteadas, y saltea el test (no falla) si no hay API key.
func newEvalAIClient(t *testing.T) providers.AIClient {
	t.Helper()
	loadDotEnvIfMissing()
	key := os.Getenv("AZURE_OPENAI_API_KEY")
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	if key == "" || endpoint == "" {
		t.Skip("AZURE_OPENAI_API_KEY/ENDPOINT no seteadas (ni en el .env del repo); salteando eval")
	}
	model := envOr("AZURE_OPENAI_MODEL", "gpt-5.4")
	apiVersion := envOr("AZURE_OPENAI_API_VERSION", "2024-12-01-preview")
	return airepo.NewAzureClient(endpoint, key, model, apiVersion)
}

// loadDotEnvIfMissing busca un .env subiendo desde el cwd y carga las claves que aún
// no estén en el entorno. Es best-effort: si no hay .env, no hace nada.
func loadDotEnvIfMissing() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for i := 0; i < 6; i++ {
		path := filepath.Join(dir, ".env")
		if f, err := os.Open(path); err == nil {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := strings.TrimSpace(sc.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				k, v, ok := strings.Cut(line, "=")
				if !ok {
					continue
				}
				k = strings.TrimSpace(k)
				v = strings.Trim(strings.TrimSpace(v), `"'`)
				if _, exists := os.LookupEnv(k); !exists {
					_ = os.Setenv(k, v)
				}
			}
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 && f <= 1 {
			return f
		}
	}
	return math.Min(math.Max(def, 0), 1)
}
