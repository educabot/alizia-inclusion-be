package dashboard_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
)

func TestGetAIUsage_AggregatesUsageAndMapsTheSummary(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	usage := &mocks.MockAIUsageProvider{
		SummarizeFn: func(_ context.Context, _ uuid.UUID, _ time.Time) (*providers.AIUsageSummary, error) {
			return &providers.AIUsageSummary{
				TotalRequests:    3,
				PromptTokens:     100,
				CompletionTokens: 40,
				TotalTokens:      140,
				ByMode: []providers.AIUsageModeSummary{
					{Mode: "assist", Requests: 2, TotalTokens: 90},
					{Mode: "recommend", Requests: 1, TotalTokens: 50},
				},
			}, nil
		},
	}

	got, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: orgID, Days: 7})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.WindowDays != 7 {
		t.Errorf("expected window_days 7, got %d", got.WindowDays)
	}
	if got.TotalTokens != 140 || got.TotalRequests != 3 {
		t.Errorf("unexpected totals: %+v", got)
	}
	if len(got.ByMode) != 2 || got.ByMode[0].Mode != "assist" {
		t.Errorf("unexpected by_mode: %+v", got.ByMode)
	}
}

func TestGetAIUsage_DefaultsTheWindowWhenDaysIsNotProvided(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	var capturedSince time.Time
	before := time.Now().AddDate(0, 0, -31)
	usage := &mocks.MockAIUsageProvider{
		SummarizeFn: func(_ context.Context, _ uuid.UUID, since time.Time) (*providers.AIUsageSummary, error) {
			capturedSince = since
			return &providers.AIUsageSummary{}, nil
		},
	}

	got, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: orgID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.WindowDays != 30 {
		t.Errorf("expected default window 30, got %d", got.WindowDays)
	}
	// The look-back must be roughly 30 days ago, i.e. after the 31-days-ago mark.
	if capturedSince.Before(before) {
		t.Errorf("since %v is older than the 30-day default window", capturedSince)
	}
}

func TestGetAIUsage_CapsAnExcessiveWindow(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	usage := &mocks.MockAIUsageProvider{}

	got, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: orgID, Days: 99999})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.WindowDays != 365 {
		t.Errorf("expected window capped at 365, got %d", got.WindowDays)
	}
}

func TestGetAIUsage_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	usage := &mocks.MockAIUsageProvider{}

	_, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: uuid.Nil})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}
