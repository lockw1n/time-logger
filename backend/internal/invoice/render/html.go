package render

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed template/invoice.html
var templates embed.FS

func HTML(invoice Invoice) ([]byte, error) {
	tpl, err := template.ParseFS(
		templates,
		"template/invoice.html",
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, invoice); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
