package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	myussd "github.com/yeboahnanaosei/x/testussd"
)

type CaseService struct {
	DB *redis.Client
}

func (cs *CaseService) GetCase(ctx context.Context, u myussd.User) (myussd.Case, error) {
	data, err := cs.DB.HGetAll(ctx, u.WaId+":case").Result()
	return myussd.Case{
		VictimName: data["name"],
		Gender:     data["gender"],
		Status:     data["status"],
		User:       u,
	}, err
}

func (cs *CaseService) SetCase(ctx context.Context, c myussd.Case) error {
	cs.DB.HSet(
		ctx,
		c.WaId+":case",
		"name", c.VictimName,
		"gender", c.Gender,
		"status", c.Status,
		"date", time.Now().Format("2006-01-02"),
	)
	return nil
}

func (cs *CaseService) ClearCase(ctx context.Context, c myussd.Case) error {
	_, err := cs.DB.Del(ctx, c.User.WaId+":case").Result()
	return err
}
