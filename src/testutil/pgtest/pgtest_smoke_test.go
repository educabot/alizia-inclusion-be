//go:build integration

package pgtest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

// TestHarness_MigratesAndConnects verifies the container starts, every migration
// applies, and key tables and the vector extension exist.
func TestHarness_MigratesAndConnects(t *testing.T) {
	tx := pgtest.Tx(t)

	var one int
	require.NoError(t, tx.Raw("SELECT 1").Scan(&one).Error)
	assert.Equal(t, 1, one)

	// A late migration's table must exist (proves all .up.sql ran in order).
	var count int64
	require.NoError(t, tx.Raw("SELECT count(*) FROM pedagogical_content").Scan(&count).Error)

	// pgvector extension must be installed (migration 000021).
	var hasVector bool
	require.NoError(t, tx.Raw("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&hasVector).Error)
	assert.True(t, hasVector, "vector extension should be enabled")
}
