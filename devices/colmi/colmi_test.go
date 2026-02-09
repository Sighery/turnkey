package colmi

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMakePacket(t *testing.T) {
	testCases := []struct {
		command byte
		data    []byte
		want    []byte
	}{
		// BlinkTwiceRequest
		{
			15, []byte{}, []byte{
				0xF, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xF,
			},
		},
		// CameraRequest
		{
			2, []byte{byte(ActionEnableCameraGesture)}, []byte{
				0x2, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
			},
		},
		// CameraResponse
		{
			2, []byte{byte(ActionTakePhoto)}, []byte{
				0x2, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4,
			},
		},
	}

	for _, tc := range testCases {
		command := tc.command
		data := tc.data
		want := tc.want

		t.Run(fmt.Sprintf("%x_%x_%x", command, data, want), func(t *testing.T) {
			got, err := MakePacket(command, data)
			if err != nil {
				t.Error(err)
			}

			if !bytes.Equal(got, want) {
				t.Errorf("Not expected bytes %x", got)
			}
		})
	}
}

func TestParseAccelerometerData(t *testing.T) {
	testCases := []struct {
		data []byte
		want SensorAccelerometerResponse
	}{
		{
			[]byte{0xa1, 0x03, 0x04, 0x0d, 0x14, 0x09, 0xe7, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xba},
			SensorAccelerometerResponse{X: -399, Y: 77, Z: 329},
		},
	}

	for _, tc := range testCases {
		data := tc.data
		want := tc.want
		t.Run(fmt.Sprintf("% x_%+v", data, want), func(t *testing.T) {
			got := ParseAccelerometerData(data)
			if got != want {
				t.Errorf("Not expected result: %+v", got)
			}
		})
	}
}

func TestParseSpo2Data(t *testing.T) {
	testCases := []struct {
		data []byte
		want SensorSpo2Response
	}{
		{
			[]byte{0xa1, 0x01, 0x06, 0x45, 0x05, 0xe1, 0x00, 0x00, 0x02, 0xf0, 0x01, 0x00, 0x00, 0x00, 0x00, 0xc6},
			SensorSpo2Response{Spo2: 1605, Spo2Max: 225, Spo2Min: 0, Spo2Diff: 240},
		},
	}

	for _, tc := range testCases {
		data := tc.data
		want := tc.want
		t.Run(fmt.Sprintf("% x_%+v", data, want), func(t *testing.T) {
			got := ParseSpo2Response(data)
			if got != want {
				t.Errorf("Not expected result: %+v", got)
			}
		})
	}
}

func TestParsePpgData(t *testing.T) {
	testCases := []struct {
		data []byte
		want SensorPpgResponse
	}{
		{
			[]byte{0xa1, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4},
			SensorPpgResponse{Ppg: 1, PpgMax: 0, PpgMin: 0, PpgDiff: 0},
		},
	}

	for _, tc := range testCases {
		data := tc.data
		want := tc.want
		t.Run(fmt.Sprintf("% x_%+v", data, want), func(t *testing.T) {
			got := ParsePpgResponse(data)
			if got != want {
				t.Errorf("Not expected result: %+v", got)
			}
		})
	}
}
