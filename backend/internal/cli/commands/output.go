package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// printRawJSON pretty-prints a raw response body without re-parsing it, so
// --json reflects exactly what the backend returned. Falls back to the raw
// bytes if they somehow aren't valid JSON to indent.
func printRawJSON(w io.Writer, raw []byte) error {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		_, err = fmt.Fprintln(w, string(raw))
		return err
	}
	_, err := fmt.Fprintln(w, buf.String())
	return err
}
