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

func TestDiagnosisRepo_ListByStudentProfile_PreloadsCatalog(t *testing.T) {
	tx := pgtest.Tx(t)
	studentID := seedStudent(t, tx, "Mati")
	profileID := seedStudentProfile(t, tx, studentID)
	diag := entities.Diagnosis{Name: "TDAH", Category: testutil.Ptr("neurodesarrollo")} // global (org nil)
	require.NoError(t, tx.Create(&diag).Error)
	require.NoError(t, tx.Create(&entities.StudentDiagnosis{
		StudentProfileID: profileID, DiagnosisID: diag.ID, Severity: testutil.Ptr("leve"),
	}).Error)
	repo := inclusion.NewDiagnosisRepo(tx)

	got, err := repo.ListByStudentProfile(context.Background(), profileID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Diagnosis, "catalog entry preloaded")
	assert.Equal(t, "TDAH", got[0].Diagnosis.Name)
}

func TestPPIRepo_GetByStudentID(t *testing.T) {
	tx := pgtest.Tx(t)
	studentID := seedStudent(t, tx, "Sofi")
	require.NoError(t, tx.Create(&entities.PPI{
		OrganizationID: testutil.TestOrgID, StudentID: studentID, Status: "vigente",
		Objectives: []string{"mejorar lectura"},
	}).Error)
	repo := inclusion.NewPPIRepo(tx)

	got, err := repo.GetByStudentID(context.Background(), testutil.TestOrgID, studentID)

	require.NoError(t, err)
	assert.Equal(t, "vigente", got.Status)
	assert.Equal(t, []string{"mejorar lectura"}, []string(got.Objectives))
}

func TestPPIRepo_GetByStudentID_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewPPIRepo(tx)

	_, err := repo.GetByStudentID(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestSituationRepo_List_GlobalPlusOwnOrgExcludesOthers(t *testing.T) {
	tx := pgtest.Tx(t)
	global := entities.Situation{Code: "no_inicia", Name: "No inicia la tarea", SortOrder: 1} // org nil
	own := entities.Situation{OrganizationID: &testutil.TestOrgID, Code: "propia", Name: "Propia de la org", SortOrder: 2}
	other := entities.Situation{OrganizationID: &pgtest.OtherOrgID, Code: "ajena", Name: "De otra org", SortOrder: 3}
	require.NoError(t, tx.Create(&global).Error)
	require.NoError(t, tx.Create(&own).Error)
	require.NoError(t, tx.Create(&other).Error)
	repo := inclusion.NewSituationRepo(tx)

	got, err := repo.List(context.Background(), testutil.TestOrgID)

	require.NoError(t, err)
	require.Len(t, got, 2, "global + own org, never another org's")
	assert.Equal(t, "no_inicia", got[0].Code, "ordered by sort_order ASC")
	assert.Equal(t, "propia", got[1].Code)
}

func TestIntegradoraAssignmentRepo_ListAndIsAssigned(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Integradora", "teacher")
	s1 := seedStudent(t, tx, "A")
	s2 := seedStudent(t, tx, "B")
	require.NoError(t, tx.Create(&entities.IntegradoraAssignment{OrganizationID: testutil.TestOrgID, UserID: userID, StudentID: s1}).Error)
	require.NoError(t, tx.Create(&entities.IntegradoraAssignment{OrganizationID: testutil.TestOrgID, UserID: userID, StudentID: s2}).Error)
	repo := inclusion.NewIntegradoraAssignmentRepo(tx)

	ids, err := repo.ListStudentIDsByUser(context.Background(), testutil.TestOrgID, userID)
	require.NoError(t, err)
	assert.Equal(t, []int64{s1, s2}, ids)

	assigned, err := repo.IsAssigned(context.Background(), testutil.TestOrgID, userID, s1)
	require.NoError(t, err)
	assert.True(t, assigned)

	notAssigned, err := repo.IsAssigned(context.Background(), testutil.TestOrgID, userID, 999)
	require.NoError(t, err)
	assert.False(t, notAssigned)
}
