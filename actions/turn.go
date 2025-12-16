package actions

import (
	"context"
)

type PageTurnDirection int

const (
	TurnDirectionForward PageTurnDirection = iota
	TurnDirectionBackward
)

func (p PageTurnDirection) String() string {
	return [...]string{"Forward", "Backward"}[p]
}

type PageTurnerProvider interface {
	TurnPage(ctx context.Context, direction PageTurnDirection) error
	Close() error
}
