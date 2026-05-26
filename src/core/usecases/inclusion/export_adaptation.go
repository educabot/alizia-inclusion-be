package inclusion

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

const (
	ExportFormatMarkdown = "md"
	ExportFormatPDF      = "pdf"
)

type ExportedDocument struct {
	Filename    string
	ContentType string
	Data        []byte
}

type ExportAdaptationRequest struct {
	OrgID        uuid.UUID
	AdaptationID int64
	Format       string
}

func (r ExportAdaptationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.AdaptationID <= 0 {
		return errAdaptationIDRequired
	}
	switch strings.ToLower(r.Format) {
	case ExportFormatMarkdown, ExportFormatPDF:
		return nil
	default:
		return errInvalidExportFormat
	}
}

type ExportAdaptation interface {
	Execute(ctx context.Context, req ExportAdaptationRequest) (*ExportedDocument, error)
}

type exportAdaptationImpl struct {
	adaptations providers.AdaptationProvider
}

func NewExportAdaptation(adaptations providers.AdaptationProvider) ExportAdaptation {
	return &exportAdaptationImpl{adaptations: adaptations}
}

func (uc *exportAdaptationImpl) Execute(ctx context.Context, req ExportAdaptationRequest) (*ExportedDocument, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	adaptation, err := uc.adaptations.Get(ctx, req.OrgID, req.AdaptationID)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(req.Format) {
	case ExportFormatMarkdown:
		return &ExportedDocument{
			Filename:    fmt.Sprintf("adaptacion-%d.md", adaptation.ID),
			ContentType: "text/markdown; charset=utf-8",
			Data:        renderAdaptationMarkdown(adaptation),
		}, nil
	case ExportFormatPDF:
		data, perr := renderAdaptationPDF(adaptation)
		if perr != nil {
			return nil, fmt.Errorf("render adaptation pdf: %w", perr)
		}
		return &ExportedDocument{
			Filename:    fmt.Sprintf("adaptacion-%d.pdf", adaptation.ID),
			ContentType: "application/pdf",
			Data:        data,
		}, nil
	default:
		return nil, errInvalidExportFormat
	}
}
