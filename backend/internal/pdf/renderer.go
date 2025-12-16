package pdf

import "context"

type Renderer interface {
	RenderHTML(ctx context.Context, html string) ([]byte, error)
}
