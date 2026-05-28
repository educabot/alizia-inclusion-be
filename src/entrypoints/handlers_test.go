package entrypoints_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	cataloguc "github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	dashuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	mgmtuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

var testOrgID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func newTenantRequest() *web.MockRequest {
	req := web.NewMockRequest()
	req.Values[middleware.OrgIDKey] = testOrgID
	req.Values[middleware.UserIDKey] = int64(42)
	return req
}

// --- Catalog usecase mocks ---

type mockListRamps struct {
	fn func(ctx context.Context, req cataloguc.ListRampsRequest) ([]entities.Ramp, error)
}

func (m *mockListRamps) Execute(ctx context.Context, req cataloguc.ListRampsRequest) ([]entities.Ramp, error) {
	return m.fn(ctx, req)
}

type mockGetRamp struct {
	fn func(ctx context.Context, req cataloguc.GetRampRequest) (*entities.Ramp, error)
}

func (m *mockGetRamp) Execute(ctx context.Context, req cataloguc.GetRampRequest) (*entities.Ramp, error) {
	return m.fn(ctx, req)
}

type mockListDevices struct {
	fn func(ctx context.Context, req cataloguc.ListDevicesRequest) ([]entities.Device, error)
}

func (m *mockListDevices) Execute(ctx context.Context, req cataloguc.ListDevicesRequest) ([]entities.Device, error) {
	return m.fn(ctx, req)
}

type mockGetDevice struct {
	fn func(ctx context.Context, req cataloguc.GetDeviceRequest) (*entities.Device, error)
}

func (m *mockGetDevice) Execute(ctx context.Context, req cataloguc.GetDeviceRequest) (*entities.Device, error) {
	return m.fn(ctx, req)
}

// --- Management usecase mocks ---

type mockListClassrooms struct {
	fn func(ctx context.Context, req mgmtuc.ListClassroomsRequest) ([]entities.Classroom, error)
}

func (m *mockListClassrooms) Execute(ctx context.Context, req mgmtuc.ListClassroomsRequest) ([]entities.Classroom, error) {
	return m.fn(ctx, req)
}

type mockGetClassroom struct {
	fn func(ctx context.Context, req mgmtuc.GetClassroomRequest) (*entities.Classroom, error)
}

func (m *mockGetClassroom) Execute(ctx context.Context, req mgmtuc.GetClassroomRequest) (*entities.Classroom, error) {
	return m.fn(ctx, req)
}

type mockCreateClassroom struct {
	fn func(ctx context.Context, req mgmtuc.CreateClassroomRequest) (*entities.Classroom, error)
}

func (m *mockCreateClassroom) Execute(ctx context.Context, req mgmtuc.CreateClassroomRequest) (*entities.Classroom, error) {
	return m.fn(ctx, req)
}

type mockDeleteClassroom struct {
	fn func(ctx context.Context, req mgmtuc.DeleteClassroomRequest) error
}

func (m *mockDeleteClassroom) Execute(ctx context.Context, req mgmtuc.DeleteClassroomRequest) error {
	return m.fn(ctx, req)
}

type mockListTeachers struct {
	fn func(ctx context.Context, req mgmtuc.ListTeachersRequest) ([]entities.User, error)
}

func (m *mockListTeachers) Execute(ctx context.Context, req mgmtuc.ListTeachersRequest) ([]entities.User, error) {
	return m.fn(ctx, req)
}

// --- Auth usecase mocks ---

type mockGetMe struct {
	fn func(ctx context.Context, req authuc.GetMeRequest) (*entities.User, error)
}

func (m *mockGetMe) Execute(ctx context.Context, req authuc.GetMeRequest) (*entities.User, error) {
	return m.fn(ctx, req)
}

// --- Dashboard usecase mocks ---

type mockGetMetrics struct {
	fn func(ctx context.Context, req dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error)
}

func (m *mockGetMetrics) Execute(ctx context.Context, req dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error) {
	return m.fn(ctx, req)
}

// ==================== Catalog Handler Tests ====================

func TestHandleListRamps_ReturnsRamps(t *testing.T) {
	container := &entrypoints.CatalogContainer{
		ListRamps: &mockListRamps{fn: func(_ context.Context, req cataloguc.ListRampsRequest) ([]entities.Ramp, error) {
			assert.Equal(t, testOrgID, req.OrgID)
			return []entities.Ramp{{ID: 1, Name: "Comunicación"}}, nil
		}},
	}

	resp := container.HandleListRamps(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleListRamps_ReturnsError(t *testing.T) {
	container := &entrypoints.CatalogContainer{
		ListRamps: &mockListRamps{fn: func(_ context.Context, _ cataloguc.ListRampsRequest) ([]entities.Ramp, error) {
			return nil, errNotFound
		}},
	}

	resp := container.HandleListRamps(newTenantRequest())

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

func TestHandleGetRamp_ReturnsByID(t *testing.T) {
	container := &entrypoints.CatalogContainer{
		GetRamp: &mockGetRamp{fn: func(_ context.Context, req cataloguc.GetRampRequest) (*entities.Ramp, error) {
			assert.Equal(t, int64(5), req.RampID)
			return &entities.Ramp{ID: 5, Name: "Sensorial"}, nil
		}},
	}
	req := newTenantRequest()
	req.Params["id"] = "5"

	resp := container.HandleGetRamp(req)

	assert.Equal(t, http.StatusOK, resp.Status)
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
	container := &entrypoints.CatalogContainer{
		ListDevices: &mockListDevices{fn: func(_ context.Context, _ cataloguc.ListDevicesRequest) ([]entities.Device, error) {
			return []entities.Device{{ID: 1, Name: "Timer"}}, nil
		}},
	}

	resp := container.HandleListDevices(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleListDevices_FiltersByRampID(t *testing.T) {
	container := &entrypoints.CatalogContainer{
		ListDevices: &mockListDevices{fn: func(_ context.Context, req cataloguc.ListDevicesRequest) ([]entities.Device, error) {
			assert.NotNil(t, req.RampID)
			assert.Equal(t, int64(3), *req.RampID)
			return []entities.Device{}, nil
		}},
	}
	req := newTenantRequest()
	req.Queries["ramp_id"] = "3"

	resp := container.HandleListDevices(req)

	assert.Equal(t, http.StatusOK, resp.Status)
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
	container := &entrypoints.CatalogContainer{
		GetDevice: &mockGetDevice{fn: func(_ context.Context, req cataloguc.GetDeviceRequest) (*entities.Device, error) {
			return &entities.Device{ID: req.DeviceID, Name: "Pictogramas"}, nil
		}},
	}
	req := newTenantRequest()
	req.Params["id"] = "7"

	resp := container.HandleGetDevice(req)

	assert.Equal(t, http.StatusOK, resp.Status)
}

// ==================== Management Handler Tests ====================

func TestHandleListClassrooms_ReturnsClassrooms(t *testing.T) {
	container := &entrypoints.ManagementContainer{
		ListClassrooms: &mockListClassrooms{fn: func(_ context.Context, _ mgmtuc.ListClassroomsRequest) ([]entities.Classroom, error) {
			return []entities.Classroom{{ID: 1, Name: "3ro A"}}, nil
		}},
	}

	resp := container.HandleListClassrooms(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleGetClassroom_ReturnsClassroom(t *testing.T) {
	container := &entrypoints.ManagementContainer{
		GetClassroom: &mockGetClassroom{fn: func(_ context.Context, req mgmtuc.GetClassroomRequest) (*entities.Classroom, error) {
			return &entities.Classroom{ID: req.ClassroomID, Name: "4to B"}, nil
		}},
	}
	req := newTenantRequest()
	req.Params["id"] = "2"

	resp := container.HandleGetClassroom(req)

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleGetClassroom_Returns404ForNotFound(t *testing.T) {
	container := &entrypoints.ManagementContainer{
		GetClassroom: &mockGetClassroom{fn: func(_ context.Context, _ mgmtuc.GetClassroomRequest) (*entities.Classroom, error) {
			return nil, errNotFound
		}},
	}
	req := newTenantRequest()
	req.Params["id"] = "999"

	resp := container.HandleGetClassroom(req)

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

func TestHandleCreateClassroom_CreatesClassroom(t *testing.T) {
	grade := "3ro"
	container := &entrypoints.ManagementContainer{
		CreateClassroom: &mockCreateClassroom{fn: func(_ context.Context, req mgmtuc.CreateClassroomRequest) (*entities.Classroom, error) {
			return &entities.Classroom{ID: 10, Name: req.Name, Grade: req.Grade}, nil
		}},
	}
	req := newTenantRequest()
	req.BindJSONFn = func(dest any) error {
		b, _ := json.Marshal(map[string]any{"name": "5to A", "grade": &grade})
		return json.Unmarshal(b, dest)
	}

	resp := container.HandleCreateClassroom(req)

	assert.Equal(t, http.StatusCreated, resp.Status)
}

func TestHandleDeleteClassroom_DeletesClassroom(t *testing.T) {
	container := &entrypoints.ManagementContainer{
		DeleteClassroom: &mockDeleteClassroom{fn: func(_ context.Context, req mgmtuc.DeleteClassroomRequest) error {
			assert.Equal(t, int64(3), req.ClassroomID)
			return nil
		}},
	}
	req := newTenantRequest()
	req.Params["id"] = "3"

	resp := container.HandleDeleteClassroom(req)

	assert.Equal(t, http.StatusNoContent, resp.Status)
}

func TestHandleListTeachers_ReturnsTeachers(t *testing.T) {
	container := &entrypoints.ManagementContainer{
		ListTeachers: &mockListTeachers{fn: func(_ context.Context, _ mgmtuc.ListTeachersRequest) ([]entities.User, error) {
			return []entities.User{{ID: 1, Name: "Ana", Email: "ana@test.com", Role: "teacher"}}, nil
		}},
	}

	resp := container.HandleListTeachers(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

// ==================== Auth Handler Tests ====================

func TestHandleGetMe_ReturnsCurrentUser(t *testing.T) {
	container := &entrypoints.AuthContainer{
		GetMe: &mockGetMe{fn: func(_ context.Context, req authuc.GetMeRequest) (*entities.User, error) {
			assert.Equal(t, int64(42), req.UserID)
			return &entities.User{ID: 42, Name: "Test", Email: "test@test.com", Role: "teacher"}, nil
		}},
	}

	resp := container.HandleGetMe(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleGetMe_ReturnsErrorWhenNotFound(t *testing.T) {
	container := &entrypoints.AuthContainer{
		GetMe: &mockGetMe{fn: func(_ context.Context, _ authuc.GetMeRequest) (*entities.User, error) {
			return nil, errNotFound
		}},
	}

	resp := container.HandleGetMe(newTenantRequest())

	assert.Equal(t, http.StatusNotFound, resp.Status)
}

// ==================== Dashboard Handler Tests ====================

func TestHandleGetMetrics_ReturnsMetrics(t *testing.T) {
	container := &entrypoints.DashboardContainer{
		GetMetrics: &mockGetMetrics{fn: func(_ context.Context, _ dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error) {
			return &dashuc.GetMetricsResponse{
				TotalStudents:        10,
				StudentsWithProfiles: 5,
			}, nil
		}},
	}

	resp := container.HandleGetMetrics(newTenantRequest())

	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestHandleGetMetrics_ReturnsError(t *testing.T) {
	container := &entrypoints.DashboardContainer{
		GetMetrics: &mockGetMetrics{fn: func(_ context.Context, _ dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error) {
			return nil, errBadRequest
		}},
	}

	resp := container.HandleGetMetrics(newTenantRequest())

	assert.Equal(t, http.StatusBadRequest, resp.Status)
}
