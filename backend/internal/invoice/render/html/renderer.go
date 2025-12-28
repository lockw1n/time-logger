package html

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/invoice.html
//go:embed templates/footer.html
var templates embed.FS

func HTML(invoice InvoiceView) ([]byte, error) {
	tpl, err := template.ParseFS(
		templates,
		"templates/invoice.html",
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

func FooterHTML() ([]byte, error) {
	tpl, err := template.ParseFS(
		templates,
		"templates/footer.html",
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
