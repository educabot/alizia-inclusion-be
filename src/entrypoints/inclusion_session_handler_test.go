package entrypoints_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	mockusecases "github.com/educabot/alizia-inclusion-be/src/mocks/usecases"
)

func TestHandleOpenSession_ReturnsGreeting(t *testing.T) {
	uc := &mockusecases.MockOpenSession{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return(&inclusionuc.OpenSessionResponse{Greeting: "¡Hola!", NeedsDimension: true}, nil)
	container := &entrypoints.InclusionContainer{OpenSession: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleOpenSession(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleCloseSession_ReturnsResponse(t *testing.T) {
	uc := &mockusecases.MockCloseSession{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return(&inclusionuc.CloseSessionResponse{}, nil)
	container := &entrypoints.InclusionContainer{CloseSession: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleCloseSession(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleBuildContext_ReturnsContext(t *testing.T) {
	uc := &mockusecases.MockBuildPromptContext{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return(&inclusionuc.PromptContext{}, nil)
	container := &entrypoints.InclusionContainer{BuildPromptContext: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleBuildContext(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleSearchContent_ReturnsResults(t *testing.T) {
	uc := &mockusecases.MockSearchPedagogicalContent{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return(&inclusionuc.SearchContentResponse{}, nil)
	container := &entrypoints.InclusionContainer{SearchPedagogicalContent: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleSearchContent(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleSearchContent_ReturnsErrorWhenUsecaseFails(t *testing.T) {
	uc := &mockusecases.MockSearchPedagogicalContent{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errBadRequest)
	container := &entrypoints.InclusionContainer{SearchPedagogicalContent: uc}
	req := newTenantRequest()
	req.BindJSONFn = noBody()

	resp := container.HandleSearchContent(req)

	assert.Equal(t, http.StatusBadRequest, resp.Status)
}

func TestHandleGetChatHistory_ReturnsConversations(t *testing.T) {
	uc := &mockusecases.MockGetChatHistory{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return([]entities.Conversation{{ID: 1, Mode: "assist"}}, nil)
	container := &entrypoints.InclusionContainer{GetChatHistory: uc}
	req := newTenantRequest()
	req.Params["contextId"] = "assist"

	resp := container.HandleGetChatHistory(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetConversation_ReturnsConversation(t *testing.T) {
	uc := &mockusecases.MockGetConversation{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.GetConversationRequest) bool {
		return req.ConversationID == 42
	})).Return(&entities.Conversation{ID: 42, Mode: "assist"}, nil)
	container := &entrypoints.InclusionContainer{GetConversation: uc}
	req := newTenantRequest()
	req.Params["id"] = "42"

	resp := container.HandleGetConversation(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetConversation_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleGetConversation(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleListAdaptationResources_ReturnsResources(t *testing.T) {
	uc := &mockusecases.MockListAdaptationResources{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req inclusionuc.ListAdaptationResourcesRequest) bool {
		return req.AdaptationID == 4
	})).Return([]entities.AdaptationResource{{ID: 1, AdaptationID: 4, Title: "Ficha"}}, nil)
	container := &entrypoints.InclusionContainer{ListAdaptationResources: uc}
	req := newTenantRequest()
	req.Params["id"] = "4"

	resp := container.HandleListAdaptationResources(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListAdaptationResources_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.InclusionContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleListAdaptationResources(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}
