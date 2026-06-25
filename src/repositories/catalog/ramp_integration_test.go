//go:build integration

package catalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/repositories/catalog"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
	"github.com/educabot/alizia-inclusion-be/src/testutil/pgtest"
)

func TestRampRepo_ListRamps_OrdersBySortOrderAndScopesToOrg(t *testing.T) {
	tx := pgtest.Tx(t)
	r2 := testutil.NewRamp(1, "Comunicación")
	r2.SortOrder = 2
	r1 := testutil.NewRamp(2, "Sensorial")
	r1.SortOrder = 1
	require.NoError(t, tx.Create(&r2).Error)
	require.NoError(t, tx.Create(&r1).Error)
	// A ramp in another tenant must not leak into the result.
	other := testutil.NewRamp(3, "Otra org")
	other.OrganizationID = pgtest.OtherOrgID
	require.NoError(t, tx.Create(&other).Error)

	repo := catalog.NewRampRepo(tx)
	got, err := repo.ListRamps(context.Background(), testutil.TestOrgID)

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Sensorial", got[0].Name, "ordered by sort_order ASC")
	assert.Equal(t, "Comunicación", got[1].Name)
}

func TestRampRepo_GetRamp_ReturnsByID(t *testing.T) {
	tx := pgtest.Tx(t)
	ramp := testutil.NewRamp(1, "Sensorial")
	require.NoError(t, tx.Create(&ramp).Error)

	repo := catalog.NewRampRepo(tx)
	got, err := repo.GetRamp(context.Background(), testutil.TestOrgID, ramp.ID)

	require.NoError(t, err)
	assert.Equal(t, "Sensorial", got.Name)
}

func TestRampRepo_GetRamp_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := catalog.NewRampRepo(tx)

	_, err := repo.GetRamp(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestRampRepo_GetRamp_NotFoundAcrossTenant(t *testing.T) {
	tx := pgtest.Tx(t)
	ramp := testutil.NewRamp(1, "De otra org")
	ramp.OrganizationID = pgtest.OtherOrgID
	require.NoError(t, tx.Create(&ramp).Error)

	repo := catalog.NewRampRepo(tx)
	_, err := repo.GetRamp(context.Background(), testutil.TestOrgID, ramp.ID)

	assert.ErrorIs(t, err, providers.ErrNotFound, "must not read another tenant's ramp")
}
