package tensorflow

import (
	"image"

	"github.com/photoprism/photoprism/pkg/fs"
)

// OpenImage opens a natively supported image file and decodes it directly.
func OpenImage(fileName string) (image.Image, error) {
	img, _, err := fs.DecodeImageFile(fileName)
	return img, err
}
