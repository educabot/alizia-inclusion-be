//go:build integration

package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

func TestStudentRepo_CreateGetWithProfilePreloaded(t *testing.T) {
	tx := pgtest.Tx(t)
	classroomID := seedClassroom(t, tx, "3ro A")
	repo := inclusion.NewStudentRepo(tx)
	s := entities.Student{OrganizationID: testutil.TestOrgID, ClassroomID: classroomID, Name: "Lucía"}
	require.NoError(t, repo.Create(context.Background(), &s))
	require.NoError(t, tx.Create(&entities.StudentProfile{StudentID: s.ID, Difficulties: []string{"se_distrae"}}).Error)

	got, err := repo.GetStudent(context.Background(), testutil.TestOrgID, s.ID)

	require.NoError(t, err)
	assert.Equal(t, "Lucía", got.Name)
	require.NotNil(t, got.Profile, "Profile is preloaded")
	assert.Equal(t, []string(got.Profile.Difficulties), []string{"se_distrae"})
}

func TestStudentRepo_GetStudent_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewStudentRepo(tx)

	_, err := repo.GetStudent(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestStudentRepo_ListByClassroom_OrdersByNameAndScopes(t *testing.T) {
	tx := pgtest.Tx(t)
	classroomID := seedClassroom(t, tx, "3ro A")
	other := seedClassroom(t, tx, "Otra")
	repo := inclusion.NewStudentRepo(tx)
	for _, n := range []string{"Zoe", "Ana"} {
		require.NoError(t, repo.Create(context.Background(), &entities.Student{OrganizationID: testutil.TestOrgID, ClassroomID: classroomID, Name: n}))
	}
	require.NoError(t, repo.Create(context.Background(), &entities.Student{OrganizationID: testutil.TestOrgID, ClassroomID: other, Name: "Fuera"}))

	got, err := repo.ListByClassroom(context.Background(), testutil.TestOrgID, classroomID)

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Ana", got[0].Name, "ordered by name ASC")
	assert.Equal(t, "Zoe", got[1].Name)
}

func TestStudentRepo_Update(t *testing.T) {
	tx := pgtest.Tx(t)
	id := seedStudent(t, tx, "Original")
	repo := inclusion.NewStudentRepo(tx)
	got, err := repo.GetStudent(context.Background(), testutil.TestOrgID, id)
	require.NoError(t, err)

	got.Name = "Cambiado"
	require.NoError(t, repo.Update(context.Background(), got))

	reloaded, err := repo.GetStudent(context.Background(), testutil.TestOrgID, id)
	require.NoError(t, err)
	assert.Equal(t, "Cambiado", reloaded.Name)
}

func TestStudentRepo_Delete(t *testing.T) {
	tx := pgtest.Tx(t)
	id := seedStudent(t, tx, "Borrar")
	repo := inclusion.NewStudentRepo(tx)

	require.NoError(t, repo.Delete(context.Background(), testutil.TestOrgID, id))

	_, err := repo.GetStudent(context.Background(), testutil.TestOrgID, id)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestStudentRepo_Delete_NotFoundWhenMissing(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewStudentRepo(tx)

	err := repo.Delete(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}
