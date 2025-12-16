package actions

import (
	"context"
	"fmt"

	"github.com/Sighery/turnkey/services/lipc"
)

const (
	WinmgrService   = "com.lab126.winmgr"
	OrientationProp = "orientationLock"
)

type KindleOrientation string

const (
	KOrientationUnlocked   KindleOrientation = ""
	KOrientationPortrait   KindleOrientation = "U"
	KOrientationLandscape  KindleOrientation = "R"
	KOrientationRPortrait  KindleOrientation = "D"
	KOrientationRLandscape KindleOrientation = "L"
)

type ScreenOrientation int

const (
	OrientationPortrait   ScreenOrientation = 0
	OrientationLandscape  ScreenOrientation = 90
	OrientationRPortrait  ScreenOrientation = 180
	OrientationRLandscape ScreenOrientation = 270
)

func (s ScreenOrientation) String() string {
	return [...]string{
		"Portrait", "Landscape", "Reverse Portrait", "Reverse Landscape",
	}[(s/90)-1]
}

func newScreenOrientationFromKindleValue(val string) (ScreenOrientation, error) {
	switch val {
	case "U":
		return OrientationPortrait, nil
	case "R":
		return OrientationLandscape, nil
	case "D":
		return OrientationRPortrait, nil
	case "L":
		return OrientationRLandscape, nil
	default:
		return OrientationPortrait, fmt.Errorf("Unknown orientation %s", val)
	}
}

func (orientation ScreenOrientation) toKindleValue() (string, error) {
	switch orientation {
	case OrientationPortrait:
		return "U", nil
	case OrientationLandscape:
		return "R", nil
	case OrientationRPortrait:
		return "D", nil
	case OrientationRLandscape:
		return "L", nil
	default:
		return "", fmt.Errorf("Unknown orientation %d", orientation)
	}
}

type OrientationProvider interface {
	CurrentOrientation(ctx context.Context) (ScreenOrientation, error)
	SetOrientation(ctx context.Context, orientation ScreenOrientation) error
}

type LipcOrientationProvider struct {
	client *lipc.LipcClient
}

func NewLipcOrientationProvider(client *lipc.LipcClient) LipcOrientationProvider {
	return LipcOrientationProvider{client: client}
}

func (p LipcOrientationProvider) CurrentOrientation(
	ctx context.Context,
) (ScreenOrientation, error) {
	val, err := p.client.GetString(ctx, WinmgrService, OrientationProp)
	if err != nil {
		return OrientationPortrait, err
	}

	return newScreenOrientationFromKindleValue(val)
}

func (p LipcOrientationProvider) SetOrientation(
	ctx context.Context, orientation ScreenOrientation,
) error {
	lipcVal, err := orientation.toKindleValue()
	if err != nil {
		return err
	}

	err = p.client.SetString(ctx, WinmgrService, OrientationProp, lipcVal)
	if err != nil {
		return err
	}
	return nil
}
