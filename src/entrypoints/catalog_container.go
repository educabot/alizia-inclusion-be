package entrypoints

import (
	cataloguc "github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
)

type CatalogContainer struct {
	ListRamps   cataloguc.ListRamps
	GetRamp     cataloguc.GetRamp
	ListDevices cataloguc.ListDevices
	GetDevice   cataloguc.GetDevice
}
