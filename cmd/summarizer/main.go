// Command summarizer es el job batch (cron) que genera los resúmenes de las
// conversaciones cerradas, para darle memoria entre clases al chat. Corre un
// lote y termina; en Railway se deploya como servicio aparte con cronSchedule.
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/app/database"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	air "github.com/educabot/alizia-inclusion-be/src/repositories/ai"
	inclusionr "github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
)

// runTimeout acota la corrida completa del lote (red + LLM por conversación).
const runTimeout = 5 * time.Minute

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("summarizer: database connection failed", "error", err)
		os.Exit(1)
	}
	closeDB := func() {
		if sqlDB, dErr := db.DB(); dErr == nil {
			_ = sqlDB.Close()
		}
	}

	uc := inclusionuc.NewSummarizeConversations(
		inclusionr.NewConversationRepo(db),
		inclusionr.NewConversationSummaryRepo(db),
		buildAIClient(cfg),
		inclusionr.NewAIUsageRepo(db),
		cfg.SummaryIdleMinutes,
		cfg.SummaryBatchLimit,
	)

	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)

	res, err := uc.Execute(ctx)
	cancel()
	if err != nil {
		slog.Error("summarizer: run failed", "error", err)
		closeDB()
		os.Exit(1)
	}
	slog.Info("summarizer: run complete", "processed", res.Processed, "failed", res.Failed)
	closeDB()
}

// buildAIClient replica el armado del cliente de IA del server web (cmd): Azure
// con circuit breaker si hay credenciales, stub si no.
func buildAIClient(cfg *config.Config) providers.AIClient {
	if cfg.AzureOpenAIKey != "" && cfg.AzureOpenAIEndpoint != "" && cfg.AzureOpenAIKey != "your-azure-openai-key" {
		return air.NewCircuitBreaker(
			air.NewAzureClient(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIKey, cfg.AzureOpenAIModel, cfg.AzureOpenAIAPIVersion),
			cfg.AICircuitFailureThreshold,
			cfg.AICircuitCooldown,
			nil,
		)
	}
	slog.Warn("summarizer: using stub AI client (no Azure credentials)")
	return air.NewStubClient()
}
