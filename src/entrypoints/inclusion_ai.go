package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type recommendBody struct {
	ConversationID int64                   `json:"conversation_id"`
	StudentID      int64                   `json:"student_id"`
	Subject        string                  `json:"subject"`
	Objective      string                  `json:"objective"`
	Duration       string                  `json:"duration"`
	Dynamic        string                  `json:"dynamic"`
	Materials      string                  `json:"materials"`
	History        []providers.ChatMessage `json:"history"`
}

type assistBody struct {
	ConversationID int64                   `json:"conversation_id"`
	ClassroomID    int64                   `json:"classroom_id"`
	StudentID      *int64                  `json:"student_id"`
	Message        string                  `json:"message"`
	Mode           string                  `json:"mode"`
	Dimension      string                  `json:"dimension"`
	History        []providers.ChatMessage `json:"history"`
}

func (c *InclusionContainer) HandleRecommendDevice(req web.Request) web.Response {
	var body recommendBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.RecommendDevice.Execute(req.Context(), inclusion.RecommendDeviceRequest{
		OrgID:          middleware.OrgID(req),
		UserID:         middleware.UserID(req),
		ConversationID: body.ConversationID,
		StudentID:      body.StudentID,
		Subject:        body.Subject,
		Objective:      body.Objective,
		Duration:       body.Duration,
		Dynamic:        body.Dynamic,
		Materials:      body.Materials,
		History:        body.History,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}

func (c *InclusionContainer) HandleAssistClassroom(req web.Request) web.Response {
	var body assistBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.AssistClassroom.Execute(req.Context(), inclusion.AssistClassroomRequest{
		OrgID:          middleware.OrgID(req),
		UserID:         middleware.UserID(req),
		ConversationID: body.ConversationID,
		ClassroomID:    body.ClassroomID,
		StudentID:      body.StudentID,
		Message:        body.Message,
		Mode:           body.Mode,
		Dimension:      body.Dimension,
		History:        body.History,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
