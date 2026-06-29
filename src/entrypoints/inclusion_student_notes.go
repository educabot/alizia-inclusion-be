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

type studentNoteResponse struct {
	ID        int64  `json:"id"`
	StudentID int64  `json:"student_id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Internal  bool   `json:"internal"`
	CreatedAt string `json:"created_at"`
}

func mapStudentNote(n entities.StudentNote) studentNoteResponse {
	return studentNoteResponse{
		ID:        n.ID,
		StudentID: n.StudentID,
		Content:   n.Content,
		Type:      n.Type,
		Internal:  n.Internal,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
}

func mapStudentNotes(ns []entities.StudentNote) []studentNoteResponse {
	out := make([]studentNoteResponse, len(ns))
	for i := range ns {
		out[i] = mapStudentNote(ns[i])
	}
	return out
}

type createStudentNoteBody struct {
	Content  string `json:"content"`
	Type     string `json:"type"`
	Internal *bool  `json:"internal"`
}

func (c *InclusionContainer) HandleListStudentNotes(req web.Request) web.Response {
	studentID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.ListStudentNotes.Execute(req.Context(), inclusion.ListStudentNotesRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: studentID,
		UserID:    middleware.UserID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentNotes(result))
}

func (c *InclusionContainer) HandleCreateStudentNote(req web.Request) web.Response {
	studentID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body createStudentNoteBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CreateStudentNote.Execute(req.Context(), inclusion.CreateStudentNoteRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: studentID,
		UserID:    middleware.UserID(req),
		Content:   body.Content,
		Type:      body.Type,
		Internal:  body.Internal,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusCreated, Body: mapStudentNote(*result)}
}
