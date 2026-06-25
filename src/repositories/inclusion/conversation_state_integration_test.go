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

func assistantMessage(t *testing.T, msgs []entities.ConversationMessage) entities.ConversationMessage {
	t.Helper()
	for _, m := range msgs {
		if m.Role == "assistant" {
			return m
		}
	}
	t.Fatal("no assistant message found")
	return entities.ConversationMessage{}
}

func TestConversationRepo_AppendTurnCreatesThenAppends(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Seño", "teacher")
	repo := inclusion.NewConversationRepo(tx)

	convID, err := repo.AppendTurn(context.Background(), providers.AppendTurnParams{
		OrgID: testutil.TestOrgID, UserID: userID, Mode: "assist",
		UserContent: "Hola", AssistantContent: "Buenas",
		Metadata: map[string]any{"recommended_device": 7},
	})
	require.NoError(t, err)
	require.NotZero(t, convID)

	got, err := repo.GetWithMessages(context.Background(), testutil.TestOrgID, convID)
	require.NoError(t, err)
	require.Len(t, got.Messages, 2)
	assert.Contains(t, string(assistantMessage(t, got.Messages).Metadata), "recommended_device", "metadata persisted as JSONB")

	// Appending to the existing conversation adds another pair of messages.
	_, err = repo.AppendTurn(context.Background(), providers.AppendTurnParams{
		ConversationID: convID, OrgID: testutil.TestOrgID, UserID: userID, Mode: "assist",
		UserContent: "Otra", AssistantContent: "Respuesta",
	})
	require.NoError(t, err)

	got2, err := repo.GetWithMessages(context.Background(), testutil.TestOrgID, convID)
	require.NoError(t, err)
	assert.Len(t, got2.Messages, 4)
}

func TestConversationRepo_GetWithMessages_NotFoundAcrossOrg(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewConversationRepo(tx)

	_, err := repo.GetWithMessages(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestConversationRepo_ListByUser_FiltersByMode(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Seño", "teacher")
	repo := inclusion.NewConversationRepo(tx)
	_, err := repo.AppendTurn(context.Background(), providers.AppendTurnParams{OrgID: testutil.TestOrgID, UserID: userID, Mode: "assist", UserContent: "a", AssistantContent: "b"})
	require.NoError(t, err)
	_, err = repo.AppendTurn(context.Background(), providers.AppendTurnParams{OrgID: testutil.TestOrgID, UserID: userID, Mode: "recommend", UserContent: "a", AssistantContent: "b"})
	require.NoError(t, err)

	assistOnly, err := repo.ListByUser(context.Background(), testutil.TestOrgID, userID, "assist")
	require.NoError(t, err)
	require.Len(t, assistOnly, 1)
	assert.Equal(t, "assist", assistOnly[0].Mode)

	all, err := repo.ListByUser(context.Background(), testutil.TestOrgID, userID, "")
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

// seedConversation inserts a bare conversation for the given user and returns its id.
func seedConversation(t *testing.T, tx *gorm.DB, userID int64, mode string) int64 {
	t.Helper()
	c := entities.Conversation{OrganizationID: testutil.TestOrgID, UserID: userID, Mode: mode}
	require.NoError(t, tx.Create(&c).Error)
	return c.ID
}

func TestConversationSummaryRepo_UpsertAndRecentByStudent(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Seño", "teacher")
	studentID := seedStudent(t, tx, "Mati")
	convID := seedConversation(t, tx, userID, "assist")
	repo := inclusion.NewConversationSummaryRepo(tx)

	sum := &entities.ConversationSummary{ConversationID: convID, Summary: "Trabajamos pausas", TopicKeywords: []string{"autorregulacion"}, TokenCount: 42}
	require.NoError(t, repo.Upsert(context.Background(), sum, []int64{studentID}, nil))

	got, err := repo.RecentByStudent(context.Background(), testutil.TestOrgID, studentID, 10)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Trabajamos pausas", got[0].Summary)

	// Idempotent: re-upserting the same conversation replaces the summary text.
	sum.Summary = "Actualizado"
	require.NoError(t, repo.Upsert(context.Background(), sum, []int64{studentID}, nil))
	got2, err := repo.RecentByStudent(context.Background(), testutil.TestOrgID, studentID, 10)
	require.NoError(t, err)
	require.Len(t, got2, 1, "still a single summary row (1:1 on conversation_id)")
	assert.Equal(t, "Actualizado", got2[0].Summary)
}

func TestConversationSummaryRepo_RecentByTopic(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Seño", "teacher")
	convID := seedConversation(t, tx, userID, "assist")
	repo := inclusion.NewConversationSummaryRepo(tx)
	require.NoError(t, repo.Upsert(context.Background(),
		&entities.ConversationSummary{ConversationID: convID, Summary: "s", TopicKeywords: []string{"fracciones", "mate"}}, nil, nil))

	got, err := repo.RecentByTopic(context.Background(), testutil.TestOrgID, "fracciones", 10)

	require.NoError(t, err)
	require.Len(t, got, 1, "topic_keywords array contains the keyword (Postgres ANY)")
}

func TestAIUsageRepo_RecordPersistsSnapshotAndSummarizes(t *testing.T) {
	tx := pgtest.Tx(t)
	userID := seedUser(t, tx, "Seño", "teacher")
	repo := inclusion.NewAIUsageRepo(tx)

	require.NoError(t, repo.Record(context.Background(), providers.AIUsageRecord{
		OrgID: testutil.TestOrgID, UserID: userID, Mode: "assist",
		PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15, Model: "gpt-4o",
		ContextSnapshot: map[string]any{"student_id": 3},
	}))
	require.NoError(t, repo.Record(context.Background(), providers.AIUsageRecord{
		OrgID: testutil.TestOrgID, UserID: userID, Mode: "assist",
		PromptTokens: 20, CompletionTokens: 10, TotalTokens: 30, Model: "gpt-4o",
	}))

	// Snapshot persisted as JSONB.
	var row entities.AIUsage
	require.NoError(t, tx.Where("user_id = ? AND total_tokens = 15", userID).First(&row).Error)
	assert.Contains(t, string(row.ContextSnapshot), "student_id")

	// Summarize aggregates the two assist turns.
	summary, err := repo.Summarize(context.Background(), testutil.TestOrgID, time.Now().Add(-time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 2, summary.TotalRequests)
	assert.Equal(t, 45, summary.TotalTokens)
	require.Len(t, summary.ByMode, 1)
	assert.Equal(t, "assist", summary.ByMode[0].Mode)
}
