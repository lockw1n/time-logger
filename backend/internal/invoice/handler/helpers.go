package handler

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/lockw1n/time-logger/internal/constants"
)

const (
	DispositionInline     = "inline"
	DispositionAttachment = "attachment"
)

var invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

func buildInvoiceFilename(fullName string, invoiceNumber string, ext string) string {
	safeName := sanitizeFilenamePart(fullName)

	return fmt.Sprintf(
		"%s_Invoice_%s.%s",
		safeName,
		invoiceNumber,
		ext,
	)
}

func sanitizeFilenamePart(s string) string {
	// Trim spaces
	s = strings.TrimSpace(s)

	// Replace spaces with underscore
	s = strings.ReplaceAll(s, " ", "_")

	// Remove invalid filename characters
	s = invalidFilenameChars.ReplaceAllString(s, "")

	// Collapse multiple underscores
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	return s
}

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(constants.InternalDateFormat, date)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return t, nil
}

func buildContentDisposition(disposition string, filename string) string {
	if disposition != DispositionInline && disposition != DispositionAttachment {
		disposition = DispositionAttachment
	}

	asciiFallback := asciiFilename(filename)
	utf8Filename := url.PathEscape(filename)

	return fmt.Sprintf(
		`%s; filename="%s"; filename*=UTF-8''%s`,
		disposition,
		asciiFallback,
		utf8Filename,
	)
}

func asciiFilename(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r <= 127 {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}
