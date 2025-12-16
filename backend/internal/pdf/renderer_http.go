package pdf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type HTTPRenderer struct {
	baseURL string
	token   string
	client  *http.Client
}

type renderRequest struct {
	HTML      string `json:"html"`
	TimeoutMs int    `json:"timeoutMs,omitempty"`
}

func NewHTTPRenderer(baseURL, token string) *HTTPRenderer {
	return &HTTPRenderer{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (r *HTTPRenderer) RenderHTML(ctx context.Context, html string) ([]byte, error) {
	if html == "" {
		return nil, ErrEmptyHTML
	}

	payload := renderRequest{
		HTML:      html,
		TimeoutMs: 30_000,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("pdf: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		r.baseURL+"/render",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("pdf: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if r.token != "" {
		req.Header.Set("X-Internal-Token", r.token)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, ErrRendererUnavailable
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf(
			"%w: status=%d body=%s",
			ErrRenderFailed,
			resp.StatusCode,
			string(errBody),
		)
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("pdf: read response: %w", err)
	}

	if len(pdfBytes) == 0 {
		return nil, ErrRenderFailed
	}

	return pdfBytes, nil
}
