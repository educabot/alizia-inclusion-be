package entrypoints_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	mockusecases "github.com/educabot/alizia-inclusion-be/src/mocks/usecases"
)

func TestHandleGetStudentProfile_ReturnsStudent(t *testing.T) {
	uc := &mockusecases.MockGetStudentProfile{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.GetStudentProfileRequest) bool {
		return req.StudentID == 5
	})).Return(&entities.Student{ID: 5, Name: "Mateo"}, nil)
	container := &entrypoints.InclusionContainer{GetStudentProfile: uc}
	req := newTenantRequest()
	req.Params["id"] = "5"

	resp := container.HandleGetStudentProfile(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetStudentProfile_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleGetStudentProfile(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleUpsertStudentProfile_UpsertsThenReturnsStudent(t *testing.T) {
	upsert := &mockusecases.MockUpsertStudentProfile{}
	upsert.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.UpsertStudentProfileRequest) bool {
		return req.StudentID == 8 && req.IsTransitory
	})).Return(&entities.StudentProfile{ID: 1, StudentID: 8}, nil)
	// The handler re-reads the student after upserting to return the full view.
	get := &mockusecases.MockGetStudentProfile{}
	get.On("Execute", mock.Anything, mock.Anything).Return(&entities.Student{ID: 8, Name: "Sofía"}, nil)
	container := &entrypoints.InclusionContainer{UpsertStudentProfile: upsert, GetStudentProfile: get}
	req := newTenantRequest()
	req.Params["id"] = "8"
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"is_transitory": true, "difficulties": []string{"se_distrae"}})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleUpsertStudentProfile(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	upsert.AssertExpectations(t)
	get.AssertExpectations(t)
}

func TestHandleUpsertStudentProfile_ReturnsErrorWhenUpsertFails(t *testing.T) {
	upsert := &mockusecases.MockUpsertStudentProfile{}
	upsert.On("Execute", mock.Anything, mock.Anything).Return(nil, errBadRequest)
	container := &entrypoints.InclusionContainer{UpsertStudentProfile: upsert}
	req := newTenantRequest()
	req.Params["id"] = "8"
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"is_transitory": false})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleUpsertStudentProfile(req)

	assert.Equal(t, http.StatusBadRequest, resp.Status)
}

func TestHandleUpsertStudentProfile_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "nope"

	resp := container.HandleUpsertStudentProfile(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleListClassroomStudents_ReturnsStudents(t *testing.T) {
	uc := &mockusecases.MockListClassroomStudents{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.ListClassroomStudentsRequest) bool {
		return req.ClassroomID == 2
	})).Return([]entities.Student{{ID: 1, Name: "Ana"}}, nil)
	container := &entrypoints.InclusionContainer{ListClassroomStudents: uc}
	req := newTenantRequest()
	req.Params["id"] = "2"

	resp := container.HandleListClassroomStudents(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListClassroomStudents_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleListClassroomStudents(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}
