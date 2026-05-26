package inclusion_test

import (
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

func TestGetChatHistory(t *testing.T) {
	ctx := context.Background()

	t.Run("returns conversations", func(t *testing.T) {
		expected := []entities.Conversation{
			testutil.NewConversation(1, 1, "recommend"),
			testutil.NewConversation(2, 1, "recommend"),
		}
		mock := &mocks.MockConversationProvider{
			ListByUserFn: func(_ context.Context, _ uuid.UUID, userID int64, mode string) ([]entities.Conversation, error) {
				if userID != 1 {
					t.Errorf("expected userID 1, got %d", userID)
				}
				if mode != "recommend" {
					t.Errorf("expected mode %q, got %q", "recommend", mode)
				}
				return expected, nil
			},
		}

		got, err := inclusion.NewGetChatHistory(mock).Execute(ctx, inclusion.GetChatHistoryRequest{
			OrgID:  testutil.TestOrgID,
			UserID: 1,
			Mode:   "recommend",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("got %d conversations, want 2", len(got))
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockConversationProvider{}
		_, err := inclusion.NewGetChatHistory(mock).Execute(ctx, inclusion.GetChatHistoryRequest{
			OrgID: uuid.Nil, UserID: 1, Mode: "recommend",
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects zero user_id", func(t *testing.T) {
		mock := &mocks.MockConversationProvider{}
		_, err := inclusion.NewGetChatHistory(mock).Execute(ctx, inclusion.GetChatHistoryRequest{
			OrgID: testutil.TestOrgID, UserID: 0, Mode: "recommend",
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects empty mode", func(t *testing.T) {
		mock := &mocks.MockConversationProvider{}
		_, err := inclusion.NewGetChatHistory(mock).Execute(ctx, inclusion.GetChatHistoryRequest{
			OrgID: testutil.TestOrgID, UserID: 1, Mode: "",
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})
}
