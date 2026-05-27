package entrypoints

import (
	"testing"
	"time"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func TestMapDevice_Downloads(t *testing.T) {
	t.Run("maps embedded resources into downloads", func(t *testing.T) {
		// Arrange
		device := entities.Device{
			ID:   5,
			Name: "Time Timer",
			Resources: []entities.DeviceResource{
				{ID: 1, DeviceID: 5, Title: "Ficha imprimible", FileURL: "/files/5/ficha.pdf", FileType: "pdf", CreatedAt: time.Unix(0, 0).UTC()},
			},
		}

		// Act
		got := mapDevice(device)

		// Assert
		if len(got.Downloads) != 1 {
			t.Fatalf("expected 1 download, got %d", len(got.Downloads))
		}
		d := got.Downloads[0]
		if d.ID != 1 || d.Title != "Ficha imprimible" || d.FileURL != "/files/5/ficha.pdf" || d.FileType != "pdf" {
			t.Errorf("unexpected download mapping: %+v", d)
		}
	})

	t.Run("omits downloads when device has no resources", func(t *testing.T) {
		// Arrange
		device := entities.Device{ID: 6, Name: "Auriculares"}

		// Act
		got := mapDevice(device)

		// Assert
		if got.Downloads != nil {
			t.Errorf("expected nil downloads, got %+v", got.Downloads)
		}
	})
}
