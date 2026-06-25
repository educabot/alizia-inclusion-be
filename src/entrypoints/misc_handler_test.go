package entrypoints_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	dashuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	mgmtuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	mockusecases "github.com/educabot/alizia-inclusion-be/src/mocks/usecases"
)

func TestHandleGetAIUsage_ReturnsUsage(t *testing.T) {
	uc := &mockusecases.MockGetAIUsage{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(&dashuc.GetAIUsageResponse{}, nil)
	container := &entrypoints.DashboardContainer{GetAIUsage: uc}

	resp := container.HandleGetAIUsage(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetAIUsage_ParsesDaysQuery(t *testing.T) {
	uc := &mockusecases.MockGetAIUsage{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req dashuc.GetAIUsageRequest) bool {
		return req.Days == 30
	})).Return(&dashuc.GetAIUsageResponse{}, nil)
	container := &entrypoints.DashboardContainer{GetAIUsage: uc}
	req := newTenantRequest()
	req.Queries["days"] = "30"

	resp := container.HandleGetAIUsage(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleUpdateClassroom_UpdatesClassroom(t *testing.T) {
	uc := &mockusecases.MockUpdateClassroom{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req mgmtuc.UpdateClassroomRequest) bool {
		return req.ClassroomID == 6
	})).Return(&entities.Classroom{ID: 6, Name: "6to A"}, nil)
	container := &entrypoints.ManagementContainer{UpdateClassroom: uc}
	req := newTenantRequest()
	req.Params["id"] = "6"
	req.BindJSONFn = func(dest any) error {
		name := "6to A"
		b, _ := json.Marshal(map[string]any{"name": &name})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleUpdateClassroom(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleUpdateClassroom_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.ManagementContainer{}
	req := newTenantRequest()
	req.Params["id"] = "x"

	resp := container.HandleUpdateClassroom(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}
