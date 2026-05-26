package entrypoints_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"

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

type mockUpdateClassroom struct {
	fn func(ctx context.Context, req mgmtuc.UpdateClassroomRequest) (*entities.Classroom, error)
}
func (m *mockUpdateClassroom) Execute(ctx context.Context, req mgmtuc.UpdateClassroomRequest) (*entities.Classroom, error) {
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

func TestHandleListRamps(t *testing.T) {
	t.Run("returns ramps", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{
			ListRamps: &mockListRamps{fn: func(_ context.Context, req cataloguc.ListRampsRequest) ([]entities.Ramp, error) {
				if req.OrgID != testOrgID {
					t.Errorf("expected org %s, got %s", testOrgID, req.OrgID)
				}
				return []entities.Ramp{{ID: 1, Name: "Comunicación"}}, nil
			}},
		}
		resp := container.HandleListRamps(newTenantRequest())
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("returns error", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{
			ListRamps: &mockListRamps{fn: func(_ context.Context, _ cataloguc.ListRampsRequest) ([]entities.Ramp, error) {
				return nil, errNotFound
			}},
		}
		resp := container.HandleListRamps(newTenantRequest())
		if resp.Status != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.Status)
		}
	})
}

func TestHandleGetRamp(t *testing.T) {
	t.Run("returns ramp by id", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{
			GetRamp: &mockGetRamp{fn: func(_ context.Context, req cataloguc.GetRampRequest) (*entities.Ramp, error) {
				if req.RampID != 5 {
					t.Errorf("expected ramp_id 5, got %d", req.RampID)
				}
				return &entities.Ramp{ID: 5, Name: "Sensorial"}, nil
			}},
		}
		req := newTenantRequest()
		req.Params["id"] = "5"
		resp := container.HandleGetRamp(req)
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("rejects invalid id", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{}
		req := newTenantRequest()
		req.Params["id"] = "abc"
		resp := container.HandleGetRamp(req)
		if resp.Status == 0 || resp.Status == http.StatusOK {
			t.Errorf("expected error status, got %d", resp.Status)
		}
	})
}

func TestHandleListDevices(t *testing.T) {
	t.Run("returns devices", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{
			ListDevices: &mockListDevices{fn: func(_ context.Context, _ cataloguc.ListDevicesRequest) ([]entities.Device, error) {
				return []entities.Device{{ID: 1, Name: "Timer"}}, nil
			}},
		}
		resp := container.HandleListDevices(newTenantRequest())
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("filters by ramp_id", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{
			ListDevices: &mockListDevices{fn: func(_ context.Context, req cataloguc.ListDevicesRequest) ([]entities.Device, error) {
				if req.RampID == nil || *req.RampID != 3 {
					t.Error("expected ramp_id 3")
				}
				return []entities.Device{}, nil
			}},
		}
		req := newTenantRequest()
		req.Queries["ramp_id"] = "3"
		resp := container.HandleListDevices(req)
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("rejects invalid ramp_id", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{}
		req := newTenantRequest()
		req.Queries["ramp_id"] = "xyz"
		resp := container.HandleListDevices(req)
		if resp.Status == 0 || resp.Status == http.StatusOK {
			t.Errorf("expected error status, got %d", resp.Status)
		}
	})
}

func TestHandleGetDevice(t *testing.T) {
	t.Run("returns device", func(t *testing.T) {
		container := &entrypoints.CatalogContainer{
			GetDevice: &mockGetDevice{fn: func(_ context.Context, req cataloguc.GetDeviceRequest) (*entities.Device, error) {
				return &entities.Device{ID: req.DeviceID, Name: "Pictogramas"}, nil
			}},
		}
		req := newTenantRequest()
		req.Params["id"] = "7"
		resp := container.HandleGetDevice(req)
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})
}

// ==================== Management Handler Tests ====================

func TestHandleListClassrooms(t *testing.T) {
	t.Run("returns classrooms", func(t *testing.T) {
		container := &entrypoints.ManagementContainer{
			ListClassrooms: &mockListClassrooms{fn: func(_ context.Context, _ mgmtuc.ListClassroomsRequest) ([]entities.Classroom, error) {
				return []entities.Classroom{{ID: 1, Name: "3ro A"}}, nil
			}},
		}
		resp := container.HandleListClassrooms(newTenantRequest())
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})
}

func TestHandleGetClassroom(t *testing.T) {
	t.Run("returns classroom", func(t *testing.T) {
		container := &entrypoints.ManagementContainer{
			GetClassroom: &mockGetClassroom{fn: func(_ context.Context, req mgmtuc.GetClassroomRequest) (*entities.Classroom, error) {
				return &entities.Classroom{ID: req.ClassroomID, Name: "4to B"}, nil
			}},
		}
		req := newTenantRequest()
		req.Params["id"] = "2"
		resp := container.HandleGetClassroom(req)
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("returns 404 for not found", func(t *testing.T) {
		container := &entrypoints.ManagementContainer{
			GetClassroom: &mockGetClassroom{fn: func(_ context.Context, _ mgmtuc.GetClassroomRequest) (*entities.Classroom, error) {
				return nil, errNotFound
			}},
		}
		req := newTenantRequest()
		req.Params["id"] = "999"
		resp := container.HandleGetClassroom(req)
		if resp.Status != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.Status)
		}
	})
}

func TestHandleCreateClassroom(t *testing.T) {
	t.Run("creates classroom", func(t *testing.T) {
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
		if resp.Status != http.StatusCreated {
			t.Errorf("expected 201, got %d", resp.Status)
		}
	})
}

func TestHandleDeleteClassroom(t *testing.T) {
	t.Run("deletes classroom", func(t *testing.T) {
		container := &entrypoints.ManagementContainer{
			DeleteClassroom: &mockDeleteClassroom{fn: func(_ context.Context, req mgmtuc.DeleteClassroomRequest) error {
				if req.ClassroomID != 3 {
					t.Errorf("expected classroom_id 3, got %d", req.ClassroomID)
				}
				return nil
			}},
		}
		req := newTenantRequest()
		req.Params["id"] = "3"
		resp := container.HandleDeleteClassroom(req)
		if resp.Status != http.StatusNoContent {
			t.Errorf("expected 204, got %d", resp.Status)
		}
	})
}

func TestHandleListTeachers(t *testing.T) {
	t.Run("returns teachers", func(t *testing.T) {
		container := &entrypoints.ManagementContainer{
			ListTeachers: &mockListTeachers{fn: func(_ context.Context, _ mgmtuc.ListTeachersRequest) ([]entities.User, error) {
				return []entities.User{{ID: 1, Name: "Ana", Email: "ana@test.com", Role: "teacher"}}, nil
			}},
		}
		resp := container.HandleListTeachers(newTenantRequest())
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})
}

// ==================== Auth Handler Tests ====================

func TestHandleGetMe(t *testing.T) {
	t.Run("returns current user", func(t *testing.T) {
		container := &entrypoints.AuthContainer{
			GetMe: &mockGetMe{fn: func(_ context.Context, req authuc.GetMeRequest) (*entities.User, error) {
				if req.UserID != 42 {
					t.Errorf("expected user_id 42, got %d", req.UserID)
				}
				return &entities.User{ID: 42, Name: "Test", Email: "test@test.com", Role: "teacher"}, nil
			}},
		}
		resp := container.HandleGetMe(newTenantRequest())
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		container := &entrypoints.AuthContainer{
			GetMe: &mockGetMe{fn: func(_ context.Context, _ authuc.GetMeRequest) (*entities.User, error) {
				return nil, errNotFound
			}},
		}
		resp := container.HandleGetMe(newTenantRequest())
		if resp.Status != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.Status)
		}
	})
}

// ==================== Dashboard Handler Tests ====================

func TestHandleGetMetrics(t *testing.T) {
	t.Run("returns metrics", func(t *testing.T) {
		container := &entrypoints.DashboardContainer{
			GetMetrics: &mockGetMetrics{fn: func(_ context.Context, _ dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error) {
				return &dashuc.GetMetricsResponse{
					TotalStudents:        10,
					StudentsWithProfiles: 5,
				}, nil
			}},
		}
		resp := container.HandleGetMetrics(newTenantRequest())
		if resp.Status != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Status)
		}
	})

	t.Run("returns error", func(t *testing.T) {
		container := &entrypoints.DashboardContainer{
			GetMetrics: &mockGetMetrics{fn: func(_ context.Context, _ dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error) {
				return nil, errBadRequest
			}},
		}
		resp := container.HandleGetMetrics(newTenantRequest())
		if resp.Status != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.Status)
		}
	})
}
