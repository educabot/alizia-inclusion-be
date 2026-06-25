package dashboard_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetAIUsage_AggregatesUsageAndMapsTheSummary(t *testing.T) {
	ctx := context.Background()
	summary := providers.AIUsageSummary{
		TotalRequests:    3,
		PromptTokens:     100,
		CompletionTokens: 40,
		TotalTokens:      140,
		ByMode: []providers.AIUsageModeSummary{
			{Mode: "assist", Requests: 2, TotalTokens: 90},
			{Mode: "recommend", Requests: 1, TotalTokens: 50},
		},
	}
	usage := new(mockproviders.MockAIUsageProvider)
	usage.On("Summarize", ctx, testutil.TestOrgID, mock.AnythingOfType("time.Time")).Return(&summary, nil)

	got, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: testutil.TestOrgID, Days: 7})

	assert.NoError(t, err)
	assert.Equal(t, 7, got.WindowDays)
	assert.Equal(t, 140, got.TotalTokens)
	assert.Equal(t, 3, got.TotalRequests)
	assert.Len(t, got.ByMode, 2)
	assert.Equal(t, "assist", got.ByMode[0].Mode)
	usage.AssertExpectations(t)
}

func TestGetAIUsage_DefaultsTheWindowWhenDaysIsNotProvided(t *testing.T) {
	ctx := context.Background()
	usage := new(mockproviders.MockAIUsageProvider)
	usage.On("Summarize", ctx, testutil.TestOrgID, mock.AnythingOfType("time.Time")).Return(&providers.AIUsageSummary{}, nil)

	got, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: testutil.TestOrgID})

	assert.NoError(t, err)
	assert.Equal(t, 30, got.WindowDays)
	usage.AssertExpectations(t)
}

func TestGetAIUsage_CapsAnExcessiveWindow(t *testing.T) {
	ctx := context.Background()
	usage := new(mockproviders.MockAIUsageProvider)
	usage.On("Summarize", ctx, testutil.TestOrgID, mock.AnythingOfType("time.Time")).Return(&providers.AIUsageSummary{}, nil)

	got, err := dashboard.NewGetAIUsage(usage).Execute(ctx, dashboard.GetAIUsageRequest{OrgID: testutil.TestOrgID, Days: 99999})

	assert.NoError(t, err)
	assert.Equal(t, 365, got.WindowDays)
	usage.AssertExpectations(t)
}

func TestGetAIUsage_RejectsNilOrgID(t *testing.T) {
	usage := new(mockproviders.MockAIUsageProvider)

	_, err := dashboard.NewGetAIUsage(usage).Execute(context.Background(), dashboard.GetAIUsageRequest{OrgID: uuid.Nil})

	assert.ErrorIs(t, err, providers.ErrValidation)
	usage.AssertNotCalled(t, "Summarize", mock.Anything, mock.Anything, mock.Anything)
}
