// Taken from clintharrison's implementation:
// https://github.com/clintharrison/bueno/blob/aebf3574fc7e10937635195da5d38ca427b254b3/colmi/colmi.go
// Colmi ring documentation: https://colmi.puxtril.com/
package colmi

import (
	"fmt"
)

var (
	CommandServiceUUID = "6e40fff0-b5a3-f393-e0a9-e50e24dcca9e"
	CommandReadUUID    = "6e400003-b5a3-f393-e0a9-e50e24dcca9e"
	CommandWriteUUID   = "6e400002-b5a3-f393-e0a9-e50e24dcca9e"
)

var (
	CommandBlinkTwice byte = 16
	CommandCamera     byte = 2
)

type CameraAction int

const (
	_ CameraAction = iota
	ActionIntoCameraUI
	ActionTakePhoto
	ActionFinish
	ActionEnableCameraGesture
	ActionKeepScreenOn
	ActionDisableCameraGesture
)

func MakePacket(command byte, data []byte) ([]byte, error) {
	if len(data) > 14 {
		return nil, fmt.Errorf("Data can't be over 14 bytes")
	}

	crc := uint32(command)

	for _, b := range data {
		crc += uint32(b)
	}

	var packet [16]uint8
	packet[0] = uint8(command)
	packet[15] = uint8(crc & 0xFF)

	for i, v := range data {
		packet[i+1] = v
	}

	return packet[:], nil
}

func BlinkTwiceRequest() []byte {
	packet, err := MakePacket(CommandBlinkTwice, []byte{})
	if err != nil {
		panic(fmt.Sprintf("Not supposed to reach this: %s", err))
	}
	return packet
}

func CameraActionRequest(action CameraAction) []byte {
	packet, err := MakePacket(CommandCamera, []byte{byte(action)})
	if err != nil {
		panic(fmt.Sprintf("Not supposed to reach this: %s", err))
	}
	return packet
}

func CheckCrc(data []byte) bool {
	if len(data) != 16 {
		return false
	}

	var sum uint32
	for _, b := range data[:15] {
		sum += uint32(b)
	}

	if data[15] != uint8(sum&0xFF) {
		return false
	}

	return true
}

func CameraActionResponse(data []byte) (CameraAction, error) {
	if !CheckCrc(data) {
		return CameraAction(0), fmt.Errorf("invalid crc or length: % x", data)
	}

	if data[0] != CommandCamera {
		return CameraAction(0), fmt.Errorf("packet is not for camera data: %x", data[0])
	}

	action := CameraAction(data[1])

	return action, nil
}
