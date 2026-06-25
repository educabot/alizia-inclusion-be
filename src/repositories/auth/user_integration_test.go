//go:build integration

package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/repositories/auth"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

func TestUserRepo_GetByID(t *testing.T) {
	tx := pgtest.Tx(t)
	u := entities.User{OrganizationID: testutil.TestOrgID, Email: "ana@test.com", Name: "Ana", Role: "teacher"}
	require.NoError(t, tx.Create(&u).Error)

	repo := auth.NewUserRepo(tx)
	got, err := repo.GetByID(context.Background(), testutil.TestOrgID, u.ID)

	require.NoError(t, err)
	assert.Equal(t, "Ana", got.Name)
	assert.Equal(t, "teacher", got.Role)
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := auth.NewUserRepo(tx)

	_, err := repo.GetByID(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUserRepo_GetByID_NotFoundAcrossTenant(t *testing.T) {
	tx := pgtest.Tx(t)
	u := entities.User{OrganizationID: pgtest.OtherOrgID, Email: "x@test.com", Name: "X", Role: "teacher"}
	require.NoError(t, tx.Create(&u).Error)

	repo := auth.NewUserRepo(tx)
	_, err := repo.GetByID(context.Background(), testutil.TestOrgID, u.ID)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUserRepo_ListByRole(t *testing.T) {
	tx := pgtest.Tx(t)
	teacher1 := entities.User{OrganizationID: testutil.TestOrgID, Email: "ana@test.com", Name: "Ana", Role: "teacher"}
	teacher2 := entities.User{OrganizationID: testutil.TestOrgID, Email: "beto@test.com", Name: "Beto", Role: "teacher"}
	admin := entities.User{OrganizationID: testutil.TestOrgID, Email: "coord@test.com", Name: "Coord", Role: "coordinator"}
	require.NoError(t, tx.Create(&teacher1).Error)
	require.NoError(t, tx.Create(&teacher2).Error)
	require.NoError(t, tx.Create(&admin).Error)

	repo := auth.NewUserRepo(tx)
	got, err := repo.ListByRole(context.Background(), testutil.TestOrgID, "teacher")

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Ana", got[0].Name, "ordered by name ASC")
	assert.Equal(t, "Beto", got[1].Name)
}
