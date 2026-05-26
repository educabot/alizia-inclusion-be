package inclusion_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestExportAdaptation(t *testing.T) {
	ctx := context.Background()

	strategy := "Usar timer visual en bloques de 10 minutos"
	notes := "Revisar cada semana con la familia"
	newAdaptation := func() *entities.Adaptation {
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

	setupMock := func(a *entities.Adaptation, getErr error) *mocks.MockAdaptationProvider {
		return &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
				if getErr != nil {
					return nil, getErr
				}
				return a, nil
			},
		}
	}

	baseRequest := inclusion.ExportAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 7,
		Format:       inclusion.ExportFormatMarkdown,
	}

	t.Run("renders markdown with adaptation content", func(t *testing.T) {
		adaptations := setupMock(newAdaptation(), nil)

		doc, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, baseRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Filename != "adaptacion-7.md" {
			t.Errorf("expected filename adaptacion-7.md, got %q", doc.Filename)
		}
		if doc.ContentType != "text/markdown; charset=utf-8" {
			t.Errorf("unexpected content type %q", doc.ContentType)
		}
		content := string(doc.Data)
		for _, want := range []string{"Matemáticas", "Lucas", "Timer Visual", strategy, notes} {
			if !contains(content, want) {
				t.Errorf("markdown missing %q", want)
			}
		}
	})

	t.Run("renders a non-empty pdf with the pdf signature", func(t *testing.T) {
		adaptations := setupMock(newAdaptation(), nil)
		req := baseRequest
		req.Format = inclusion.ExportFormatPDF

		doc, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Filename != "adaptacion-7.pdf" {
			t.Errorf("expected filename adaptacion-7.pdf, got %q", doc.Filename)
		}
		if doc.ContentType != "application/pdf" {
			t.Errorf("unexpected content type %q", doc.ContentType)
		}
		if !bytes.HasPrefix(doc.Data, []byte("%PDF-")) {
			t.Errorf("expected PDF signature, got %q", firstBytes(doc.Data, 8))
		}
	})

	t.Run("propagates not found from repository", func(t *testing.T) {
		adaptations := setupMock(nil, providers.ErrAdaptationNotFound)

		_, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, baseRequest)
		if !errors.Is(err, providers.ErrAdaptationNotFound) {
			t.Errorf("expected ErrAdaptationNotFound, got: %v", err)
		}
	})

	t.Run("rejects unsupported format", func(t *testing.T) {
		adaptations := setupMock(newAdaptation(), nil)
		req := baseRequest
		req.Format = "docx"

		_, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		adaptations := setupMock(newAdaptation(), nil)
		req := baseRequest
		req.OrgID = uuid.Nil

		_, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects zero adaptation_id", func(t *testing.T) {
		adaptations := setupMock(newAdaptation(), nil)
		req := baseRequest
		req.AdaptationID = 0

		_, err := inclusion.NewExportAdaptation(adaptations).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})
}

func contains(haystack, needle string) bool {
	return bytes.Contains([]byte(haystack), []byte(needle))
}

func firstBytes(b []byte, n int) []byte {
	if len(b) < n {
		return b
	}
	return b[:n]
}
