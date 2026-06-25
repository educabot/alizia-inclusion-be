//go:build integration

package management_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/repositories/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

func TestClassroomRepo_CreateThenGet(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := management.NewClassroomRepo(tx)
	c := entities.Classroom{OrganizationID: testutil.TestOrgID, Name: "3ro A", Grade: testutil.Ptr("3ro")}

	require.NoError(t, repo.Create(context.Background(), &c))
	require.NotZero(t, c.ID, "Create assigns the BIGSERIAL id")

	got, err := repo.Get(context.Background(), testutil.TestOrgID, c.ID)
	require.NoError(t, err)
	assert.Equal(t, "3ro A", got.Name)
	require.NotNil(t, got.Grade)
	assert.Equal(t, "3ro", *got.Grade)
}

func TestClassroomRepo_List_OrdersByNameAndScopesToOrg(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := management.NewClassroomRepo(tx)
	b := entities.Classroom{OrganizationID: testutil.TestOrgID, Name: "B"}
	a := entities.Classroom{OrganizationID: testutil.TestOrgID, Name: "A"}
	other := entities.Classroom{OrganizationID: pgtest.OtherOrgID, Name: "C"}
	require.NoError(t, repo.Create(context.Background(), &b))
	require.NoError(t, repo.Create(context.Background(), &a))
	require.NoError(t, repo.Create(context.Background(), &other))

	got, err := repo.List(context.Background(), testutil.TestOrgID)

	require.NoError(t, err)
	require.Len(t, got, 2, "other tenant excluded")
	assert.Equal(t, "A", got[0].Name, "ordered by name ASC")
	assert.Equal(t, "B", got[1].Name)
}

func TestClassroomRepo_Update(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := management.NewClassroomRepo(tx)
	c := entities.Classroom{OrganizationID: testutil.TestOrgID, Name: "Original"}
	require.NoError(t, repo.Create(context.Background(), &c))

	c.Name = "Renombrada"
	require.NoError(t, repo.Update(context.Background(), &c))

	got, err := repo.Get(context.Background(), testutil.TestOrgID, c.ID)
	require.NoError(t, err)
	assert.Equal(t, "Renombrada", got.Name)
}

func TestClassroomRepo_Delete(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := management.NewClassroomRepo(tx)
	c := entities.Classroom{OrganizationID: testutil.TestOrgID, Name: "Borrar"}
	require.NoError(t, repo.Create(context.Background(), &c))

	require.NoError(t, repo.Delete(context.Background(), testutil.TestOrgID, c.ID))

	_, err := repo.Get(context.Background(), testutil.TestOrgID, c.ID)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestClassroomRepo_Delete_NotFoundWhenMissing(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := management.NewClassroomRepo(tx)

	err := repo.Delete(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestClassroomRepo_Get_NotFoundAcrossTenant(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := management.NewClassroomRepo(tx)
	c := entities.Classroom{OrganizationID: pgtest.OtherOrgID, Name: "De otra org"}
	require.NoError(t, repo.Create(context.Background(), &c))

	_, err := repo.Get(context.Background(), testutil.TestOrgID, c.ID)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}
