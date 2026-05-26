package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// --- RampProvider ---

type MockRampProvider struct {
	ListRampsFn func(ctx context.Context, orgID uuid.UUID) ([]entities.Ramp, error)
	GetRampFn   func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error)
}

func (m *MockRampProvider) ListRamps(ctx context.Context, orgID uuid.UUID) ([]entities.Ramp, error) {
	return m.ListRampsFn(ctx, orgID)
}

func (m *MockRampProvider) GetRamp(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error) {
	return m.GetRampFn(ctx, orgID, id)
}

// --- DeviceProvider ---

type MockDeviceProvider struct {
	ListDevicesFn func(ctx context.Context, orgID uuid.UUID, rampID *int64) ([]entities.Device, error)
	GetDeviceFn   func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Device, error)
}

func (m *MockDeviceProvider) ListDevices(ctx context.Context, orgID uuid.UUID, rampID *int64) ([]entities.Device, error) {
	return m.ListDevicesFn(ctx, orgID, rampID)
}

func (m *MockDeviceProvider) GetDevice(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Device, error) {
	return m.GetDeviceFn(ctx, orgID, id)
}

// --- StudentProfileProvider ---

type MockStudentProfileProvider struct {
	GetByStudentIDFn func(ctx context.Context, studentID int64) (*entities.StudentProfile, error)
	UpsertFn         func(ctx context.Context, profile *entities.StudentProfile) error
}

func (m *MockStudentProfileProvider) GetByStudentID(ctx context.Context, studentID int64) (*entities.StudentProfile, error) {
	return m.GetByStudentIDFn(ctx, studentID)
}

func (m *MockStudentProfileProvider) Upsert(ctx context.Context, profile *entities.StudentProfile) error {
	return m.UpsertFn(ctx, profile)
}

// --- StudentProvider ---

type MockStudentProvider struct {
	GetStudentFn      func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Student, error)
	ListByClassroomFn func(ctx context.Context, orgID uuid.UUID, classroomID int64) ([]entities.Student, error)
	ListFn            func(ctx context.Context, orgID uuid.UUID) ([]entities.Student, error)
	CreateFn          func(ctx context.Context, student *entities.Student) error
	UpdateFn          func(ctx context.Context, student *entities.Student) error
	DeleteFn          func(ctx context.Context, orgID uuid.UUID, id int64) error
}

func (m *MockStudentProvider) GetStudent(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Student, error) {
	return m.GetStudentFn(ctx, orgID, id)
}

func (m *MockStudentProvider) ListByClassroom(ctx context.Context, orgID uuid.UUID, classroomID int64) ([]entities.Student, error) {
	return m.ListByClassroomFn(ctx, orgID, classroomID)
}

func (m *MockStudentProvider) List(ctx context.Context, orgID uuid.UUID) ([]entities.Student, error) {
	return m.ListFn(ctx, orgID)
}

func (m *MockStudentProvider) Create(ctx context.Context, student *entities.Student) error {
	return m.CreateFn(ctx, student)
}

func (m *MockStudentProvider) Update(ctx context.Context, student *entities.Student) error {
	return m.UpdateFn(ctx, student)
}

func (m *MockStudentProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	return m.DeleteFn(ctx, orgID, id)
}

// --- AIClient ---

type MockAIClient struct {
	GenerateFn       func(ctx context.Context, params providers.GenerateParams) (string, error)
	ChatFn           func(ctx context.Context, messages []providers.ChatMessage) (*providers.ChatResponse, error)
	ChatWithToolsFn  func(ctx context.Context, messages []providers.ChatMessage, tools []providers.ToolDefinition) (*providers.ChatResponse, error)
}

func (m *MockAIClient) Generate(ctx context.Context, params providers.GenerateParams) (string, error) {
	return m.GenerateFn(ctx, params)
}

func (m *MockAIClient) Chat(ctx context.Context, messages []providers.ChatMessage) (*providers.ChatResponse, error) {
	return m.ChatFn(ctx, messages)
}

func (m *MockAIClient) ChatWithTools(ctx context.Context, messages []providers.ChatMessage, tools []providers.ToolDefinition) (*providers.ChatResponse, error) {
	return m.ChatWithToolsFn(ctx, messages, tools)
}

// --- AdaptationProvider ---

type MockAdaptationProvider struct {
	ListFn        func(ctx context.Context, orgID uuid.UUID, studentID *int64) ([]entities.Adaptation, error)
	GetFn         func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error)
	CreateFn      func(ctx context.Context, adaptation *entities.Adaptation) error
	UpdateFn      func(ctx context.Context, adaptation *entities.Adaptation) error
	DeleteFn      func(ctx context.Context, orgID uuid.UUID, id int64) error
	SetDevicesFn  func(ctx context.Context, adaptationID int64, deviceIDs []int64) error
	CountSinceFn  func(ctx context.Context, orgID uuid.UUID, since time.Time) (int, error)
	TopDevicesFn  func(ctx context.Context, orgID uuid.UUID, limit int) ([]providers.DeviceUsageStat, error)
}

func (m *MockAdaptationProvider) List(ctx context.Context, orgID uuid.UUID, studentID *int64) ([]entities.Adaptation, error) {
	return m.ListFn(ctx, orgID, studentID)
}

func (m *MockAdaptationProvider) Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
	return m.GetFn(ctx, orgID, id)
}

func (m *MockAdaptationProvider) Create(ctx context.Context, adaptation *entities.Adaptation) error {
	return m.CreateFn(ctx, adaptation)
}

func (m *MockAdaptationProvider) Update(ctx context.Context, adaptation *entities.Adaptation) error {
	return m.UpdateFn(ctx, adaptation)
}

func (m *MockAdaptationProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	return m.DeleteFn(ctx, orgID, id)
}

func (m *MockAdaptationProvider) SetDevices(ctx context.Context, adaptationID int64, deviceIDs []int64) error {
	return m.SetDevicesFn(ctx, adaptationID, deviceIDs)
}

func (m *MockAdaptationProvider) CountSince(ctx context.Context, orgID uuid.UUID, since time.Time) (int, error) {
	return m.CountSinceFn(ctx, orgID, since)
}

func (m *MockAdaptationProvider) TopDevices(ctx context.Context, orgID uuid.UUID, limit int) ([]providers.DeviceUsageStat, error) {
	return m.TopDevicesFn(ctx, orgID, limit)
}

// --- AdaptationResourceProvider ---

type MockAdaptationResourceProvider struct {
	ListByAdaptationFn func(ctx context.Context, adaptationID int64) ([]entities.AdaptationResource, error)
}

func (m *MockAdaptationResourceProvider) ListByAdaptation(ctx context.Context, adaptationID int64) ([]entities.AdaptationResource, error) {
	return m.ListByAdaptationFn(ctx, adaptationID)
}

// --- UserProvider ---

type MockUserProvider struct {
	GetByIDFn    func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error)
	ListByRoleFn func(ctx context.Context, orgID uuid.UUID, role string) ([]entities.User, error)
}

func (m *MockUserProvider) GetByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error) {
	return m.GetByIDFn(ctx, orgID, id)
}

func (m *MockUserProvider) ListByRole(ctx context.Context, orgID uuid.UUID, role string) ([]entities.User, error) {
	return m.ListByRoleFn(ctx, orgID, role)
}

// --- ClassroomProvider ---

type MockClassroomProvider struct {
	ListFn   func(ctx context.Context, orgID uuid.UUID) ([]entities.Classroom, error)
	GetFn    func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error)
	CreateFn func(ctx context.Context, classroom *entities.Classroom) error
	UpdateFn func(ctx context.Context, classroom *entities.Classroom) error
	DeleteFn func(ctx context.Context, orgID uuid.UUID, id int64) error
}

func (m *MockClassroomProvider) List(ctx context.Context, orgID uuid.UUID) ([]entities.Classroom, error) {
	return m.ListFn(ctx, orgID)
}

func (m *MockClassroomProvider) Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error) {
	return m.GetFn(ctx, orgID, id)
}

func (m *MockClassroomProvider) Create(ctx context.Context, classroom *entities.Classroom) error {
	return m.CreateFn(ctx, classroom)
}

func (m *MockClassroomProvider) Update(ctx context.Context, classroom *entities.Classroom) error {
	return m.UpdateFn(ctx, classroom)
}

func (m *MockClassroomProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	return m.DeleteFn(ctx, orgID, id)
}

// --- ConversationProvider ---

type MockConversationProvider struct {
	ListByUserFn  func(ctx context.Context, orgID uuid.UUID, userID int64, mode string) ([]entities.Conversation, error)
	AppendTurnFn  func(ctx context.Context, params providers.AppendTurnParams) (int64, error)
}

func (m *MockConversationProvider) ListByUser(ctx context.Context, orgID uuid.UUID, userID int64, mode string) ([]entities.Conversation, error) {
	return m.ListByUserFn(ctx, orgID, userID, mode)
}

func (m *MockConversationProvider) AppendTurn(ctx context.Context, params providers.AppendTurnParams) (int64, error) {
	if m.AppendTurnFn == nil {
		return params.ConversationID, nil
	}
	return m.AppendTurnFn(ctx, params)
}

// --- AIUsageProvider ---

type MockAIUsageProvider struct {
	RecordFn    func(ctx context.Context, record providers.AIUsageRecord) error
	SummarizeFn func(ctx context.Context, orgID uuid.UUID, since time.Time) (*providers.AIUsageSummary, error)
}

func (m *MockAIUsageProvider) Record(ctx context.Context, record providers.AIUsageRecord) error {
	if m.RecordFn == nil {
		return nil
	}
	return m.RecordFn(ctx, record)
}

func (m *MockAIUsageProvider) Summarize(ctx context.Context, orgID uuid.UUID, since time.Time) (*providers.AIUsageSummary, error) {
	if m.SummarizeFn == nil {
		return &providers.AIUsageSummary{}, nil
	}
	return m.SummarizeFn(ctx, orgID, since)
}
