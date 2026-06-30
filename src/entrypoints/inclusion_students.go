package entrypoints

import (
	"net/http"
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type createStudentBody struct {
	Name        string `json:"name"`
	ClassroomID *int64 `json:"classroom_id"`
}

type updateStudentBody struct {
	Name        *string `json:"name"`
	ClassroomID *int64  `json:"classroom_id"`
}

func (c *InclusionContainer) HandleListStudents(req web.Request) web.Response {
	var classroomID *int64
	if v := req.Query("classroom_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return rest.HandleError(err)
		}
		classroomID = &id
	}

	result, err := c.ListStudents.Execute(req.Context(), inclusion.ListStudentsRequest{
		OrgID:       middleware.OrgID(req),
		ClassroomID: classroomID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentsWithProfiles(result))
}

func (c *InclusionContainer) HandleGetStudent(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetStudentProfile.Execute(req.Context(), inclusion.GetStudentProfileRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentWithProfile(*result))
}

func (c *InclusionContainer) HandleCreateStudent(req web.Request) web.Response {
	var body createStudentBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CreateStudent.Execute(req.Context(), inclusion.CreateStudentRequest{
		OrgID:       middleware.OrgID(req),
		ClassroomID: body.ClassroomID,
		Name:        body.Name,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusCreated, Body: mapStudentWithProfile(*result)}
}

func (c *InclusionContainer) HandleUpdateStudent(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body updateStudentBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.UpdateStudent.Execute(req.Context(), inclusion.UpdateStudentRequest{
		OrgID:       middleware.OrgID(req),
		StudentID:   id,
		Name:        body.Name,
		ClassroomID: body.ClassroomID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapStudentWithProfile(*result))
}

func (c *InclusionContainer) HandleDeleteStudent(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	if err := c.DeleteStudent.Execute(req.Context(), inclusion.DeleteStudentRequest{
		OrgID:     middleware.OrgID(req),
		StudentID: id,
	}); err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusNoContent}
}
