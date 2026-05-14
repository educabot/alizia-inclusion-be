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
		},
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, bcfg.Environment(cfg.Env)),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
