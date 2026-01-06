package format

import (
	"encoding/json"
	"io"

	"github.com/bsalunke/dockplugin/pkg/image"
)

// JSONFormatter formats images as JSON
type JSONFormatter struct {
	Writer io.Writer
}

// Format writes images in JSON format
func (f *JSONFormatter) Format(images []*image.ImageInfo) error {
	encoder := json.NewEncoder(f.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(images)
}
