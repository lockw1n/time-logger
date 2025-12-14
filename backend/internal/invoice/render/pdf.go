package render

import (
	"bytes"
	"fmt"
	"os/exec"
)

func PDF(html []byte) ([]byte, error) {
	bin, err := exec.LookPath("wkhtmltopdf")
	if err != nil {
		return nil, fmt.Errorf("wkhtmltopdf not found in PATH: %w", err)
	}

	cmd := exec.Command(
		bin,
		"--quiet",
		"--disable-smart-shrinking",
		"--encoding", "utf-8",
		"--page-size", "A4",
		"--margin-top", "10mm",
		"--margin-bottom", "10mm",
		"--margin-left", "10mm",
		"--margin-right", "10mm",
		"-", // stdin
		"-", // stdout
	)

	cmd.Stdin = bytes.NewReader(html)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(
			"wkhtmltopdf failed: %w: %s",
			err,
			stderr.String(),
		)
	}

	return stdout.Bytes(), nil
}
