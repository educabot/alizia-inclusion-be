package entrypoints

import (
	"testing"
	"time"

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

	if len(got.Downloads) != 1 {
		t.Fatalf("expected 1 download, got %d", len(got.Downloads))
	}
	d := got.Downloads[0]
	if d.ID != 1 || d.Title != "Ficha imprimible" || d.FileURL != "/files/5/ficha.pdf" || d.FileType != "pdf" {
		t.Errorf("unexpected download mapping: %+v", d)
	}
}

func TestMapDevice_OmitsDownloadsWhenDeviceHasNoResources(t *testing.T) {
	device := entities.Device{ID: 6, Name: "Auriculares"}

	got := mapDevice(device)

	if got.Downloads != nil {
		t.Errorf("expected nil downloads, got %+v", got.Downloads)
	}
}
