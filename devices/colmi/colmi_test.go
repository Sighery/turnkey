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
