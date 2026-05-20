package main

import (
	bcfg "github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

const jwtIssuer = "alizia-inclusion"

func NewHandlers(uc *UseCases, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret, jwtIssuer)

	return &entrypoints.WebHandlerContainer{
		Auth: &entrypoints.AuthContainer{
			Toker:   toker,
			LoginUC: uc.Login,
			GetMe:   uc.GetMe,
		},
		Catalog: &entrypoints.CatalogContainer{
			ListRamps:   uc.ListRamps,
			GetRamp:     uc.GetRamp,
			ListDevices: uc.ListDevices,
			GetDevice:   uc.GetDevice,
		},
		Inclusion: &entrypoints.InclusionContainer{
			GetStudentProfile:     uc.GetStudentProfile,
			UpsertStudentProfile:  uc.UpsertStudentProfile,
			ListClassroomStudents: uc.ListClassroomStudents,
			RecommendDevice:       uc.RecommendDevice,
			AssistClassroom:       uc.AssistClassroom,
			ListStudents:          uc.ListStudents,
			CreateStudent:         uc.CreateStudent,
			UpdateStudent:         uc.UpdateStudent,
			DeleteStudent:         uc.DeleteStudent,
			ListAdaptations:         uc.ListAdaptations,
			GetAdaptation:           uc.GetAdaptation,
			CreateAdaptation:        uc.CreateAdaptation,
			UpdateAdaptation:        uc.UpdateAdaptation,
			DeleteAdaptation:        uc.DeleteAdaptation,
			ListAdaptationResources: uc.ListAdaptationResources,
			GetChatHistory:          uc.GetChatHistory,
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
		},
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, bcfg.Environment(cfg.Env)),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
