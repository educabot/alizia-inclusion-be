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

func TestStudentProfileRepo_UpsertInsertsThenUpdates(t *testing.T) {
	tx := pgtest.Tx(t)
	studentID := seedStudent(t, tx, "Mateo")
	repo := inclusion.NewStudentProfileRepo(tx)

	// First upsert inserts.
	p := &entities.StudentProfile{StudentID: studentID, IsTransitory: true, Difficulties: []string{"se_distrae"}}
	require.NoError(t, repo.Upsert(context.Background(), p))

	got, err := repo.GetByStudentID(context.Background(), studentID)
	require.NoError(t, err)
	assert.True(t, got.IsTransitory)
	assert.Equal(t, []string{"se_distrae"}, []string(got.Difficulties))

	// Second upsert on the same student_id updates in place (UNIQUE student_id + ON CONFLICT).
	p2 := &entities.StudentProfile{StudentID: studentID, IsTransitory: false, Difficulties: []string{"no_inicia_tarea"}}
	require.NoError(t, repo.Upsert(context.Background(), p2))

	updated, err := repo.GetByStudentID(context.Background(), studentID)
	require.NoError(t, err)
	assert.False(t, updated.IsTransitory)
	assert.Equal(t, []string{"no_inicia_tarea"}, []string(updated.Difficulties))
}

func TestStudentProfileRepo_GetByStudentID_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewStudentProfileRepo(tx)

	_, err := repo.GetByStudentID(context.Background(), 999)

	assert.ErrorIs(t, err, providers.ErrProfileNotFound)
}

func TestTeacherProfileRepo_GetByUserID(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Docente", "teacher")
	require.NoError(t, tx.Create(&entities.TeacherProfile{
		UserID: userID, OrganizationID: testutil.TestOrgID, Specialization: testutil.Ptr("Primaria"),
	}).Error)
	repo := inclusion.NewTeacherProfileRepo(tx)

	got, err := repo.GetByUserID(context.Background(), testutil.TestOrgID, userID)

	require.NoError(t, err)
	require.NotNil(t, got.Specialization)
	assert.Equal(t, "Primaria", *got.Specialization)
}

func TestTeacherProfileRepo_GetByUserID_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewTeacherProfileRepo(tx)

	_, err := repo.GetByUserID(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}
