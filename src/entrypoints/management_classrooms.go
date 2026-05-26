package entrypoints

import (
	"net/http"
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type classroomResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Grade        *string `json:"grade,omitempty"`
	Section      *string `json:"section,omitempty"`
	StudentCount int     `json:"student_count"`
}

func mapClassroom(c entities.Classroom) classroomResponse {
	return classroomResponse{
		ID:           c.ID,
		Name:         c.Name,
		Grade:        c.Grade,
		Section:      c.Section,
		StudentCount: len(c.Students),
	}
}

func mapClassrooms(cs []entities.Classroom) []classroomResponse {
	out := make([]classroomResponse, len(cs))
	for i := range cs {
		out[i] = mapClassroom(cs[i])
	}
	return out
}

type createClassroomBody struct {
	Name    string  `json:"name"`
	Grade   *string `json:"grade"`
	Section *string `json:"section"`
}

type updateClassroomBody struct {
	Name    *string `json:"name"`
	Grade   *string `json:"grade"`
	Section *string `json:"section"`
}

func (c *ManagementContainer) HandleListClassrooms(req web.Request) web.Response {
	result, err := c.ListClassrooms.Execute(req.Context(), management.ListClassroomsRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapClassrooms(result))
}

func (c *ManagementContainer) HandleGetClassroom(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetClassroom.Execute(req.Context(), management.GetClassroomRequest{
		OrgID:       middleware.OrgID(req),
		ClassroomID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapClassroom(*result))
}

func (c *ManagementContainer) HandleCreateClassroom(req web.Request) web.Response {
	var body createClassroomBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CreateClassroom.Execute(req.Context(), management.CreateClassroomRequest{
		OrgID:   middleware.OrgID(req),
		Name:    body.Name,
		Grade:   body.Grade,
		Section: body.Section,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusCreated, Body: mapClassroom(*result)}
}

func (c *ManagementContainer) HandleUpdateClassroom(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body updateClassroomBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.UpdateClassroom.Execute(req.Context(), management.UpdateClassroomRequest{
		OrgID:       middleware.OrgID(req),
		ClassroomID: id,
		Name:        body.Name,
		Grade:       body.Grade,
		Section:     body.Section,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapClassroom(*result))
}

func (c *ManagementContainer) HandleDeleteClassroom(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	if err := c.DeleteClassroom.Execute(req.Context(), management.DeleteClassroomRequest{
		OrgID:       middleware.OrgID(req),
		ClassroomID: id,
	}); err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusNoContent}
}
