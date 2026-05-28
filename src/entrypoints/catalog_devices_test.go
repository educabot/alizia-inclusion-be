package entrypoints

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func TestMapDevice_MapsEmbeddedResourcesIntoDownloads(t *testing.T) {
	device := entities.Device{
		ID:   5,
		Name: "Time Timer",
		Resources: []entities.DeviceResource{
			{ID: 1, DeviceID: 5, Title: "Ficha imprimible", FileURL: "/files/5/ficha.pdf", FileType: "pdf", CreatedAt: time.Unix(0, 0).UTC()},
		},
	}

	got := mapDevice(device)

	assert.Len(t, got.Downloads, 1)
	d := got.Downloads[0]
	assert.Equal(t, int64(1), d.ID)
	assert.Equal(t, "Ficha imprimible", d.Title)
	assert.Equal(t, "/files/5/ficha.pdf", d.FileURL)
	assert.Equal(t, "pdf", d.FileType)
}

func TestMapDevice_OmitsDownloadsWhenDeviceHasNoResources(t *testing.T) {
	device := entities.Device{ID: 6, Name: "Auriculares"}

	got := mapDevice(device)

	assert.Nil(t, got.Downloads)
}
