package management_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestDeleteClassroom_DeletesClassroom(t *testing.T) {
	ctx := context.Background()
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Delete", ctx, testutil.TestOrgID, int64(1)).Return(nil)

	uc := management.NewDeleteClassroom(classrooms)
	err := uc.Execute(ctx, management.DeleteClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 1,
	})

	assert.NoError(t, err)
	classrooms.AssertExpectations(t)
}

func TestDeleteClassroom_RejectsNilOrgID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewDeleteClassroom(classrooms)

	err := uc.Execute(context.Background(), management.DeleteClassroomRequest{
		OrgID:       uuid.Nil,
		ClassroomID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteClassroom_RejectsZeroClassroomID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewDeleteClassroom(classrooms)

	err := uc.Execute(context.Background(), management.DeleteClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteClassroom_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Delete", ctx, testutil.TestOrgID, int64(99)).Return(errClassroomNotFound)

	uc := management.NewDeleteClassroom(classrooms)
	err := uc.Execute(ctx, management.DeleteClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 99,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	classrooms.AssertExpectations(t)
}
