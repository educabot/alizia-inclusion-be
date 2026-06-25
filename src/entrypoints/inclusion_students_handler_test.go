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

func TestHandleListStudents_ReturnsStudents(t *testing.T) {
	uc := &mockusecases.MockListStudents{}
	uc.On("Execute", mock.Anything, mock.Anything).Return([]entities.Student{{ID: 1, Name: "Lucía"}}, nil)
	container := &entrypoints.InclusionContainer{ListStudents: uc}

	resp := container.HandleListStudents(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleListStudents_FiltersByClassroomID(t *testing.T) {
	uc := &mockusecases.MockListStudents{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.ListStudentsRequest) bool {
		return req.ClassroomID != nil && *req.ClassroomID == 8
	})).Return([]entities.Student{}, nil)
	container := &entrypoints.InclusionContainer{ListStudents: uc}
	req := newTenantRequest()
	req.Queries["classroom_id"] = "8"

	resp := container.HandleListStudents(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListStudents_RejectsInvalidClassroomID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Queries["classroom_id"] = "abc"

	resp := container.HandleListStudents(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleGetStudent_ReturnsStudent(t *testing.T) {
	uc := &mockusecases.MockGetStudentProfile{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.GetStudentProfileRequest) bool {
		return req.StudentID == 5
	})).Return(&entities.Student{ID: 5, Name: "Mateo"}, nil)
	container := &entrypoints.InclusionContainer{GetStudentProfile: uc}
	req := newTenantRequest()
	req.Params["id"] = "5"

	resp := container.HandleGetStudent(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetStudent_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleGetStudent(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleCreateStudent_CreatesStudent(t *testing.T) {
	uc := &mockusecases.MockCreateStudent{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.CreateStudentRequest) bool {
		return req.Name == "Sofía" && req.ClassroomID == 3
	})).Return(&entities.Student{ID: 9, Name: "Sofía", ClassroomID: 3}, nil)
	container := &entrypoints.InclusionContainer{CreateStudent: uc}
	req := newTenantRequest()
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"name": "Sofía", "classroom_id": 3})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleCreateStudent(req)

	assert.Equal(t, http.StatusCreated, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleCreateStudent_ReturnsErrorWhenUsecaseFails(t *testing.T) {
	uc := &mockusecases.MockCreateStudent{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errBadRequest)
	container := &entrypoints.InclusionContainer{CreateStudent: uc}
	req := newTenantRequest()
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"name": "", "classroom_id": 0})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleCreateStudent(req)

	assert.Equal(t, http.StatusBadRequest, resp.Status)
}

func TestHandleUpdateStudent_UpdatesStudent(t *testing.T) {
	uc := &mockusecases.MockUpdateStudent{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.UpdateStudentRequest) bool {
		return req.StudentID == 7
	})).Return(&entities.Student{ID: 7, Name: "Renamed"}, nil)
	container := &entrypoints.InclusionContainer{UpdateStudent: uc}
	req := newTenantRequest()
	req.Params["id"] = "7"
	req.BindJSONFn = func(dest any) error {
		name := "Renamed"
		b, _ := json.Marshal(map[string]any{"name": &name})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleUpdateStudent(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleUpdateStudent_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "nope"

	resp := container.HandleUpdateStudent(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleDeleteStudent_DeletesStudent(t *testing.T) {
	uc := &mockusecases.MockDeleteStudent{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.DeleteStudentRequest) bool {
		return req.StudentID == 4
	})).Return(nil)
	container := &entrypoints.InclusionContainer{DeleteStudent: uc}
	req := newTenantRequest()
	req.Params["id"] = "4"

	resp := container.HandleDeleteStudent(req)

	assert.Equal(t, http.StatusNoContent, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleDeleteStudent_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "bad"

	resp := container.HandleDeleteStudent(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}
