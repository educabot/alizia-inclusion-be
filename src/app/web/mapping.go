package web

import (
	"net/http"
	"strconv"

	webgin "github.com/educabot/team-ai-toolkit/web/gin"
	"github.com/gin-gonic/gin"

	"github.com/educabot/alizia-inclusion-be/config"
	"github.com/educabot/alizia-inclusion-be/src/app/web/static"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func ConfigureMappings(engine *gin.Engine, h *entrypoints.WebHandlerContainer, cfg *config.Config) {
	// Assets estáticos públicos (imágenes de devices). Fuera del grupo /api/v1:
	// son públicos (el <img src> del FE no manda Authorization) y no llevan
	// auth ni tenant. Las imágenes están embebidas (ver package static), keyed
	// por product_code: /images/devices/ETE-XXXX-EB.png.
	registerStaticAssets(engine)

	api := engine.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(h.AuthMiddleware))
	api.Use(webgin.AdaptMiddleware(h.TenantMiddleware))

	aiRateLimit := webgin.AdaptMiddleware(middleware.RateLimitMiddleware(cfg.AIRateLimitPerHour))

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
	studentByID.GET("/notes", webgin.Adapt(h.Inclusion.HandleListStudentNotes))
	studentByID.POST("/notes", webgin.Adapt(h.Inclusion.HandleCreateStudentNote))

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
	api.DELETE("/chat/conversation/:id", webgin.Adapt(h.Inclusion.HandleDeleteConversation))
	api.PATCH("/chat/conversation/:id", webgin.Adapt(h.Inclusion.HandleRenameConversation))

	// Dashboard
	api.GET("/dashboard/metrics", webgin.Adapt(h.Dashboard.HandleGetMetrics))
	api.GET("/dashboard/ai-usage", webgin.Adapt(h.Dashboard.HandleGetAIUsage))

	// Apertura de sesión (router / Prompt 0) — sin LLM, no requiere rate limit de IA
	api.POST("/inclusion/open", webgin.Adapt(h.Inclusion.HandleOpenSession))

	// Context Assembler (HU-2) — arma el contexto del alumno/valija/tema; sin LLM
	api.POST("/inclusion/context", webgin.Adapt(h.Inclusion.HandleBuildContext))

	// RAG de contenido pedagógico (HU-3) — búsqueda keyword/full-text; sin LLM
	api.POST("/inclusion/search-content", webgin.Adapt(h.Inclusion.HandleSearchContent))

	// RAG híbrido (vector + texto + conceptos) sobre el corpus rag_*; sin LLM
	api.POST("/inclusion/search-content/hybrid", webgin.Adapt(h.Inclusion.HandleHybridSearch))

	// AI endpoints (rate-limited per organization)
	api.POST("/inclusion/recommend", aiRateLimit, webgin.Adapt(h.Inclusion.HandleRecommendDevice))
	api.POST("/inclusion/assist", aiRateLimit, webgin.Adapt(h.Inclusion.HandleAssistClassroom))
}

// registerStaticAssets monta el file server embebido en /images/*. Sirve los
// binarios desde el embed.FS del package static y agrega Cache-Control: las
// imágenes de devices son inmutables por nombre (product_code), así que el
// browser/CDN puede cachearlas con holgura.
func registerStaticAssets(engine *gin.Engine) {
	fileServer := http.StripPrefix("/images/", http.FileServer(http.FS(static.Images())))
	engine.GET("/images/*filepath", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=86400")
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
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

		// disposition=inline abre el archivo en el navegador (preview, tipo Drive)
		// en vez de forzar la descarga. Default: attachment.
		disposition := "attachment"
		if c.Query("disposition") == "inline" {
			disposition = "inline"
		}
		c.Header("Content-Disposition", disposition+`; filename="`+doc.Filename+`"`)
		c.Data(http.StatusOK, doc.ContentType, doc.Data)
	}
}
