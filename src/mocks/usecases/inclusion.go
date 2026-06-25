package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	inclusionuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
)

type MockGetStudentProfile struct{ mock.Mock }

func (m *MockGetStudentProfile) Execute(ctx context.Context, req inclusionuc.GetStudentProfileRequest) (*entities.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Student), args.Error(1)
}

type MockUpsertStudentProfile struct{ mock.Mock }

func (m *MockUpsertStudentProfile) Execute(ctx context.Context, req inclusionuc.UpsertStudentProfileRequest) (*entities.StudentProfile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.StudentProfile), args.Error(1)
}

type MockListClassroomStudents struct{ mock.Mock }

func (m *MockListClassroomStudents) Execute(ctx context.Context, req inclusionuc.ListClassroomStudentsRequest) ([]entities.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Student), args.Error(1)
}

type MockRecommendDevice struct{ mock.Mock }

func (m *MockRecommendDevice) Execute(ctx context.Context, req inclusionuc.RecommendDeviceRequest) (*inclusionuc.RecommendDeviceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.RecommendDeviceResponse), args.Error(1)
}

type MockAssistClassroom struct{ mock.Mock }

func (m *MockAssistClassroom) Execute(ctx context.Context, req inclusionuc.AssistClassroomRequest) (*inclusionuc.AssistClassroomResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.AssistClassroomResponse), args.Error(1)
}

type MockOpenSession struct{ mock.Mock }

func (m *MockOpenSession) Execute(ctx context.Context, req inclusionuc.OpenSessionRequest) (*inclusionuc.OpenSessionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.OpenSessionResponse), args.Error(1)
}

type MockCloseSession struct{ mock.Mock }

func (m *MockCloseSession) Execute(ctx context.Context, req inclusionuc.CloseSessionRequest) (*inclusionuc.CloseSessionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.CloseSessionResponse), args.Error(1)
}

type MockBuildPromptContext struct{ mock.Mock }

func (m *MockBuildPromptContext) Execute(ctx context.Context, req inclusionuc.BuildContextRequest) (*inclusionuc.PromptContext, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.PromptContext), args.Error(1)
}

type MockSearchPedagogicalContent struct{ mock.Mock }

func (m *MockSearchPedagogicalContent) Execute(ctx context.Context, req inclusionuc.SearchContentRequest) (*inclusionuc.SearchContentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.SearchContentResponse), args.Error(1)
}

type MockListStudents struct{ mock.Mock }

func (m *MockListStudents) Execute(ctx context.Context, req inclusionuc.ListStudentsRequest) ([]entities.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Student), args.Error(1)
}

type MockCreateStudent struct{ mock.Mock }

func (m *MockCreateStudent) Execute(ctx context.Context, req inclusionuc.CreateStudentRequest) (*entities.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Student), args.Error(1)
}

type MockUpdateStudent struct{ mock.Mock }

func (m *MockUpdateStudent) Execute(ctx context.Context, req inclusionuc.UpdateStudentRequest) (*entities.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Student), args.Error(1)
}

type MockDeleteStudent struct{ mock.Mock }

func (m *MockDeleteStudent) Execute(ctx context.Context, req inclusionuc.DeleteStudentRequest) error {
	return m.Called(ctx, req).Error(0)
}

type MockListAdaptations struct{ mock.Mock }

func (m *MockListAdaptations) Execute(ctx context.Context, req inclusionuc.ListAdaptationsRequest) ([]entities.Adaptation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Adaptation), args.Error(1)
}

type MockGetAdaptation struct{ mock.Mock }

func (m *MockGetAdaptation) Execute(ctx context.Context, req inclusionuc.GetAdaptationRequest) (*entities.Adaptation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Adaptation), args.Error(1)
}

type MockCreateAdaptation struct{ mock.Mock }

func (m *MockCreateAdaptation) Execute(ctx context.Context, req inclusionuc.CreateAdaptationRequest) (*entities.Adaptation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Adaptation), args.Error(1)
}

type MockUpdateAdaptation struct{ mock.Mock }

func (m *MockUpdateAdaptation) Execute(ctx context.Context, req inclusionuc.UpdateAdaptationRequest) (*entities.Adaptation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Adaptation), args.Error(1)
}

type MockDeleteAdaptation struct{ mock.Mock }

func (m *MockDeleteAdaptation) Execute(ctx context.Context, req inclusionuc.DeleteAdaptationRequest) error {
	return m.Called(ctx, req).Error(0)
}

type MockListAdaptationResources struct{ mock.Mock }

func (m *MockListAdaptationResources) Execute(ctx context.Context, req inclusionuc.ListAdaptationResourcesRequest) ([]entities.AdaptationResource, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.AdaptationResource), args.Error(1)
}

type MockExportAdaptation struct{ mock.Mock }

func (m *MockExportAdaptation) Execute(ctx context.Context, req inclusionuc.ExportAdaptationRequest) (*inclusionuc.ExportedDocument, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inclusionuc.ExportedDocument), args.Error(1)
}

type MockGetChatHistory struct{ mock.Mock }

func (m *MockGetChatHistory) Execute(ctx context.Context, req inclusionuc.GetChatHistoryRequest) ([]entities.Conversation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Conversation), args.Error(1)
}

type MockGetConversation struct{ mock.Mock }

func (m *MockGetConversation) Execute(ctx context.Context, req inclusionuc.GetConversationRequest) (*entities.Conversation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Conversation), args.Error(1)
}
