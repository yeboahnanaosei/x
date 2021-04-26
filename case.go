package testussd

import (
	"context"
	"fmt"
)

type Case struct {
	VictimName string
	Gender     string
	Status     string
	Date       string
	User
}

func (c Case) String() string {
	return fmt.Sprintf("Name: %s\nGender: %s\nStatus: %s\n", c.VictimName, c.Gender, c.Status)
}

type CaseService interface {
	GetCase(ctx context.Context, u User) (Case, error)
	SetCase(ctx context.Context, c Case) error
	ClearCase(ctx context.Context, c Case) error
}
