package inclusion_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func newExportAdaptation() *entities.Adaptation {
	strategy := "Usar timer visual en bloques de 10 minutos"
	notes := "Revisar cada semana con la familia"
	student := testutil.NewStudent(1, 1, "Lucas")
	device := testutil.NewDevice(1, 1, "Timer Visual")
	return &entities.Adaptation{
		ID:                 7,
		OrganizationID:     testutil.TestOrgID,
		StudentID:          1,
		Subject:            "Matemáticas",
		AdaptationType:     "actividad_adaptada",
		AdaptationStrategy: &strategy,
		Notes:              &notes,
		Status:             "en_curso",
		Student:            &student,
		Devices:            []entities.Device{device},
	}
}

var baseExportRequest = inclusion.ExportAdaptationRequest{
	OrgID:        testutil.TestOrgID,
	AdaptationID: 7,
	Format:       inclusion.ExportFormatMarkdown,
}

func TestExportAdaptation_RendersMarkdownWithAdaptationContent(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(7)).Return(newExportAdaptation(), nil)

	doc, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, baseExportRequest)

	require.NoError(t, err)
	assert.Equal(t, "adaptacion-7.md", doc.Filename)
	assert.Equal(t, "text/markdown; charset=utf-8", doc.ContentType)
	content := string(doc.Data)
	for _, want := range []string{"Matemáticas", "Lucas", "Timer Visual", "Usar timer visual en bloques de 10 minutos", "Revisar cada semana con la familia"} {
		assert.Contains(t, content, want)
	}
	adaptations.AssertExpectations(t)
}

func TestExportAdaptation_RendersNonEmptyPDFWithPDFSignature(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(7)).Return(newExportAdaptation(), nil)
	req := baseExportRequest
	req.Format = inclusion.ExportFormatPDF

	doc, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "adaptacion-7.pdf", doc.Filename)
	assert.Equal(t, "application/pdf", doc.ContentType)
	assert.True(t, bytes.HasPrefix(doc.Data, []byte("%PDF-")))
	adaptations.AssertExpectations(t)
}

func TestExportAdaptation_PropagatesNotFoundFromRepository(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(7)).Return(nil, providers.ErrAdaptationNotFound)

	_, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, baseExportRequest)

	assert.ErrorIs(t, err, providers.ErrAdaptationNotFound)
	adaptations.AssertExpectations(t)
}

func TestExportAdaptation_RejectsUnsupportedFormat(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	req := baseExportRequest
	req.Format = "docx"

	_, err := inclusion.NewExportAdaptation(adaptations).Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestExportAdaptation_RejectsNilOrgID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	req := baseExportRequest
	req.OrgID = uuid.Nil

	_, err := inclusion.NewExportAdaptation(adaptations).Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestExportAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	req := baseExportRequest
	req.AdaptationID = 0

	_, err := inclusion.NewExportAdaptation(adaptations).Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}
