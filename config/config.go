package config

import (
	"math"
	"time"

	bcfg "github.com/educabot/team-ai-toolkit/config"
)

type Config struct {
	bcfg.BaseConfig
	AzureOpenAIKey      string
	AzureOpenAIEndpoint string
	AzureOpenAIModel    string

	// AuthPublicKey es la clave pública RS256 de auth-service (PEM). Env canónica
	// AUTH_PUBLIC_KEY (unificada con alizia-be/seguridad); acepta el alias legacy
	// JWT_PUBLIC_KEY como fallback durante la migración.
	AuthPublicKey string

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

	// EmbeddingDim is the dimension of the pedagogical-content embedding vectors.
	// Single source of truth for the (Futuro) vector search: the embedding column
	// stays inert in the MVP (keyword/full-text first), and the vector index needs
	// a fixed dimension, so the real value must be confirmed against the Azure model
	// before enabling it. Default 1536 (OpenAI/Azure text-embedding-3-small).
	EmbeddingDim int
}

func Load() *Config {
	base := bcfg.LoadBase()
	return &Config{
		BaseConfig:          base,
		AzureOpenAIKey:      bcfg.EnvOr("AZURE_OPENAI_API_KEY", ""),
		AzureOpenAIEndpoint: bcfg.EnvOr("AZURE_OPENAI_ENDPOINT", ""),
		AzureOpenAIModel:    bcfg.EnvOr("AZURE_OPENAI_MODEL", "gpt-4o-mini"),

		AuthPublicKey: bcfg.EnvOr("AUTH_PUBLIC_KEY", bcfg.EnvOr("JWT_PUBLIC_KEY", "")),

		DBMaxOpenConns:    boundedUintToInt(bcfg.GetEnvUint("DB_MAX_OPEN_CONNS", "25")),
		DBMaxIdleConns:    boundedUintToInt(bcfg.GetEnvUint("DB_MAX_IDLE_CONNS", "10")),
		DBConnMaxLifetime: time.Duration(boundedUintToInt(bcfg.GetEnvUint("DB_CONN_MAX_LIFETIME_MIN", "30"))) * time.Minute,
		DBConnMaxIdleTime: time.Duration(boundedUintToInt(bcfg.GetEnvUint("DB_CONN_MAX_IDLE_TIME_MIN", "5"))) * time.Minute,

		AIRateLimitPerHour:        boundedUintToInt(bcfg.GetEnvUint("AI_RATE_LIMIT_PER_HOUR", "0")),
		AICircuitFailureThreshold: boundedUintToInt(bcfg.GetEnvUint("AI_CIRCUIT_FAILURE_THRESHOLD", "5")),
		AICircuitCooldown:         time.Duration(boundedUintToInt(bcfg.GetEnvUint("AI_CIRCUIT_COOLDOWN_SEC", "30"))) * time.Second,

		AIAgenticEnabled: bcfg.EnvOr("AI_AGENTIC_ENABLED", "false") == "true",

		EmbeddingDim: boundedUintToInt(bcfg.GetEnvUint("EMBEDDING_DIM", "1536")),
	}
}

func boundedUintToInt(v uint) int {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	return int(v)
}
