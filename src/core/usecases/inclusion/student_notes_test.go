package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestCreateStudentNote_CreatesWithDefaults(t *testing.T) {
	notes := new(mockproviders.MockStudentNoteProvider)
	ctx := context.Background()
	var captured *entities.StudentNote
	notes.On("Create", ctx, mock.AnythingOfType("*entities.StudentNote")).
		Run(func(args mock.Arguments) {
			n, ok := args.Get(1).(*entities.StudentNote)
			require.True(t, ok)
			n.ID = 1
			captured = n
		}).
		Return(nil)

	result, err := inclusion.NewCreateStudentNote(notes).Execute(ctx, inclusion.CreateStudentNoteRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 5,
		UserID:    7,
		Content:   "le tiembla la mano al escribir",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "seguimiento", captured.Type) // default
	assert.True(t, captured.Internal)             // default
	assert.Equal(t, int64(5), captured.StudentID)
	assert.Equal(t, int64(7), captured.UserID) // dueño de la nota
	notes.AssertExpectations(t)
}

func TestCreateStudentNote_RejectsEmptyContent(t *testing.T) {
	notes := new(mockproviders.MockStudentNoteProvider)
	_, err := inclusion.NewCreateStudentNote(notes).Execute(context.Background(), inclusion.CreateStudentNoteRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 5,
		Content:   "",
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	notes.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateStudentNote_RejectsInvalidType(t *testing.T) {
	notes := new(mockproviders.MockStudentNoteProvider)
	_, err := inclusion.NewCreateStudentNote(notes).Execute(context.Background(), inclusion.CreateStudentNoteRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 5,
		Content:   "x",
		Type:      "no_existe",
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestListStudentNotes_ReturnsNotes(t *testing.T) {
	notes := new(mockproviders.MockStudentNoteProvider)
	ctx := context.Background()
	want := []entities.StudentNote{{ID: 1, StudentID: 5, Content: "n1"}}
	notes.On("ListByStudent", ctx, testutil.TestOrgID, int64(5), int64(7)).Return(want, nil)

	got, err := inclusion.NewListStudentNotes(notes).Execute(ctx, inclusion.ListStudentNotesRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 5,
		UserID:    7,
	})
	require.NoError(t, err)
	assert.Equal(t, want, got)
	notes.AssertExpectations(t)
}

func TestListStudentNotes_RejectsBadStudent(t *testing.T) {
	notes := new(mockproviders.MockStudentNoteProvider)
	_, err := inclusion.NewListStudentNotes(notes).Execute(context.Background(), inclusion.ListStudentNotesRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 0,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}
