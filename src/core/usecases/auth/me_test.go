package auth_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetMe_ReturnsUser(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := auth.NewGetMe(users)
	ctx := context.Background()
	want := testutil.NewUser(42, "Ana")
	users.On("GetByID", ctx, testutil.TestOrgID, int64(42)).Return(&want, nil)

	got, err := uc.Execute(ctx, auth.GetMeRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 42,
	})

	assert.NoError(t, err)
	assert.Equal(t, &want, got)
	users.AssertExpectations(t)
}

func TestGetMe_RejectsNilOrgID(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := auth.NewGetMe(users)

	_, err := uc.Execute(context.Background(), auth.GetMeRequest{
		OrgID:  uuid.Nil,
		UserID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	users.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetMe_RejectsZeroUserID(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := auth.NewGetMe(users)

	_, err := uc.Execute(context.Background(), auth.GetMeRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	users.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetMe_ReturnsNotFound(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := auth.NewGetMe(users)
	ctx := context.Background()
	users.On("GetByID", ctx, testutil.TestOrgID, int64(99)).Return(nil, errUserNotFound)

	got, err := uc.Execute(ctx, auth.GetMeRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 99,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	assert.Nil(t, got)
	users.AssertExpectations(t)
}
