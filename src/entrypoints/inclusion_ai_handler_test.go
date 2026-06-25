package entrypoints_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	mockusecases "github.com/educabot/alizia-inclusion-be/src/mocks/usecases"
)

func noBody() func(any) error { return func(any) error { return nil } }

func TestHandleRecommendDevice_ReturnsRecommendation(t *testing.T) {
	uc := &mockusecases.MockRecommendDevice{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return(&inclusionuc.RecommendDeviceResponse{Response: "Probá el timer visual"}, nil)
	container := &entrypoints.InclusionContainer{RecommendDevice: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleRecommendDevice(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleRecommendDevice_ReturnsErrorWhenUsecaseFails(t *testing.T) {
	uc := &mockusecases.MockRecommendDevice{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errBadRequest)
	container := &entrypoints.InclusionContainer{RecommendDevice: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleRecommendDevice(req)

	assert.Equal(t, http.StatusBadRequest, resp.Status)
}

func TestHandleAssistClassroom_ReturnsResponse(t *testing.T) {
	uc := &mockusecases.MockAssistClassroom{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return(&inclusionuc.AssistClassroomResponse{Response: "Contame qué pasó"}, nil)
	container := &entrypoints.InclusionContainer{AssistClassroom: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleAssistClassroom(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleAssistClassroom_ReturnsErrorWhenUsecaseFails(t *testing.T) {
	uc := &mockusecases.MockAssistClassroom{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errNotFound)
	container := &entrypoints.InclusionContainer{AssistClassroom: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleAssistClassroom(req)

	assert.Equal(t, http.StatusNotFound, resp.Status)
}
