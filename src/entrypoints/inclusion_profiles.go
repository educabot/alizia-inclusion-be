package entrypoints

import (
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type studentProfileResponse struct {
	ID              int64    `json:"id"`
	StudentID       int64    `json:"student_id"`
	StudentName     string   `json:"student_name"`
	IsTransitory    bool     `json:"is_transitory"`
	Difficulties    []string `json:"difficulties"`
	FreeDescription *string  `json:"free_description,omitempty"`
}

type studentWithProfileResponse struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	ClassroomID *int64                  `json:"classroom_id,omitempty"`
	Profile     *studentProfileResponse `json:"profile,omitempty"`
	CreatedAt   string                  `json:"created_at"`
}

func mapStudentWithProfile(s entities.Student) studentWithProfileResponse {
	resp := studentWithProfileResponse{
		ID:          s.ID,
		Name:        s.Name,
		ClassroomID: s.ClassroomID,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
	}
	if s.Profile != nil {
		diffs := make([]string, len(s.Profile.Difficulties))
		copy(diffs, s.Profile.Difficulties)
		resp.Profile = &studentProfileResponse{
			ID:              s.Profile.ID,
			StudentID:       s.Profile.StudentID,
			StudentName:     s.Name,
			IsTransitory:    s.Profile.IsTransitory,
			Difficulties:    diffs,
			FreeDescription: s.Profile.FreeDescription,
		}
	}
	return resp
}

func mapStudentsWithProfiles(ss []entities.Student) []studentWithProfileResponse {
	out := make([]studentWithProfileResponse, len(ss))
	for i := range ss {
		out[i] = mapStudentWithProfile(ss[i])
	}
	return out
}

type upsertProfileBody struct {
	IsTransitory    bool     `json:"is_transitory"`
	Difficulties    []string `json:"difficulties"`
	FreeDescription *string  `json:"free_description"`
}

func (c *InclusionContainer) HandleGetStudentProfile(req web.Request) web.Response {
	studentID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetStudentProfile.Execute(req.Context(), inclusion.GetStudentProfileRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: studentID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentWithProfile(*result))
}

func (c *InclusionContainer) HandleUpsertStudentProfile(req web.Request) web.Response {
	studentID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body upsertProfileBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	_, err = c.UpsertStudentProfile.Execute(req.Context(), inclusion.UpsertStudentProfileRequest{
		OrgID:           middleware.OrgID(req),
		StudentID:       studentID,
		IsTransitory:    body.IsTransitory,
		Difficulties:    body.Difficulties,
		FreeDescription: body.FreeDescription,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	student, err := c.GetStudentProfile.Execute(req.Context(), inclusion.GetStudentProfileRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: studentID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentWithProfile(*student))
}

func (c *InclusionContainer) HandleListClassroomStudents(req web.Request) web.Response {
	classroomID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.ListClassroomStudents.Execute(req.Context(), inclusion.ListClassroomStudentsRequest{
		OrgID:       middleware.OrgID(req),
		ClassroomID: classroomID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentsWithProfiles(result))
}
