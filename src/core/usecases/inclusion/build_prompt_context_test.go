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

type ctxMocks struct {
	students    *mockproviders.MockStudentProvider
	teachers    *mockproviders.MockTeacherProfileProvider
	situations  *mockproviders.MockSituationCatalogProvider
	diagnoses   *mockproviders.MockDiagnosisProvider
	ppi         *mockproviders.MockPPIProvider
	adaptations *mockproviders.MockAdaptationProvider
	classrooms  *mockproviders.MockClassroomProvider
	devices     *mockproviders.MockDeviceProvider
	summaries   *mockproviders.MockConversationSummaryProvider
}

func newCtxMocks() ctxMocks {
	return ctxMocks{
		students:    new(mockproviders.MockStudentProvider),
		teachers:    new(mockproviders.MockTeacherProfileProvider),
		situations:  new(mockproviders.MockSituationCatalogProvider),
		diagnoses:   new(mockproviders.MockDiagnosisProvider),
		ppi:         new(mockproviders.MockPPIProvider),
		adaptations: new(mockproviders.MockAdaptationProvider),
		classrooms:  new(mockproviders.MockClassroomProvider),
		devices:     new(mockproviders.MockDeviceProvider),
		summaries:   new(mockproviders.MockConversationSummaryProvider),
	}
}

func (m ctxMocks) usecase() inclusion.BuildPromptContext {
	return inclusion.NewBuildPromptContext(
		m.students, m.teachers, m.situations, m.diagnoses, m.ppi,
		m.adaptations, m.classrooms, m.devices, m.summaries,
	)
}

// expectStatic configura el bloque estático (devices + situaciones) que siempre se carga.
func (m ctxMocks) expectStatic() {
	m.devices.On("ListDevices", mock.Anything, testutil.TestOrgID, (*int64)(nil)).
		Return([]entities.Device{testutil.NewDevice(1, 1, "Time Timer")}, nil)
	m.situations.On("List", mock.Anything, testutil.TestOrgID).
		Return([]entities.Situation{{ID: 1, Code: "no_inicia_tarea", Name: "No inicia la tarea"}}, nil)
}

func TestBuildContext_RejectsNilOrgID(t *testing.T) {
	m := newCtxMocks()
	_, err := m.usecase().Execute(context.Background(), inclusion.BuildContextRequest{UserID: 1})
	assert.Error(t, err)
}

func TestBuildContext_LoadsStaticAndTeacher(t *testing.T) {
	m := newCtxMocks()
	m.expectStatic()
	teacher := &entities.TeacherProfile{ID: 5, UserID: 1, OrganizationID: testutil.TestOrgID}
	m.teachers.On("GetByUserID", mock.Anything, testutil.TestOrgID, int64(1)).Return(teacher, nil)

	got, err := m.usecase().Execute(context.Background(), inclusion.BuildContextRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionToolkit,
	})

	require.NoError(t, err)
	assert.Len(t, got.DeviceCatalog, 1)
	assert.Len(t, got.Situations, 1)
	require.NotNil(t, got.Teacher)
	assert.Nil(t, got.TargetStudent) // valija no carga alumno
	assert.NotContains(t, got.MissingData, "perfil_docente")
}

func TestBuildContext_MissingTeacherProfileIsSuggestedNotFatal(t *testing.T) {
	m := newCtxMocks()
	m.expectStatic()
	m.teachers.On("GetByUserID", mock.Anything, testutil.TestOrgID, int64(1)).
		Return(nil, providers.ErrNotFound)

	got, err := m.usecase().Execute(context.Background(), inclusion.BuildContextRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionTopic,
	})

	require.NoError(t, err)
	assert.Nil(t, got.Teacher)
	assert.Contains(t, got.MissingData, "perfil_docente")
}

func TestBuildContext_StudentDimensionLoadsFullContext(t *testing.T) {
	m := newCtxMocks()
	m.expectStatic()
	m.teachers.On("GetByUserID", mock.Anything, testutil.TestOrgID, int64(1)).
		Return(&entities.TeacherProfile{ID: 5, UserID: 1, OrganizationID: testutil.TestOrgID}, nil)

	student := testutil.NewStudentWithProfile(7, 3, "Pedro", []string{"se_distrae"})
	m.students.On("GetStudent", mock.Anything, testutil.TestOrgID, int64(7)).Return(&student, nil)
	m.diagnoses.On("ListByStudentProfile", mock.Anything, int64(7)).
		Return([]entities.StudentDiagnosis{{ID: 1, StudentProfileID: 7, DiagnosisID: 2}}, nil)
	m.ppi.On("GetByStudentID", mock.Anything, testutil.TestOrgID, int64(7)).
		Return(&entities.PPI{ID: 9, StudentID: 7, OrganizationID: testutil.TestOrgID, Status: "active"}, nil)
	m.adaptations.On("List", mock.Anything, testutil.TestOrgID, testutil.Ptr(int64(7)), mock.Anything).
		Return([]entities.Adaptation{testutil.NewAdaptation(1, 7, 1)}, nil)
	classroom := testutil.NewClassroom(3, "4to A")
	m.classrooms.On("Get", mock.Anything, testutil.TestOrgID, int64(3)).Return(&classroom, nil)
	m.students.On("ListByClassroom", mock.Anything, testutil.TestOrgID, int64(3)).
		Return([]entities.Student{student}, nil)
	m.summaries.On("RecentByStudent", mock.Anything, testutil.TestOrgID, int64(7), 10).
		Return([]entities.ConversationSummary{{ConversationID: 99}}, nil)

	got, err := m.usecase().Execute(context.Background(), inclusion.BuildContextRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionStudent,
		StudentID: testutil.Ptr(int64(7)),
	})

	require.NoError(t, err)
	require.NotNil(t, got.TargetStudent)
	assert.Equal(t, "Pedro", got.TargetStudent.Name)
	assert.Len(t, got.Diagnoses, 1)
	require.NotNil(t, got.PPI)
	assert.Len(t, got.PastAdaptations, 1)
	require.NotNil(t, got.Classroom)
	assert.Len(t, got.ClassroomStudents, 1)
	assert.Len(t, got.PriorSummaries, 1)
	assert.NotContains(t, got.MissingData, "ppi")
	m.students.AssertExpectations(t)
	m.ppi.AssertExpectations(t)
}

func TestBuildContext_StudentWithoutProfileAndPPIDegradesGracefully(t *testing.T) {
	m := newCtxMocks()
	m.expectStatic()
	m.teachers.On("GetByUserID", mock.Anything, testutil.TestOrgID, int64(1)).
		Return(&entities.TeacherProfile{ID: 5, UserID: 1, OrganizationID: testutil.TestOrgID}, nil)

	student := testutil.NewStudent(8, 0, "Sin Perfil") // sin Profile, sin classroom
	m.students.On("GetStudent", mock.Anything, testutil.TestOrgID, int64(8)).Return(&student, nil)
	m.ppi.On("GetByStudentID", mock.Anything, testutil.TestOrgID, int64(8)).
		Return(nil, providers.ErrNotFound)
	m.adaptations.On("List", mock.Anything, testutil.TestOrgID, testutil.Ptr(int64(8)), mock.Anything).
		Return([]entities.Adaptation{}, nil)
	m.summaries.On("RecentByStudent", mock.Anything, testutil.TestOrgID, int64(8), 10).
		Return([]entities.ConversationSummary{}, nil)

	got, err := m.usecase().Execute(context.Background(), inclusion.BuildContextRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    1,
		Dimension: inclusion.DimensionStudent,
		StudentID: testutil.Ptr(int64(8)),
	})

	require.NoError(t, err)
	require.NotNil(t, got.TargetStudent)
	assert.Nil(t, got.PPI)
	assert.Contains(t, got.MissingData, "perfil_alumno")
	assert.Contains(t, got.MissingData, "ppi")
	assert.Nil(t, got.Classroom) // classroom_id 0 => no se consulta
	m.diagnoses.AssertNotCalled(t, "ListByStudentProfile")
}

func TestBuildContext_SnapshotHasNoPII(t *testing.T) {
	student := testutil.NewStudentWithProfile(7, 3, "Pedro Sensible", []string{"x"})
	pc := &inclusion.PromptContext{
		Dimension:     inclusion.DimensionStudent,
		Teacher:       &entities.TeacherProfile{ID: 5, UserID: 1},
		TargetStudent: &student,
		PPI:           &entities.PPI{ID: 9, StudentID: 7},
		Diagnoses:     []entities.StudentDiagnosis{{ID: 1}},
		MissingData:   []string{"ppi"},
	}

	snap := pc.Snapshot()

	require.NotNil(t, snap.TargetStudentID)
	assert.Equal(t, int64(7), *snap.TargetStudentID)
	require.NotNil(t, snap.StudentProfileID)
	assert.Equal(t, int64(7), *snap.StudentProfileID)
	require.NotNil(t, snap.PPIID)
	assert.Equal(t, 1, snap.DiagnosesCount)
	// El snapshot solo lleva IDs/contadores: no debe haber forma de leer el nombre.
	assert.NotContains(t, snap.MissingData, "Pedro")
}
