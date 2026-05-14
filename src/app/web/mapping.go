package web

import (
	"github.com/gin-gonic/gin"
	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
)

func ConfigureMappings(engine *gin.Engine, h *entrypoints.WebHandlerContainer, _ *config.Config) {
	api := engine.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(h.AuthMiddleware))
	api.Use(webgin.AdaptMiddleware(h.TenantMiddleware))

	// Catalog: ramps & devices (any authenticated user)
	api.GET("/ramps", webgin.Adapt(h.Catalog.HandleListRamps))
	api.GET("/ramps/:id", webgin.Adapt(h.Catalog.HandleGetRamp))
	api.GET("/devices", webgin.Adapt(h.Catalog.HandleListDevices))
	api.GET("/devices/:id", webgin.Adapt(h.Catalog.HandleGetDevice))

	// Students & profiles
	api.GET("/students/:student_id/profile", webgin.Adapt(h.Inclusion.HandleGetStudentProfile))
	api.PUT("/students/:student_id/profile", webgin.Adapt(h.Inclusion.HandleUpsertStudentProfile))
	api.GET("/classrooms/:classroom_id/students", webgin.Adapt(h.Inclusion.HandleListClassroomStudents))

	// AI endpoints
	api.POST("/inclusion/recommend", webgin.Adapt(h.Inclusion.HandleRecommendDevice))
	api.POST("/inclusion/assist", webgin.Adapt(h.Inclusion.HandleAssistClassroom))
}
