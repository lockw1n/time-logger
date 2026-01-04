package service

import "time"

type LoginResult struct {
	ConsultantID uint64
	AccessToken  string
	ExpiresAt    time.Time
}
