package entrypoints_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	cataloguc "github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	dashuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	mgmtuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	mockusecases "github.com/educabot/alizia-inclusion-be/src/mocks/usecases"
)

var testOrgID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func newTenantRequest() *web.MockRequest {
	req := web.NewMockRequest()
	req.Values[middleware.OrgIDKey] = testOrgID
	req.Values[middleware.UserIDKey] = int64(42)
	return req
}

// ==================== Catalog Handler Tests ====================

func TestHandleListRamps_ReturnsRamps(t *testing.T) {
	uc := &mockusecases.MockListRamps{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req cataloguc.ListRampsRequest) bool {
		return req.OrgID == testOrgID
	})).Return([]entities.Ramp{{ID: 1, Name: "Comunicación"}}, nil)
	container := &entrypoints.CatalogContainer{ListRamps: uc}

	resp := container.HandleListRamps(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListRamps_ReturnsError(t *testing.T) {
	uc := &mockusecases.MockListRamps{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errNotFound)
	container := &entrypoints.CatalogContainer{ListRamps: uc}

	resp := container.HandleListRamps(newTenantRequest())

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

func TestHandleGetRamp_ReturnsByID(t *testing.T) {
	uc := &mockusecases.MockGetRamp{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req cataloguc.GetRampRequest) bool {
		return req.RampID == 5
	})).Return(&entities.Ramp{ID: 5, Name: "Sensorial"}, nil)
	container := &entrypoints.CatalogContainer{GetRamp: uc}
	req := newTenantRequest()
	req.Params["id"] = "5"

	resp := container.HandleGetRamp(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetRamp_RejectsInvalidID(t *testing.T) {
	container := &entrypoints.CatalogContainer{}
	req := newTenantRequest()
	req.Params["id"] = "abc"

	resp := container.HandleGetRamp(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleListDevices_ReturnsDevices(t *testing.T) {
	uc := &mockusecases.MockListDevices{}
	uc.On("Execute", mock.Anything, mock.Anything).Return([]entities.Device{{ID: 1, Name: "Timer"}}, nil)
	container := &entrypoints.CatalogContainer{ListDevices: uc}

	resp := container.HandleListDevices(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleListDevices_FiltersByRampID(t *testing.T) {
	uc := &mockusecases.MockListDevices{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req cataloguc.ListDevicesRequest) bool {
		return req.RampID != nil && *req.RampID == 3
	})).Return([]entities.Device{}, nil)
	container := &entrypoints.CatalogContainer{ListDevices: uc}
	req := newTenantRequest()
	req.Queries["ramp_id"] = "3"

	resp := container.HandleListDevices(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListDevices_RejectsInvalidRampID(t *testing.T) {
	container := &entrypoints.CatalogContainer{}
	req := newTenantRequest()
	req.Queries["ramp_id"] = "xyz"

	resp := container.HandleListDevices(req)

	assert.NotEqual(t, http.StatusOK, resp.Status)
	assert.NotEqual(t, 0, resp.Status)
}

func TestHandleGetDevice_ReturnsDevice(t *testing.T) {
	uc := &mockusecases.MockGetDevice{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req cataloguc.GetDeviceRequest) bool {
		return req.DeviceID == 7
	})).Return(&entities.Device{ID: 7, Name: "Pictogramas"}, nil)
	container := &entrypoints.CatalogContainer{GetDevice: uc}
	req := newTenantRequest()
	req.Params["id"] = "7"

	resp := container.HandleGetDevice(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

// ==================== Management Handler Tests ====================

func TestHandleListClassrooms_ReturnsClassrooms(t *testing.T) {
	uc := &mockusecases.MockListClassrooms{}
	uc.On("Execute", mock.Anything, mock.Anything).Return([]entities.Classroom{{ID: 1, Name: "3ro A"}}, nil)
	container := &entrypoints.ManagementContainer{ListClassrooms: uc}

	resp := container.HandleListClassrooms(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleGetClassroom_ReturnsClassroom(t *testing.T) {
	uc := &mockusecases.MockGetClassroom{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req mgmtuc.GetClassroomRequest) bool {
		return req.ClassroomID == 2
	})).Return(&entities.Classroom{ID: 2, Name: "4to B"}, nil)
	container := &entrypoints.ManagementContainer{GetClassroom: uc}
	req := newTenantRequest()
	req.Params["id"] = "2"

	resp := container.HandleGetClassroom(req)

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetClassroom_Returns404ForNotFound(t *testing.T) {
	uc := &mockusecases.MockGetClassroom{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errNotFound)
	container := &entrypoints.ManagementContainer{GetClassroom: uc}
	req := newTenantRequest()
	req.Params["id"] = "999"

	resp := container.HandleGetClassroom(req)

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

func TestHandleCreateClassroom_CreatesClassroom(t *testing.T) {
	grade := "3ro"
	uc := &mockusecases.MockCreateClassroom{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req mgmtuc.CreateClassroomRequest) bool {
		return req.Name == "5to A"
	})).Return(&entities.Classroom{ID: 10, Name: "5to A", Grade: &grade}, nil)
	container := &entrypoints.ManagementContainer{CreateClassroom: uc}
	req := newTenantRequest()
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"name": "5to A", "grade": &grade})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleCreateClassroom(req)

	assert.Equal(t, http.StatusCreated, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleDeleteClassroom_DeletesClassroom(t *testing.T) {
	uc := &mockusecases.MockDeleteClassroom{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req mgmtuc.DeleteClassroomRequest) bool {
		return req.ClassroomID == 3
	})).Return(nil)
	container := &entrypoints.ManagementContainer{DeleteClassroom: uc}
	req := newTenantRequest()
	req.Params["id"] = "3"

	resp := container.HandleDeleteClassroom(req)

	assert.Equal(t, http.StatusNoContent, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleListTeachers_ReturnsTeachers(t *testing.T) {
	uc := &mockusecases.MockListTeachers{}
	uc.On("Execute", mock.Anything, mock.Anything).
		Return([]entities.User{{ID: 1, Name: "Ana", Email: "ana@test.com", Role: "teacher"}}, nil)
	container := &entrypoints.ManagementContainer{ListTeachers: uc}

	resp := container.HandleListTeachers(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

// ==================== Auth Handler Tests ====================

func TestHandleGetMe_ReturnsCurrentUser(t *testing.T) {
	uc := &mockusecases.MockGetMe{}
	uc.On("Execute", mock.Anything, mock.MatchedBy(func(req authuc.GetMeRequest) bool {
		return req.UserID == 42
	})).Return(&entities.User{ID: 42, Name: "Test", Email: "test@test.com", Role: "teacher"}, nil)
	container := &entrypoints.AuthContainer{GetMe: uc}

	resp := container.HandleGetMe(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
	uc.AssertExpectations(t)
}

func TestHandleGetMe_ReturnsErrorWhenNotFound(t *testing.T) {
	uc := &mockusecases.MockGetMe{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errNotFound)
	container := &entrypoints.AuthContainer{GetMe: uc}

	resp := container.HandleGetMe(newTenantRequest())

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

// ==================== Dashboard Handler Tests ====================

func TestHandleGetMetrics_ReturnsMetrics(t *testing.T) {
	uc := &mockusecases.MockGetMetrics{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(&dashuc.GetMetricsResponse{
		TotalStudents:        10,
		StudentsWithProfiles: 5,
	}, nil)
	container := &entrypoints.DashboardContainer{GetMetrics: uc}

	resp := container.HandleGetMetrics(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleGetMetrics_ReturnsError(t *testing.T) {
	uc := &mockusecases.MockGetMetrics{}
	uc.On("Execute", mock.Anything, mock.Anything).Return(nil, errBadRequest)
	container := &entrypoints.DashboardContainer{GetMetrics: uc}

	resp := container.HandleGetMetrics(newTenantRequest())

	assert.Equal(t, http.StatusBadRequest, resp.Status)
}
