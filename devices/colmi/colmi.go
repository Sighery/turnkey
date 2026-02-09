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
	CommandRawSensor  byte = 161
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

type RawSensorAction int

const (
	_ RawSensorAction = iota
	_
	ActionDisableRawSensor
	_
	ActionEnableRawSensor
)

type CommandResponse interface {
	Type() byte
}

type CameraResponse struct {
	Action CameraAction
}

func (CameraResponse) Type() byte { return CommandCamera }

type SensorSpo2Response struct {
	Spo2     uint16
	Spo2Max  uint8
	Spo2Min  uint8
	Spo2Diff uint8
}

func (SensorSpo2Response) Type() byte { return CommandRawSensor }

type SensorPpgResponse struct {
	Ppg     uint16
	PpgMax  uint16
	PpgMin  uint16
	PpgDiff uint16
}

func (SensorPpgResponse) Type() byte { return CommandRawSensor }

type SensorAccelerometerResponse struct {
	X int16
	Y int16
	Z int16
}

func (SensorAccelerometerResponse) Type() byte { return CommandRawSensor }

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

func ParseCameraResponse(data []byte) CameraResponse {
	return CameraResponse{Action: CameraAction(data[1])}
}

func ParseSpo2Response(data []byte) SensorSpo2Response {
	return SensorSpo2Response{
		Spo2:     ((uint16)(data[2]) << 8) | uint16(data[3]),
		Spo2Max:  (uint8)(data[5]),
		Spo2Min:  (uint8)(data[7]),
		Spo2Diff: (uint8)(data[9]),
	}
}

func ParsePpgResponse(data []byte) SensorPpgResponse {
	return SensorPpgResponse{
		Ppg:     ((uint16)(data[2]) << 8) | uint16(data[3]),
		PpgMax:  ((uint16)(data[4]) << 8) | uint16(data[5]),
		PpgMin:  ((uint16)(data[6]) << 8) | uint16(data[7]),
		PpgDiff: ((uint16)(data[8]) << 8) | uint16(data[9]),
	}
}

func ParseAccelerometerData(data []byte) SensorAccelerometerResponse {
	x := int16((int16(data[6]) << 4) | (int16(data[7]) & 0x0F))
	if x&0x800 != 0 {
		x -= 1 << 12
	}

	y := int16((int16(data[2]) << 4) | (int16(data[3]) & 0x0F))
	if y&0x800 != 0 {
		y -= 1 << 12
	}

	z := int16((int16(data[4]) << 4) | (int16(data[5]) & 0x0F))
	if z&0x800 != 0 {
		z -= 1 << 12
	}

	return SensorAccelerometerResponse{X: x, Y: y, Z: z}
}

func ParseCommandResponse(data []byte) (CommandResponse, error) {
	if !CheckCrc(data) {
		return CameraResponse{}, fmt.Errorf("invalid length or crc: % x", data)
	}

	switch {
	case data[0] == CommandCamera:
		return ParseCameraResponse(data), nil
	case data[0] == CommandRawSensor && data[1] == 0x01:
		return ParseSpo2Response(data), nil
	case data[0] == CommandRawSensor && data[1] == 0x02:
		return ParsePpgResponse(data), nil
	case data[0] == CommandRawSensor && data[1] == 0x03:
		return ParseAccelerometerData(data), nil
	}

	return CameraResponse{}, fmt.Errorf("Command not yet implemented: % x", data)
}

func RawSensorsRequest(action RawSensorAction) []byte {
	packet, err := MakePacket(CommandRawSensor, []byte{byte(action)})
	if err != nil {
		panic(fmt.Sprintf("Not supposed to reach this: %s", err))
	}
	return packet
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
