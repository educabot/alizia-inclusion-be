package inclusion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestCreateAdaptation_CreatesAdaptationWithoutDevices(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Create", ctx, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			a, ok := args.Get(1).(*entities.Adaptation)
			require.True(t, ok)
			a.ID = 1
		}).
		Return(nil)
	got := testutil.NewAdaptation(1, 1, 1)
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&got, nil)

	result, err := inclusion.NewCreateAdaptation(adaptations).Execute(ctx, inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      testutil.Ptr(int64(1)),
		TeacherID:      1,
		Subject:        "Matematicas",
		AdaptationType: "actividad_adaptada",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	adaptations.AssertExpectations(t)
	adaptations.AssertNotCalled(t, "SetDevices", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateAdaptation_CreatesAdaptationWithDevices(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Create", ctx, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			a, ok := args.Get(1).(*entities.Adaptation)
			require.True(t, ok)
			a.ID = 1
		}).
		Return(nil)
	adaptations.On("SetDevices", ctx, int64(1), []int64{10, 20}).Return(nil)
	got := testutil.NewAdaptation(1, 1, 1)
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&got, nil)

	result, err := inclusion.NewCreateAdaptation(adaptations).Execute(ctx, inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      testutil.Ptr(int64(1)),
		TeacherID:      1,
		Subject:        "Matematicas",
		AdaptationType: "actividad_adaptada",
		DeviceIDs:      []int64{10, 20},
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	adaptations.AssertExpectations(t)
}

func TestCreateAdaptation_DefaultsAdaptationTypeAndPersistsTitleWhenTypeOmitted(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	var captured *entities.Adaptation
	adaptations.On("Create", ctx, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			a, ok := args.Get(1).(*entities.Adaptation)
			require.True(t, ok)
			a.ID = 1
			captured = a
		}).
		Return(nil)
	got := testutil.NewAdaptation(1, 1, 1)
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&got, nil)

	_, err := inclusion.NewCreateAdaptation(adaptations).Execute(ctx, inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      testutil.Ptr(int64(1)),
		TeacherID:      1,
		Title:          "Secuencia con apoyos visuales",
		Subject:        "Matematicas",
		AdaptationType: "",
	})

	require.NoError(t, err)
	require.NotNil(t, captured)
	assert.Equal(t, "actividad_adaptada", captured.AdaptationType)
	assert.Equal(t, "Secuencia con apoyos visuales", captured.Title)
	adaptations.AssertExpectations(t)
}

func TestCreateAdaptation_RejectsNilOrgID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewCreateAdaptation(adaptations).Execute(context.Background(), inclusion.CreateAdaptationRequest{
		OrgID:     uuid.Nil,
		StudentID: testutil.Ptr(int64(1)),
		TeacherID: 1,
		Subject:   "Matematicas",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateAdaptation_AllowsAdaptationWithoutStudent(t *testing.T) {
	// Recurso asociado a una situación, sin alumno registrado (student_id nil).
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Create", ctx, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			a, ok := args.Get(1).(*entities.Adaptation)
			require.True(t, ok)
			a.ID = 7
		}).
		Return(nil)
	got := testutil.NewAdaptation(7, 1, 1)
	got.StudentID = nil
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(7)).Return(&got, nil)

	result, err := inclusion.NewCreateAdaptation(adaptations).Execute(ctx, inclusion.CreateAdaptationRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: nil,
		TeacherID: 1,
		Subject:   "Situación de aula",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.StudentID)
	adaptations.AssertExpectations(t)
}

func TestCreateAdaptation_RejectsExplicitZeroStudentID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewCreateAdaptation(adaptations).Execute(context.Background(), inclusion.CreateAdaptationRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: testutil.Ptr(int64(0)),
		TeacherID: 1,
		Subject:   "Matematicas",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateAdaptation_RejectsZeroTeacherID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewCreateAdaptation(adaptations).Execute(context.Background(), inclusion.CreateAdaptationRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: testutil.Ptr(int64(1)),
		TeacherID: 0,
		Subject:   "Matematicas",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateAdaptation_RejectsInvalidType(t *testing.T) {
	// GAP 2: un adaptation_type no vacío debe pertenecer al enum; si no, error de validación.
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewCreateAdaptation(adaptations).Execute(context.Background(), inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      testutil.Ptr(int64(1)),
		TeacherID:      1,
		Subject:        "Matematicas",
		AdaptationType: "tipo_invalido",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// El subject (materia) es opcional: el flujo del docente lo descartó y el guardado
// desde el chat no lo envía. Un subject vacío debe crear la adaptación igual.
func TestCreateAdaptation_AllowsEmptySubject(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Create", ctx, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			a, ok := args.Get(1).(*entities.Adaptation)
			require.True(t, ok)
			a.ID = 1
		}).
		Return(nil)
	got := testutil.NewAdaptation(1, 1, 1)
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&got, nil)

	result, err := inclusion.NewCreateAdaptation(adaptations).Execute(ctx, inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      testutil.Ptr(int64(1)),
		TeacherID:      1,
		Subject:        "",
		AdaptationType: "estrategia_aula",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	adaptations.AssertExpectations(t)
}
