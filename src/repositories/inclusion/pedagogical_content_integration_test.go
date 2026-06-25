//go:build integration

package inclusion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/repositories/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

// seedContent inserts a published pedagogical document with one chunk and returns
// the content id. keywords and tags must be non-nil: the columns are nullable but
// array_to_string(NULL) would null out the whole weighted tsvector.
func seedContent(t *testing.T, tx *gorm.DB, orgID *uuid.UUID, title string, keywords []string, chunkText string, tags []string) int64 {
	t.Helper()
	c := entities.PedagogicalContent{
		Title:          &title,
		Type:           testutil.Ptr("material"),
		Status:         "published",
		Keywords:       append(pq.StringArray{}, keywords...),
		OrganizationID: orgID,
	}
	require.NoError(t, tx.Create(&c).Error)
	chunk := entities.PedagogicalContentChunk{
		ContentID: c.ID,
		ChunkText: &chunkText,
		Tags:      append(pq.StringArray{}, tags...),
	}
	require.NoError(t, tx.Create(&chunk).Error)
	return c.ID
}

func TestPedagogicalContentRepo_SearchMatchesByKeyword(t *testing.T) {
	tx := pgtest.Tx(t)
	seedContent(t, tx, &testutil.TestOrgID, "Guía TEA", []string{"autismo"}, "cuerpo del texto", []string{"aula"})
	repo := inclusion.NewPedagogicalContentRepo(tx)

	got, err := repo.SearchChunks(context.Background(), testutil.TestOrgID, "autismo", 10)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Guía TEA", got[0].Title)
	assert.Greater(t, got[0].Score, float64(0))
}

func TestPedagogicalContentRepo_KeywordOutranksBodyMatch(t *testing.T) {
	tx := pgtest.Tx(t)
	// Doc A matches "dislexia" in its keywords (weight A); Doc B only in the body (weight C).
	seedContent(t, tx, &testutil.TestOrgID, "Doc A", []string{"dislexia"}, "contenido sin relación", []string{"x"})
	seedContent(t, tx, &testutil.TestOrgID, "Doc B", []string{"otra"}, "estrategias para la dislexia en el aula", []string{"y"})
	repo := inclusion.NewPedagogicalContentRepo(tx)

	got, err := repo.SearchChunks(context.Background(), testutil.TestOrgID, "dislexia", 10)

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Doc A", got[0].Title, "keyword match (weight A) ranks above body match (weight C)")
	assert.Greater(t, got[0].Score, got[1].Score)
}

func TestPedagogicalContentRepo_OrSemanticsMatchesAnyTerm(t *testing.T) {
	tx := pgtest.Tx(t)
	seedContent(t, tx, &testutil.TestOrgID, "Fracciones", []string{"fracciones"}, "matemática", []string{"primaria"})
	repo := inclusion.NewPedagogicalContentRepo(tx)

	// "inexistente" appears nowhere; AND semantics would return nothing, OR returns the doc.
	got, err := repo.SearchChunks(context.Background(), testutil.TestOrgID, "fracciones inexistente", 10)

	require.NoError(t, err)
	require.Len(t, got, 1, "orTSQuery gives OR semantics: any term matching is enough")
	assert.Equal(t, "Fracciones", got[0].Title)
}

func TestPedagogicalContentRepo_SearchEmptyAndNoMatch(t *testing.T) {
	tx := pgtest.Tx(t)
	seedContent(t, tx, &testutil.TestOrgID, "Algo", []string{"autismo"}, "texto", []string{"aula"})
	repo := inclusion.NewPedagogicalContentRepo(tx)

	empty, err := repo.SearchChunks(context.Background(), testutil.TestOrgID, "   ", 10)
	require.NoError(t, err)
	assert.Empty(t, empty, "blank query short-circuits to empty")

	noMatch, err := repo.SearchChunks(context.Background(), testutil.TestOrgID, "palabrainexistente", 10)
	require.NoError(t, err)
	assert.Empty(t, noMatch)
}

func TestPedagogicalContentRepo_ScopesToPublishedGlobalAndOwnOrg(t *testing.T) {
	tx := pgtest.Tx(t)
	// global (published) + own org (published) match; draft and other-org must not.
	seedContent(t, tx, nil, "Global", []string{"sensorial"}, "x", []string{"a"})
	seedContent(t, tx, &testutil.TestOrgID, "Propio", []string{"sensorial"}, "x", []string{"a"})
	seedContent(t, tx, &pgtest.OtherOrgID, "Ajeno", []string{"sensorial"}, "x", []string{"a"})
	draft := entities.PedagogicalContent{Title: testutil.Ptr("Borrador"), Status: "draft", Keywords: pq.StringArray{"sensorial"}, OrganizationID: &testutil.TestOrgID}
	require.NoError(t, tx.Create(&draft).Error)
	require.NoError(t, tx.Create(&entities.PedagogicalContentChunk{ContentID: draft.ID, ChunkText: testutil.Ptr("x"), Tags: pq.StringArray{"a"}}).Error)
	repo := inclusion.NewPedagogicalContentRepo(tx)

	got, err := repo.SearchChunks(context.Background(), testutil.TestOrgID, "sensorial", 10)

	require.NoError(t, err)
	titles := []string{}
	for _, r := range got {
		titles = append(titles, r.Title)
	}
	assert.ElementsMatch(t, []string{"Global", "Propio"}, titles, "draft and other-org excluded")
}

func TestPedagogicalContentRepo_GetContentWithChunks(t *testing.T) {
	tx := pgtest.Tx(t)
	id := seedContent(t, tx, &testutil.TestOrgID, "Doc", []string{"k"}, "cuerpo", []string{"t"})
	repo := inclusion.NewPedagogicalContentRepo(tx)

	got, err := repo.GetContent(context.Background(), testutil.TestOrgID, id)

	require.NoError(t, err)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Doc", *got.Title)
	require.Len(t, got.Chunks, 1, "chunks preloaded")
}

func TestPedagogicalContentRepo_GetContent_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := inclusion.NewPedagogicalContentRepo(tx)

	_, err := repo.GetContent(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}
