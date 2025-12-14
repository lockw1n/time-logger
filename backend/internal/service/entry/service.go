package entry

import (
	"context"

	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
)

type Service interface {
	Create(ctx context.Context, req entrydto.Request) (*entrydto.Response, error)
	Update(ctx context.Context, id uint64, req entrydto.Request) (*entrydto.Response, error)
	Delete(ctx context.Context, id uint64) error
}
