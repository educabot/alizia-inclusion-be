package config

import (
	"math"
	"time"

	bcfg "github.com/educabot/team-ai-toolkit/config"
)

type Config struct {
	bcfg.BaseConfig
	AzureOpenAIKey        string
	AzureOpenAIEndpoint   string
	AzureOpenAIModel      string
	AzureOpenAIAPIVersion string

	// Azure embeddings (búsqueda híbrida RAG). Recurso Azure potencialmente distinto
	// del de chat, por eso endpoint/key/deployment propios. EmbeddingDim ya existe.
	AzureEmbeddingEndpoint   string
	AzureEmbeddingAPIKey     string
	AzureEmbeddingDeployment string
	AzureEmbeddingAPIVersion string

	JWTPublicKey string

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	// AIRateLimitPerHour caps AI requests per organization per hour. 0 = unlimited.
	AIRateLimitPerHour int
	// AICircuitFailureThreshold is the number of consecutive AI client failures
	// that trip the circuit breaker open.
	AICircuitFailureThreshold int
	// AICircuitCooldown is how long the circuit stays open before allowing a trial call.
	AICircuitCooldown time.Duration

	// AIAgenticEnabled turns on the tool-calling (function calling) loop for the
	// assist endpoint. Off by default until validated against the live AI provider.
	AIAgenticEnabled bool

	// ChatTraceVerbose enables full-text logging of the chat pipeline (prompts,
	// responses, student names). Default true; set false in prod to avoid PII in logs
	// (metadata —ids, counts, scores, tokens— is always logged regardless).
	ChatTraceVerbose bool

	// EmbeddingDim is the dimension of the pedagogical-content embedding vectors.
	// Single source of truth for the (Futuro) vector search: the embedding column
	// stays inert in the MVP (keyword/full-text first), and the vector index needs
	// a fixed dimension, so the real value must be confirmed against the Azure model
	// before enabling it. Default 1536 (OpenAI/Azure text-embedding-3-small).
	EmbeddingDim int

	// SummaryIdleMinutes is the minimum age (minutes) of a conversation's last
	// message before the summarizer cron treats it as "closed" and resumes it. Default 20.
	SummaryIdleMinutes int
	// SummaryBatchLimit caps how many conversations the summarizer cron processes
	// per run. Default 50.
	SummaryBatchLimit int
}

func Load() *Config {
	// Desde team-ai-toolkit v1.10.0, LoadBase() ya NO exige JWT_SECRET (lo lee opcional):
	// este servicio no firma tokens HS256, la auth es RS256 contra el auth-service.
	base := bcfg.LoadBase()
	return &Config{
		BaseConfig:          base,
		AzureOpenAIKey:      bcfg.EnvOr("AZURE_OPENAI_API_KEY", ""),
		AzureOpenAIEndpoint: bcfg.EnvOr("AZURE_OPENAI_ENDPOINT", ""),
		AzureOpenAIModel:    bcfg.EnvOr("AZURE_OPENAI_MODEL", "gpt-5.4"),
		// Reasoning (gpt-5.x) requiere api-version 2024-12-01-preview o posterior.
		AzureOpenAIAPIVersion: bcfg.EnvOr("AZURE_OPENAI_API_VERSION", "2024-12-01-preview"),

		AzureEmbeddingEndpoint:   bcfg.EnvOr("AZURE_OPENAI_EMBEDDING_ENDPOINT", ""),
		AzureEmbeddingAPIKey:     bcfg.EnvOr("AZURE_OPENAI_EMBEDDING_API_KEY", ""),
		AzureEmbeddingDeployment: bcfg.EnvOr("AZURE_OPENAI_EMBEDDING_DEPLOYMENT", "text-embedding-3-small"),
		AzureEmbeddingAPIVersion: bcfg.EnvOr("AZURE_OPENAI_EMBEDDING_API_VERSION", "2024-02-01"),

		// AUTH_PUBLIC_KEY es el nombre nuevo (auth-service actualizado); JWT_PUBLIC_KEY queda como fallback.
		JWTPublicKey: bcfg.EnvOr("AUTH_PUBLIC_KEY", bcfg.EnvOr("JWT_PUBLIC_KEY", "")),

		DBMaxOpenConns:    boundedUintToInt(bcfg.GetEnvUint("DB_MAX_OPEN_CONNS", "25")),
		DBMaxIdleConns:    boundedUintToInt(bcfg.GetEnvUint("DB_MAX_IDLE_CONNS", "10")),
		DBConnMaxLifetime: time.Duration(boundedUintToInt(bcfg.GetEnvUint("DB_CONN_MAX_LIFETIME_MIN", "30"))) * time.Minute,
		DBConnMaxIdleTime: time.Duration(boundedUintToInt(bcfg.GetEnvUint("DB_CONN_MAX_IDLE_TIME_MIN", "5"))) * time.Minute,

		AIRateLimitPerHour:        boundedUintToInt(bcfg.GetEnvUint("AI_RATE_LIMIT_PER_HOUR", "0")),
		AICircuitFailureThreshold: boundedUintToInt(bcfg.GetEnvUint("AI_CIRCUIT_FAILURE_THRESHOLD", "5")),
		AICircuitCooldown:         time.Duration(boundedUintToInt(bcfg.GetEnvUint("AI_CIRCUIT_COOLDOWN_SEC", "30"))) * time.Second,

		AIAgenticEnabled: bcfg.EnvOr("AI_AGENTIC_ENABLED", "false") == "true",

		ChatTraceVerbose: bcfg.EnvOr("CHAT_TRACE_VERBOSE", "true") == "true",

		EmbeddingDim: boundedUintToInt(bcfg.GetEnvUint("EMBEDDING_DIM", "1536")),

		SummaryIdleMinutes: boundedUintToInt(bcfg.GetEnvUint("SUMMARY_IDLE_MINUTES", "20")),
		SummaryBatchLimit:  boundedUintToInt(bcfg.GetEnvUint("SUMMARY_BATCH_LIMIT", "50")),
	}
}

func boundedUintToInt(v uint) int {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	return int(v)
}
