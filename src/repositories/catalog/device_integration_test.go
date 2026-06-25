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

func TestDeviceRepo_ListDevices_AllAndFilteredByRamp(t *testing.T) {
	tx := pgtest.Tx(t)
	ramp := testutil.NewRamp(1, "Sensorial")
	require.NoError(t, tx.Create(&ramp).Error)
	ramp2 := testutil.NewRamp(2, "Comunicación")
	require.NoError(t, tx.Create(&ramp2).Error)

	d1 := testutil.NewDevice(1, ramp.ID, "Timer")
	d2 := testutil.NewDevice(2, ramp2.ID, "Pictogramas")
	require.NoError(t, tx.Create(&d1).Error)
	require.NoError(t, tx.Create(&d2).Error)

	repo := catalog.NewDeviceRepo(tx)

	all, err := repo.ListDevices(context.Background(), testutil.TestOrgID, nil)
	require.NoError(t, err)
	assert.Len(t, all, 2)

	filtered, err := repo.ListDevices(context.Background(), testutil.TestOrgID, &ramp.ID)
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	assert.Equal(t, "Timer", filtered[0].Name)
	require.NotNil(t, filtered[0].Ramp, "Ramp is preloaded")
	assert.Equal(t, "Sensorial", filtered[0].Ramp.Name)
}

func TestDeviceRepo_GetDevice_ReturnsByID(t *testing.T) {
	tx := pgtest.Tx(t)
	ramp := testutil.NewRamp(1, "Sensorial")
	require.NoError(t, tx.Create(&ramp).Error)
	dev := testutil.NewDevice(1, ramp.ID, "Auriculares")
	require.NoError(t, tx.Create(&dev).Error)

	repo := catalog.NewDeviceRepo(tx)
	got, err := repo.GetDevice(context.Background(), testutil.TestOrgID, dev.ID)

	require.NoError(t, err)
	assert.Equal(t, "Auriculares", got.Name)
}

func TestDeviceRepo_GetDevice_NotFound(t *testing.T) {
	tx := pgtest.Tx(t)
	repo := catalog.NewDeviceRepo(tx)

	_, err := repo.GetDevice(context.Background(), testutil.TestOrgID, 999)

	assert.ErrorIs(t, err, providers.ErrNotFound)
}
