package actions

import (
	"fmt"
	"log/slog"

	"github.com/holoplot/go-evdev"
)

type TouchscreenInputDevice struct {
	name   string
	path   string
	rangeX inputDeviceRange
	rangeY inputDeviceRange
}

func FindTouchscreenInputDevice() (TouchscreenInputDevice, error) {
	devices, err := evdev.ListDevicePaths()
	if err != nil {
		return TouchscreenInputDevice{}, err
	}

	for _, deviceInfo := range devices {
		device, err := evdev.Open(deviceInfo.Path)
		if err != nil {
			slog.Warn(fmt.Sprintf("Couldn't open evdev %#v", deviceInfo))
			continue
		}
		defer device.Close()

		absInfo, err := device.AbsInfos()
		if err != nil || len(absInfo) == 0 {
			slog.Warn(fmt.Sprintf("Couldn't get ABS Info for evdev %#v", deviceInfo))
			continue
		}

		xData, xFound := absInfo[evdev.ABS_MT_POSITION_X]
		yData, yFound := absInfo[evdev.ABS_MT_POSITION_Y]

		if xFound == false || yFound == false {
			continue
		}

		return TouchscreenInputDevice{
			name:   deviceInfo.Name,
			path:   deviceInfo.Path,
			rangeX: inputDeviceRange{start: xData.Minimum, end: xData.Maximum},
			rangeY: inputDeviceRange{start: yData.Minimum, end: yData.Maximum},
		}, nil
	}

	return TouchscreenInputDevice{}, fmt.Errorf("Couldn't find a valid input device")
}
