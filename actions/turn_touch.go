package actions

import (
	"context"
	"fmt"
	"log/slog"
	"syscall"
	"time"

	"github.com/holoplot/go-evdev"
)

func currentTimeval() syscall.Timeval {
	now := time.Now()
	return syscall.Timeval{
		Sec:  int32(now.Unix()),
		Usec: int32(now.UnixNano() / 1000 % 1000),
	}
}

// For the Y coordinate, start seems to be top, while end is bottom
// For the X coordinate, start seems to be left, while end is right
type inputDeviceRange struct {
	start int32
	end   int32
}

func calculatePoint(start int32, end int32, percentage float32) int32 {
	return int32(((float32(end) - float32(start)) * (percentage / 100.0)) + float32(start))
}

type coordinate struct {
	x int32
	y int32
}

func calculateCoordinate(
	xRange, yRange inputDeviceRange, xPercentage, yPercentage float32,
	orientation ScreenOrientation,
) coordinate {
	var px, py float32
	switch orientation {
	case OrientationPortrait:
		px = xPercentage
		py = yPercentage
	case OrientationLandscape:
		px = yPercentage
		py = 100.0 - xPercentage
	case OrientationRPortrait:
		px = 100.0 - xPercentage
		py = 100.0 - yPercentage
	case OrientationRLandscape:
		px = 100.0 - yPercentage
		py = xPercentage
	}

	xCoor := calculatePoint(xRange.start, xRange.end, px)
	yCoor := calculatePoint(yRange.start, yRange.end, py)
	return coordinate{x: xCoor, y: yCoor}
}

type InputDevicePageTurnerProvider struct {
	device              *evdev.InputDevice
	rangeX              inputDeviceRange
	rangeY              inputDeviceRange
	orientationProvider OrientationProvider
}

func NewInputDevicePageTurnerProvider(
	device TouchscreenInputDevice, oProvider OrientationProvider,
) (InputDevicePageTurnerProvider, error) {
	evDev, err := evdev.Open(device.path)
	if err != nil {
		return InputDevicePageTurnerProvider{}, err
	}

	return InputDevicePageTurnerProvider{
		device:              evDev,
		rangeX:              device.rangeX,
		rangeY:              device.rangeY,
		orientationProvider: oProvider,
	}, nil
}

func (p InputDevicePageTurnerProvider) Close() error {
	return p.device.Close()
}

func (device InputDevicePageTurnerProvider) touchCoordinate(xCoor int32, yCoor int32) error {
	events := []evdev.InputEvent{
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_MT_SLOT,
			Value: 0,
		},
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_MT_TRACKING_ID,
			Value: 0,
		},
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_MT_POSITION_X,
			Value: xCoor,
		},
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_MT_POSITION_Y,
			Value: yCoor,
		},
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_SYN,
			Code:  evdev.SYN_REPORT,
			Value: 0,
		},
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_MT_TRACKING_ID,
			Value: -1,
		},
		evdev.InputEvent{
			Time:  currentTimeval(),
			Type:  evdev.EV_SYN,
			Code:  evdev.SYN_REPORT,
			Value: 0,
		},
	}

	for _, event := range events {
		err := device.device.WriteOne(&event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p InputDevicePageTurnerProvider) TurnPage(
	ctx context.Context, direction PageTurnDirection,
) error {
	orientation, err := p.orientationProvider.CurrentOrientation(ctx)
	if err != nil {
		return err
	}

	var xPer, yPer float32
	switch direction {
	case TurnDirectionForward:
		xPer = 95
		yPer = 50
	case TurnDirectionBackward:
		xPer = 5
		yPer = 50
	default:
		return fmt.Errorf("Unknown page turn direction %d", direction)
	}

	coor := calculateCoordinate(p.rangeX, p.rangeY, xPer, yPer, orientation)
	slog.Warn("Touching coordinate", "xCoor", coor.x, "yCoor", coor.y)
	return p.touchCoordinate(coor.x, coor.y)
}
