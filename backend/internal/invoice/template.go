package invoice

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed template.html
var templateHTML string

type Item struct {
	Description string  `json:"description"`
	Hours       float64 `json:"hours"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
}

type Data struct {
	ConsultantName      string  `json:"consultant_name"`
	ConsultantAddress   string  `json:"consultant_address"`
	ConsultantTaxNumber string  `json:"consultant_tax_number"`
	CompanyName         string  `json:"company_name"`
	CompanyUID          string  `json:"company_uid"`
	CompanyStreet       string  `json:"company_street"`
	CompanyCity         string  `json:"company_city"`
	CompanyCountry      string  `json:"company_country"`
	OrderNumber         string  `json:"order_number"`
	InvoiceNumber       string  `json:"invoice_number"`
	InvoiceDate         string  `json:"invoice_date"`
	PeriodStart         string  `json:"period_start"`
	PeriodEnd           string  `json:"period_end"`
	HourlyRate          float64 `json:"hourly_rate"`
	Currency            string  `json:"currency"`
	PaymentCondition    string  `json:"payment_condition"`
	BankName            string  `json:"bank_name"`
	BankAddress         string  `json:"bank_address"`
	IBAN                string  `json:"iban"`
	BIC                 string  `json:"bic"`
	BankCountry         string  `json:"bank_country"`
	Items               []Item  `json:"items"`
	TotalHours          float64 `json:"total_hours"`
	TotalAmount         float64 `json:"total_amount"`
}

func moneyFormatter(currency string) func(float64) string {
	return func(value float64) string {
		formatted := fmt.Sprintf("%.2f", value)
		// use comma as decimal separator like the provided sample
		formatted = strings.ReplaceAll(formatted, ".", ",")
		return currency + formatted
	}
}

func hoursFormatter(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)
	return strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
}

// RenderHTML merges the template with data and returns the HTML string.
func RenderHTML(data Data) (string, error) {
	funcs := template.FuncMap{
		"money": moneyFormatter(data.Currency),
		"hours": hoursFormatter,
	}

	tpl, err := template.New("invoice").Funcs(funcs).Parse(templateHTML)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WKHTMLToPDF wraps calling wkhtmltopdf on the provided HTML file path and returns the generated PDF bytes.
func WKHTMLToPDF(htmlPath string) ([]byte, error) {
	bin, err := exec.LookPath("wkhtmltopdf")
	if err != nil {
		return nil, fmt.Errorf("wkhtmltopdf not found in PATH: %w", err)
	}

	outputPath := filepath.Join(filepath.Dir(htmlPath), filepath.Base(htmlPath)+".pdf")

	cmd := exec.Command(bin,
		"--quiet",
		"--disable-smart-shrinking",
		"--margin-top", "10mm",
		"--margin-bottom", "10mm",
		"--margin-left", "10mm",
		"--margin-right", "10mm",
		htmlPath,
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("wkhtmltopdf failed: %v: %s", err, stderr.String())
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("reading generated PDF failed: %w", err)
	}
	_ = os.Remove(outputPath)
	return content, nil
}

// GeneratePDF writes html to a temporary file, renders it with wkhtmltopdf, and returns the PDF bytes.
func GeneratePDF(html string) ([]byte, error) {
	tmp, err := os.CreateTemp("", "invoice-*.html")
	if err != nil {
		return nil, fmt.Errorf("unable to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(html); err != nil {
		return nil, fmt.Errorf("unable to write temp HTML: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("unable to close temp HTML: %w", err)
	}

	return WKHTMLToPDF(tmp.Name())
}
