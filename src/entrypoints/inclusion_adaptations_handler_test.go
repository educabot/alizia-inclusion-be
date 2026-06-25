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

func TestHandleListAdaptations_ReturnsAdaptations(t *testing.T) {
	uc := &mockusecases.MockListAdaptations{}
	uc.On("Execute", mock.Anything, mock.Anything).Return([]entities.Adaptation{{ID: 1, Subject: "Matemática"}}, nil)
	container := &entrypoints.InclusionContainer{ListAdaptations: uc}

	resp := container.HandleListAdaptations(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleListAdaptations_DefaultsToCurrentTeacherScope(t *testing.T) {
	uc := &mockusecases.MockListAdaptations{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.ListAdaptationsRequest) bool {
		return req.TeacherID != nil && *req.TeacherID == 42
	})).Return([]entities.Adaptation{}, nil)
	container := &entrypoints.InclusionContainer{ListAdaptations: uc}

	resp := container.HandleListAdaptations(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListAdaptations_AllTrueDropsTeacherScope(t *testing.T) {
	uc := &mockusecases.MockListAdaptations{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.ListAdaptationsRequest) bool {
		return req.TeacherID == nil
	})).Return([]entities.Adaptation{}, nil)
	container := &entrypoints.InclusionContainer{ListAdaptations: uc}
	req := newTenantRequest()
	req.Queries["all"] = "true"

	resp := container.HandleListAdaptations(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListAdaptations_RejectsInvalidStudentID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Queries["student_id"] = "abc"

	resp := container.HandleListAdaptations(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleGetAdaptation_ReturnsAdaptation(t *testing.T) {
	uc := &mockusecases.MockGetAdaptation{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.GetAdaptationRequest) bool {
		return req.AdaptationID == 5
	})).Return(&entities.Adaptation{ID: 5, Subject: "Lengua"}, nil)
	container := &entrypoints.InclusionContainer{GetAdaptation: uc}
	req := newTenantRequest()
	req.Params["id"] = "5"

	resp := container.HandleGetAdaptation(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetAdaptation_Returns404ForNotFound(t *testing.T) {
	uc := &mockusecases.MockGetAdaptation{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errNotFound)
	container := &entrypoints.InclusionContainer{GetAdaptation: uc}
	req := newTenantRequest()
	req.Params["id"] = "999"

	resp := container.HandleGetAdaptation(req)

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

func TestHandleGetAdaptation_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleGetAdaptation(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleCreateAdaptation_CreatesAdaptation(t *testing.T) {
	uc := &mockusecases.MockCreateAdaptation{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.CreateAdaptationRequest) bool {
		return req.StudentID == 1 && req.TeacherID == 42
	})).Return(&entities.Adaptation{ID: 10, Subject: "Matemática"}, nil)
	container := &entrypoints.InclusionContainer{CreateAdaptation: uc}
	req := newTenantRequest()
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"student_id": 1, "subject": "Matemática", "adaptation_type": "estrategia_aula"})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleCreateAdaptation(req)

	assert.Equal(t, http.StatusCreated, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleUpdateAdaptation_UpdatesAdaptation(t *testing.T) {
	uc := &mockusecases.MockUpdateAdaptation{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.UpdateAdaptationRequest) bool {
		return req.AdaptationID == 7
	})).Return(&entities.Adaptation{ID: 7, Subject: "Lengua"}, nil)
	container := &entrypoints.InclusionContainer{UpdateAdaptation: uc}
	req := newTenantRequest()
	req.Params["id"] = "7"
	req.BindJSONFn = func(dest any) error {
		status := "funciono"
		b, _ := json.Marshal(map[string]any{"status": &status})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleUpdateAdaptation(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleUpdateAdaptation_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "nope"

	resp := container.HandleUpdateAdaptation(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleDeleteAdaptation_DeletesAdaptation(t *testing.T) {
	uc := &mockusecases.MockDeleteAdaptation{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.DeleteAdaptationRequest) bool {
		return req.AdaptationID == 3
	})).Return(nil)
	container := &entrypoints.InclusionContainer{DeleteAdaptation: uc}
	req := newTenantRequest()
	req.Params["id"] = "3"

	resp := container.HandleDeleteAdaptation(req)

	assert.Equal(t, http.StatusNoContent, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleDeleteAdaptation_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "bad"

	resp := container.HandleDeleteAdaptation(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}
