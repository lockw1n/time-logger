package service

import "time"

type clock struct{}

func NewClock() Clock {
	return &clock{}
}

func (clock) Now() time.Time {
	return time.Now().UTC()
}
