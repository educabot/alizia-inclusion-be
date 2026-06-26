package entrypoints

import (
	"net/http"
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type adaptationResponse struct {
	ID                   int64    `json:"id"`
	StudentID            *int64   `json:"student_id,omitempty"`
	StudentName          string   `json:"student_name"`
	TeacherID            int64    `json:"teacher_id"`
	TeacherName          string   `json:"teacher_name"`
	DeviceID             *int64   `json:"device_id,omitempty"`
	DeviceName           *string  `json:"device_name,omitempty"`
	DeviceIDs            []int64  `json:"device_ids"`
	DeviceNames          []string `json:"device_names"`
	Title                string   `json:"title"`
	Subject              string   `json:"subject"`
	ActivityDescription  *string  `json:"activity_description,omitempty"`
	AdaptationStrategy   *string  `json:"adaptation_strategy,omitempty"`
	AdaptationType       string   `json:"adaptation_type"`
	Outcome              *string  `json:"outcome,omitempty"`
	Notes                *string  `json:"notes,omitempty"`
	Status               string   `json:"status"`
	SourceConversationID *int64   `json:"source_conversation_id,omitempty"`
	SourceMessageID      *int64   `json:"source_message_id,omitempty"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}

func mapAdaptation(a entities.Adaptation) adaptationResponse {
	resp := adaptationResponse{
		ID:                   a.ID,
		StudentID:            a.StudentID,
		TeacherID:            a.TeacherID,
		DeviceID:             a.DeviceID,
		Title:                a.Title,
		Subject:              a.Subject,
		ActivityDescription:  a.ActivityDescription,
		AdaptationStrategy:   a.AdaptationStrategy,
		AdaptationType:       a.AdaptationType,
		Outcome:              a.Outcome,
		Notes:                a.Notes,
		Status:               a.Status,
		SourceConversationID: a.SourceConversationID,
		SourceMessageID:      a.SourceMessageID,
		CreatedAt:            a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            a.UpdatedAt.Format(time.RFC3339),
		DeviceIDs:            make([]int64, 0),
		DeviceNames:          make([]string, 0),
	}
	if a.Student != nil {
		resp.StudentName = a.Student.Name
	}
	if a.Teacher != nil {
		resp.TeacherName = a.Teacher.Name
	}
	if a.Device != nil {
		resp.DeviceName = &a.Device.Name
	}
	for i := range a.Devices {
		d := &a.Devices[i]
		resp.DeviceIDs = append(resp.DeviceIDs, d.ID)
		resp.DeviceNames = append(resp.DeviceNames, d.Name)
	}
	return resp
}

func mapAdaptations(as []entities.Adaptation) []adaptationResponse {
	out := make([]adaptationResponse, len(as))
	for i := range as {
		out[i] = mapAdaptation(as[i])
	}
	return out
}

type createAdaptationBody struct {
	StudentID            *int64  `json:"student_id"`
	DeviceID             *int64  `json:"device_id"`
	DeviceIDs            []int64 `json:"device_ids"`
	Title                *string `json:"title"`
	Subject              string  `json:"subject"`
	ActivityDescription  *string `json:"activity_description"`
	AdaptationStrategy   *string `json:"adaptation_strategy"`
	AdaptationType       string  `json:"adaptation_type"`
	Notes                *string `json:"notes"`
	SourceConversationID *int64  `json:"source_conversation_id"`
	SourceMessageID      *int64  `json:"source_message_id"`
}

type updateAdaptationBody struct {
	DeviceID            *int64   `json:"device_id"`
	DeviceIDs           *[]int64 `json:"device_ids"`
	Title               *string  `json:"title"`
	Subject             *string  `json:"subject"`
	ActivityDescription *string  `json:"activity_description"`
	AdaptationStrategy  *string  `json:"adaptation_strategy"`
	AdaptationType      *string  `json:"adaptation_type"`
	Outcome             *string  `json:"outcome"`
	Notes               *string  `json:"notes"`
	Status              *string  `json:"status"`
}

func (c *InclusionContainer) HandleListAdaptations(req web.Request) web.Response {
	var studentID *int64
	if v := req.Query("student_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return rest.HandleError(err)
		}
		studentID = &id
	}

	result, err := c.ListAdaptations.Execute(req.Context(), inclusion.ListAdaptationsRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: studentID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapAdaptations(result))
}

func (c *InclusionContainer) HandleGetAdaptation(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetAdaptation.Execute(req.Context(), inclusion.GetAdaptationRequest{
		OrgID:        middleware.OrgID(req),
		AdaptationID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapAdaptation(*result))
}

func (c *InclusionContainer) HandleCreateAdaptation(req web.Request) web.Response {
	var body createAdaptationBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	title := ""
	if body.Title != nil {
		title = *body.Title
	}

	result, err := c.CreateAdaptation.Execute(req.Context(), inclusion.CreateAdaptationRequest{
		OrgID:                middleware.OrgID(req),
		StudentID:            body.StudentID,
		TeacherID:            middleware.UserID(req),
		DeviceID:             body.DeviceID,
		DeviceIDs:            body.DeviceIDs,
		Title:                title,
		Subject:              body.Subject,
		ActivityDescription:  body.ActivityDescription,
		AdaptationStrategy:   body.AdaptationStrategy,
		AdaptationType:       body.AdaptationType,
		Notes:                body.Notes,
		SourceConversationID: body.SourceConversationID,
		SourceMessageID:      body.SourceMessageID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusCreated, Body: mapAdaptation(*result)}
}

func (c *InclusionContainer) HandleUpdateAdaptation(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body updateAdaptationBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.UpdateAdaptation.Execute(req.Context(), inclusion.UpdateAdaptationRequest{
		OrgID:               middleware.OrgID(req),
		AdaptationID:        id,
		DeviceID:            body.DeviceID,
		DeviceIDs:           body.DeviceIDs,
		Title:               body.Title,
		Subject:             body.Subject,
		ActivityDescription: body.ActivityDescription,
		AdaptationStrategy:  body.AdaptationStrategy,
		AdaptationType:      body.AdaptationType,
		Outcome:             body.Outcome,
		Notes:               body.Notes,
		Status:              body.Status,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapAdaptation(*result))
}

func (c *InclusionContainer) HandleDeleteAdaptation(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	if err := c.DeleteAdaptation.Execute(req.Context(), inclusion.DeleteAdaptationRequest{
		OrgID:        middleware.OrgID(req),
		AdaptationID: id,
	}); err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusNoContent}
}
