package main

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	air "github.com/educabot/alizia-inclusion-be/src/repositories/ai"
	authr "github.com/educabot/alizia-inclusion-be/src/repositories/auth"
	catalogr "github.com/educabot/alizia-inclusion-be/src/repositories/catalog"
	inclusionr "github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
	mgmtr "github.com/educabot/alizia-inclusion-be/src/repositories/management"
)

type Repositories struct {
	Ramps               providers.RampProvider
	Devices             providers.DeviceProvider
	Students            providers.StudentProvider
	StudentProfiles     providers.StudentProfileProvider
	AI                  providers.AIClient
	Users               providers.UserProvider
	Classrooms          providers.ClassroomProvider
	Adaptations         providers.AdaptationProvider
	AdaptationResources providers.AdaptationResourceProvider
	Conversations       providers.ConversationProvider
	AIUsage             providers.AIUsageProvider
}

func NewRepositories(db *gorm.DB, cfg *config.Config) *Repositories {
	var aiClient providers.AIClient
	if cfg.AzureOpenAIKey != "" && cfg.AzureOpenAIEndpoint != "" && cfg.AzureOpenAIKey != "your-azure-openai-key" {
		aiClient = air.NewAzureClient(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIKey, cfg.AzureOpenAIModel)
		slog.Info("using Azure OpenAI client", "endpoint", cfg.AzureOpenAIEndpoint, "model", cfg.AzureOpenAIModel)
	} else {
		aiClient = air.NewStubClient()
		slog.Warn("using stub AI client, set AZURE_OPENAI_API_KEY and AZURE_OPENAI_ENDPOINT for real AI")
	}

	return &Repositories{
		Ramps:               catalogr.NewRampRepo(db),
		Devices:             catalogr.NewDeviceRepo(db),
		Students:            inclusionr.NewStudentRepo(db),
		StudentProfiles:     inclusionr.NewStudentProfileRepo(db),
		AI:                  aiClient,
		Users:               authr.NewUserRepo(db),
		Classrooms:          mgmtr.NewClassroomRepo(db),
		Adaptations:         inclusionr.NewAdaptationRepo(db),
		AdaptationResources: inclusionr.NewAdaptationResourceRepo(db),
		Conversations:       inclusionr.NewConversationRepo(db),
		AIUsage:             inclusionr.NewAIUsageRepo(db),
	}
}
