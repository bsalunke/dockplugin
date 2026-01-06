package format

import (
	"bytes"
	"io"
	"text/template"

	"github.com/bsalunke/dockplugin/pkg/image"
)

// CustomFormatter formats images using a Go template
type CustomFormatter struct {
	Template string
	Writer   io.Writer
}

// Format writes images using custom template format
func (f *CustomFormatter) Format(images []*image.ImageInfo) error {
	tmpl, err := template.New("custom").Parse(f.Template)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, img := range images {
		if err := tmpl.Execute(&buf, img); err != nil {
			return err
		}
		buf.WriteString("\n")
	}

	_, err = f.Writer.Write(buf.Bytes())
	return err
}
