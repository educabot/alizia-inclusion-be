// Package static sirve assets embebidos en el binario (imágenes de devices, etc.).
// Embeber los archivos evita depender del filesystem en runtime: el deploy es
// atómico (Railway no necesita montar volúmenes) y la imagen vive junto al código.
package static

import (
	"embed"
	"io/fs"
)

//go:embed images
var assets embed.FS

// Images devuelve el subárbol de assets a partir de "images", de modo que la
// raíz del FS sea el directorio servido (ej: /images/devices/ETE-XXXX-EB.png
// mapea a images/devices/ETE-XXX-EB.png dentro del embed).
func Images() fs.FS {
	sub, err := fs.Sub(assets, "images")
	if err != nil {
		// El path es estático y embebido en compile-time; si fs.Sub falla acá es
		// un bug de build, no una condición de runtime recuperable.
		panic(err)
	}
	return sub
}
