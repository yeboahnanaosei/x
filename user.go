package testussd

import (
	"context"
	"errors"
)

var ErrStageNotFound = errors.New("stage not found")

type User struct {
	ProfileName string
	WaId        string
	LastStage       string
	Input   string
}

type UserService interface {
	SetStage(ctx context.Context, u User, stage string) error
	GetStage(ctx context.Context, u User) (string, error)
	ClearStage(ctx context.Context, u User) error
}
