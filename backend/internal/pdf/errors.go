package pdf

import "errors"

var (
	ErrEmptyHTML           = errors.New("pdf: empty html")
	ErrRenderFailed        = errors.New("pdf: render failed")
	ErrRendererUnavailable = errors.New("pdf: renderer unavailable")
)
