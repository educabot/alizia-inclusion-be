//go:build integration

package inclusion_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

// seedAdaptation inserts an adaptation for the given student/teacher with the given
// subject and returns it.
func seedAdaptation(t *testing.T, tx *gorm.DB, studentID, teacherID int64, subject string) entities.Adaptation {
	t.Helper()
	a := entities.Adaptation{
		OrganizationID: testutil.TestOrgID,
		StudentID:      studentID,
		TeacherID:      teacherID,
		Subject:        subject,
		Status:         "en_curso",
	}
	require.NoError(t, tx.Create(&a).Error)
	return a
}

func TestAdaptationRepo_CreateGetDelete(t *testing.T) {
	tx := pgtest.Tx(t)
	studentID := seedStudent(t, tx, "Lucía")
	teacherID := seedUser(t, tx, "Seño", "teacher")
	repo := inclusion.NewAdaptationRepo(tx)

	a := entities.Adaptation{OrganizationID: testutil.TestOrgID, StudentID: studentID, TeacherID: teacherID, Subject: "Matemática", Status: "en_curso"}
	require.NoError(t, repo.Create(context.Background(), &a))
	require.NotZero(t, a.ID)

	got, err := repo.Get(context.Background(), testutil.TestOrgID, a.ID)
	require.NoError(t, err)
	assert.Equal(t, "Matemática", got.Subject)
	require.NotNil(t, got.Student, "Student preloaded")
	assert.Equal(t, "Lucía", got.Student.Name)

	require.NoError(t, repo.Delete(context.Background(), testutil.TestOrgID, a.ID))
	_, err = repo.Get(context.Background(), testutil.TestOrgID, a.ID)
	assert.ErrorIs(t, err, providers.ErrAdaptationNotFound)
}

func TestAdaptationRepo_Get_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewAdaptationRepo(tx)

	_, err := repo.Get(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrAdaptationNotFound)
}

func TestAdaptationRepo_List_FiltersByTeacherAndStudent(t *testing.T) {
	tx := pgtest.Tx(t)
	t1 := seedUser(t, tx, "T1", "teacher")
	t2 := seedUser(t, tx, "T2", "teacher")
	s1 := seedStudent(t, tx, "S1")
	s2 := seedStudent(t, tx, "S2")
	seedAdaptation(t, tx, s1, t1, "Mate")
	seedAdaptation(t, tx, s2, t1, "Lengua")
	seedAdaptation(t, tx, s1, t2, "Ciencias")
	repo := inclusion.NewAdaptationRepo(tx)

	byTeacher, err := repo.List(context.Background(), testutil.TestOrgID, providers.AdaptationFilter{TeacherID: &t1})
	require.NoError(t, err)
	assert.Len(t, byTeacher, 2)

	byStudent, err := repo.List(context.Background(), testutil.TestOrgID, providers.AdaptationFilter{StudentID: &s1})
	require.NoError(t, err)
	assert.Len(t, byStudent, 2)
}

func TestAdaptationRepo_List_QueryIsCaseInsensitive(t *testing.T) {
	tx := pgtest.Tx(t)
	teacherID := seedUser(t, tx, "T", "teacher")
	studentID := seedStudent(t, tx, "S")
	seedAdaptation(t, tx, studentID, teacherID, "Matemática")
	seedAdaptation(t, tx, studentID, teacherID, "Lengua")
	repo := inclusion.NewAdaptationRepo(tx)

	// ILIKE: lowercase query must match the capitalized subject (Postgres-only behaviour).
	got, err := repo.List(context.Background(), testutil.TestOrgID, providers.AdaptationFilter{Query: "matem"})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Matemática", got[0].Subject)
}

func TestAdaptationRepo_SetDevicesAndPreload(t *testing.T) {
	tx := pgtest.Tx(t)
	teacherID := seedUser(t, tx, "T", "teacher")
	studentID := seedStudent(t, tx, "S")
	ramp := testutil.NewRamp(0, "Sensorial")
	require.NoError(t, tx.Create(&ramp).Error)
	d1 := testutil.NewDevice(0, ramp.ID, "Timer")
	d2 := testutil.NewDevice(0, ramp.ID, "Pictogramas")
	require.NoError(t, tx.Create(&d1).Error)
	require.NoError(t, tx.Create(&d2).Error)
	a := seedAdaptation(t, tx, studentID, teacherID, "Mate")
	repo := inclusion.NewAdaptationRepo(tx)

	require.NoError(t, repo.SetDevices(context.Background(), a.ID, []int64{d1.ID, d2.ID}))

	got, err := repo.Get(context.Background(), testutil.TestOrgID, a.ID)
	require.NoError(t, err)
	assert.Len(t, got.Devices, 2, "m2m devices preloaded")
}

func TestAdaptationRepo_CountSinceAndTopDevices(t *testing.T) {
	tx := pgtest.Tx(t)
	teacherID := seedUser(t, tx, "T", "teacher")
	studentID := seedStudent(t, tx, "S")
	ramp := testutil.NewRamp(0, "Sensorial")
	require.NoError(t, tx.Create(&ramp).Error)
	dev := testutil.NewDevice(0, ramp.ID, "Timer")
	require.NoError(t, tx.Create(&dev).Error)
	a := seedAdaptation(t, tx, studentID, teacherID, "Mate")
	repo := inclusion.NewAdaptationRepo(tx)
	require.NoError(t, repo.SetDevices(context.Background(), a.ID, []int64{dev.ID}))

	count, err := repo.CountSince(context.Background(), testutil.TestOrgID, time.Now().Add(-time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	top, err := repo.TopDevices(context.Background(), testutil.TestOrgID, 5)
	require.NoError(t, err)
	require.Len(t, top, 1)
	assert.Equal(t, "Timer", top[0].DeviceName)
	assert.Equal(t, 1, top[0].Count)
}

func TestAdaptationResourceRepo_ListByAdaptation(t *testing.T) {
	tx := pgtest.Tx(t)
	teacherID := seedUser(t, tx, "T", "teacher")
	studentID := seedStudent(t, tx, "S")
	a := seedAdaptation(t, tx, studentID, teacherID, "Mate")
	require.NoError(t, tx.Create(&entities.AdaptationResource{AdaptationID: a.ID, Title: "Ficha", FileURL: "https://x/f.pdf", FileType: "pdf"}).Error)
	repo := inclusion.NewAdaptationResourceRepo(tx)

	got, err := repo.ListByAdaptation(context.Background(), a.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Ficha", got[0].Title)
}
