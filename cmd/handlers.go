package main

import (
	bcfg "github.com/educabot/team-ai-toolkit/config"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

func NewHandlers(uc *UseCases, cfg *config.Config) *entrypoints.WebHandlerContainer {
	return &entrypoints.WebHandlerContainer{
		Auth: &entrypoints.AuthContainer{
			GetMe: uc.GetMe,
		},
		Catalog: &entrypoints.CatalogContainer{
			ListRamps:   uc.ListRamps,
			GetRamp:     uc.GetRamp,
			ListDevices: uc.ListDevices,
			GetDevice:   uc.GetDevice,
		},
		Inclusion: &entrypoints.InclusionContainer{
			GetStudentProfile:        uc.GetStudentProfile,
			UpsertStudentProfile:     uc.UpsertStudentProfile,
			ListClassroomStudents:    uc.ListClassroomStudents,
			RecommendDevice:          uc.RecommendDevice,
			AssistClassroom:          uc.AssistClassroom,
			OpenSession:              uc.OpenSession,
			BuildPromptContext:       uc.BuildPromptContext,
			SearchPedagogicalContent: uc.SearchPedagogicalContent,
			ListStudents:             uc.ListStudents,
			CreateStudent:            uc.CreateStudent,
			UpdateStudent:            uc.UpdateStudent,
			DeleteStudent:            uc.DeleteStudent,
			ListAdaptations:          uc.ListAdaptations,
			GetAdaptation:            uc.GetAdaptation,
			CreateAdaptation:         uc.CreateAdaptation,
			UpdateAdaptation:         uc.UpdateAdaptation,
			DeleteAdaptation:         uc.DeleteAdaptation,
			ListAdaptationResources:  uc.ListAdaptationResources,
			ExportAdaptation:         uc.ExportAdaptation,
			GetChatHistory:           uc.GetChatHistory,
			DeleteConversation:       uc.DeleteConversation,
			RenameConversation:       uc.RenameConversation,
		},
		Management: &entrypoints.ManagementContainer{
			ListClassrooms:  uc.ListClassrooms,
			GetClassroom:    uc.GetClassroom,
			CreateClassroom: uc.CreateClassroom,
			UpdateClassroom: uc.UpdateClassroom,
			DeleteClassroom: uc.DeleteClassroom,
			ListTeachers:    uc.ListTeachers,
		},
		Dashboard: &entrypoints.DashboardContainer{
			GetMetrics: uc.GetMetrics,
			GetAIUsage: uc.GetAIUsage,
		},
		AuthMiddleware:   middleware.RS256AuthMiddleware(cfg.JWTPublicKey, bcfg.Environment(cfg.Env)),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
