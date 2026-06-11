package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

type openMocks struct {
	students  *mockproviders.MockStudentProvider
	summaries *mockproviders.MockConversationSummaryProvider
}

func newOpenMocks() openMocks {
	return openMocks{
		students:  new(mockproviders.MockStudentProvider),
		summaries: new(mockproviders.MockConversationSummaryProvider),
	}
}

func (m openMocks) usecase() inclusion.OpenSession {
	return inclusion.NewOpenSession(m.students, m.summaries)
}

func ptrInt64(v int64) *int64 { return &v }

func TestOpenSession_RejectsNilOrgID(t *testing.T) {
	m := newOpenMocks()

	_, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{UserID: 1})

	assert.Error(t, err)
	m.students.AssertNotCalled(t, "GetStudent")
}

func TestOpenSession_NoDimensionAsksWhichOne(t *testing.T) {
	m := newOpenMocks()

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 1,
	})

	require.NoError(t, err)
	assert.True(t, got.NeedsDimension)
	assert.Contains(t, got.Greeting, "alumno")
	m.students.AssertNotCalled(t, "GetStudent")
}

func TestOpenSession_AmbiguousDimensionReasks(t *testing.T) {
	m := newOpenMocks()

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: "cualquier cosa",
	})

	require.NoError(t, err)
	assert.True(t, got.NeedsDimension)
	assert.Empty(t, got.Dimension)
}

func TestOpenSession_StudentWithoutIDReasks(t *testing.T) {
	m := newOpenMocks()

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionStudent,
	})

	require.NoError(t, err)
	assert.True(t, got.NeedsDimension)
	m.students.AssertNotCalled(t, "GetStudent")
}

func TestOpenSession_StudentLoadsContextAndPriorSummaries(t *testing.T) {
	m := newOpenMocks()
	student := testutil.NewStudentWithProfile(7, 1, "Pedro", []string{"se_distrae"})
	prior := []entities.ConversationSummary{
		{ConversationID: 99, Summary: "Última charla sobre Pedro", TopicKeywords: []string{"TDAH"}},
	}
	m.students.On("GetStudent", mock.Anything, testutil.TestOrgID, int64(7)).Return(&student, nil)
	m.summaries.On("RecentByStudent", mock.Anything, testutil.TestOrgID, int64(7), 10).Return(prior, nil)

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionStudent,
		StudentID: ptrInt64(7),
	})

	require.NoError(t, err)
	assert.False(t, got.NeedsDimension)
	assert.Equal(t, inclusion.DimensionStudent, got.Dimension)
	require.NotNil(t, got.Student)
	assert.Equal(t, "Pedro", got.Student.Name)
	assert.Contains(t, got.Greeting, "Pedro")
	require.Len(t, got.PriorSummaries, 1)
	assert.Equal(t, int64(99), got.PriorSummaries[0].ConversationID)
	m.students.AssertExpectations(t)
	m.summaries.AssertExpectations(t)
}

func TestOpenSession_TopicWithoutKeywordReasks(t *testing.T) {
	m := newOpenMocks()

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionTopic,
	})

	require.NoError(t, err)
	assert.True(t, got.NeedsDimension)
	m.summaries.AssertNotCalled(t, "RecentByTopic")
}

func TestOpenSession_TopicRetrievesByKeyword(t *testing.T) {
	m := newOpenMocks()
	m.summaries.On("RecentByTopic", mock.Anything, testutil.TestOrgID, "TEA", 10).
		Return([]entities.ConversationSummary{{ConversationID: 5}}, nil)

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionTopic,
		Topic:     "TEA",
	})

	require.NoError(t, err)
	assert.Equal(t, inclusion.DimensionTopic, got.Dimension)
	require.Len(t, got.PriorSummaries, 1)
	m.summaries.AssertExpectations(t)
}

func TestOpenSession_ToolkitDoesNotLoadStudent(t *testing.T) {
	m := newOpenMocks()

	got, err := m.usecase().Execute(context.Background(), inclusion.OpenSessionRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionToolkit,
	})

	require.NoError(t, err)
	assert.Equal(t, inclusion.DimensionToolkit, got.Dimension)
	assert.False(t, got.NeedsDimension)
	m.students.AssertNotCalled(t, "GetStudent")
	m.summaries.AssertNotCalled(t, "RecentByDevice")
}
