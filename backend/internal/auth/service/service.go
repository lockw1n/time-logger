package service

import "context"

type Service interface {
	Login(ctx context.Context, in LoginInput) (LoginResult, error)
}
