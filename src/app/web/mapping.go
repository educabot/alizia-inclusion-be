package web

import (
	"net/http"
	"strconv"

	webgin "github.com/educabot/team-ai-toolkit/web/gin"
	"github.com/gin-gonic/gin"

	"github.com/educabot/alizia-inclusion-be/config"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func ConfigureMappings(engine *gin.Engine, h *entrypoints.WebHandlerContainer, _ *config.Config) {
	api := engine.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(h.AuthMiddleware))
	api.Use(webgin.AdaptMiddleware(h.TenantMiddleware))

	// Auth (authenticated)
	api.GET("/auth/me", webgin.Adapt(h.Auth.HandleGetMe))

	// Management: classrooms
	api.GET("/classrooms", webgin.Adapt(h.Management.HandleListClassrooms))
	api.POST("/classrooms", webgin.Adapt(h.Management.HandleCreateClassroom))

	classroomByID := api.Group("/classrooms/:id")
	classroomByID.GET("", webgin.Adapt(h.Management.HandleGetClassroom))
	classroomByID.PUT("", webgin.Adapt(h.Management.HandleUpdateClassroom))
	classroomByID.DELETE("", webgin.Adapt(h.Management.HandleDeleteClassroom))
	classroomByID.GET("/students", webgin.Adapt(h.Inclusion.HandleListClassroomStudents))

	// Management: teachers
	api.GET("/teachers", webgin.Adapt(h.Management.HandleListTeachers))

	// Students CRUD
	api.GET("/students", webgin.Adapt(h.Inclusion.HandleListStudents))
	api.POST("/students", webgin.Adapt(h.Inclusion.HandleCreateStudent))

	studentByID := api.Group("/students/:id")
	studentByID.GET("", webgin.Adapt(h.Inclusion.HandleGetStudent))
	studentByID.PUT("", webgin.Adapt(h.Inclusion.HandleUpdateStudent))
	studentByID.DELETE("", webgin.Adapt(h.Inclusion.HandleDeleteStudent))
	studentByID.GET("/profile", webgin.Adapt(h.Inclusion.HandleGetStudentProfile))
	studentByID.PUT("/profile", webgin.Adapt(h.Inclusion.HandleUpsertStudentProfile))

	// Catalog: ramps & devices
	api.GET("/ramps", webgin.Adapt(h.Catalog.HandleListRamps))
	api.GET("/ramps/:id", webgin.Adapt(h.Catalog.HandleGetRamp))
	api.GET("/devices", webgin.Adapt(h.Catalog.HandleListDevices))
	api.GET("/devices/:id", webgin.Adapt(h.Catalog.HandleGetDevice))

	// Adaptations
	api.GET("/adaptations", webgin.Adapt(h.Inclusion.HandleListAdaptations))
	api.POST("/adaptations", webgin.Adapt(h.Inclusion.HandleCreateAdaptation))
	api.GET("/adaptations/:id", webgin.Adapt(h.Inclusion.HandleGetAdaptation))
	api.PUT("/adaptations/:id", webgin.Adapt(h.Inclusion.HandleUpdateAdaptation))
	api.DELETE("/adaptations/:id", webgin.Adapt(h.Inclusion.HandleDeleteAdaptation))
	api.GET("/adaptations/:id/resources", webgin.Adapt(h.Inclusion.HandleListAdaptationResources))
	api.GET("/adaptations/:id/export", exportAdaptationRoute(h.Inclusion.ExportAdaptation))

	// Chat history
	api.GET("/chat/history/:contextId", webgin.Adapt(h.Inclusion.HandleGetChatHistory))

	// Dashboard
	api.GET("/dashboard/metrics", webgin.Adapt(h.Dashboard.HandleGetMetrics))

	// AI endpoints
	api.POST("/inclusion/recommend", webgin.Adapt(h.Inclusion.HandleRecommendDevice))
	api.POST("/inclusion/assist", webgin.Adapt(h.Inclusion.HandleAssistClassroom))
}

// exportAdaptationRoute serves a binary document download. It bypasses the
// JSON-only web.Response adapter so it can set Content-Type and
// Content-Disposition headers directly on the gin response.
func exportAdaptationRoute(uc inclusionuc.ExportAdaptation) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := webgin.NewRequest(c)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "invalid_id", "message": "invalid adaptation id"})
			return
		}

		format := c.Query("format")
		if format == "" {
			format = inclusionuc.ExportFormatPDF
		}

		doc, err := uc.Execute(req.Context(), inclusionuc.ExportAdaptationRequest{
			OrgID:        middleware.OrgID(req),
			AdaptationID: id,
			Format:       format,
		})
		if err != nil {
			resp := rest.HandleError(err)
			c.JSON(resp.Status, resp.Body)
			return
		}

		c.Header("Content-Disposition", `attachment; filename="`+doc.Filename+`"`)
		c.Data(http.StatusOK, doc.ContentType, doc.Data)
	}
}
