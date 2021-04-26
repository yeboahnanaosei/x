package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	myussd "github.com/yeboahnanaosei/x/testussd"
)

type UserService struct {
	DB *redis.Client
}

func (s *UserService) SetStage(ctx context.Context, u myussd.User, stage string) error {
	_, err := s.DB.Set(ctx, u.WaId, stage, time.Minute).Result()
	return err
}

func (s *UserService) GetStage(ctx context.Context, u myussd.User) (string, error) {
	res, err := s.DB.Get(ctx, u.WaId).Result()
	if errors.Is(err, redis.Nil) {
		return res, myussd.ErrStageNotFound
	}
	return res, err
}

func (s *UserService) ClearStage(ctx context.Context, u myussd.User) error {
	_, err := s.DB.Del(ctx, u.WaId).Result()
	return err
}
