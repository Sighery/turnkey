package actions

import (
	"fmt"
	"testing"
)

func TestNewScreenOrientationFromKindleValue(t *testing.T) {
	testCases := []struct {
		val  string
		want ScreenOrientation
	}{
		{val: "U", want: OrientationPortrait},
		{val: "R", want: OrientationLandscape},
		{val: "D", want: OrientationRPortrait},
		{val: "L", want: OrientationRLandscape},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("val:%s,want:%d", tc.val, tc.want), func(t *testing.T) {
			got, err := newScreenOrientationFromKindleValue(tc.val)
			if err != nil {
				t.Error(err)
			}
			if got != tc.want {
				t.Errorf("Invalid. Want: %d, got: %d", tc.want, got)
			}
		})
	}
}

func TestNewScreenOrientationFromKindleValueInvalid(t *testing.T) {
	testCases := []string{"A", "G", "F", "Testing"}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc), func(t *testing.T) {
			_, err := newScreenOrientationFromKindleValue(tc)
			if err == nil {
				t.Errorf("Invalid value %s passed", tc)
			}
		})
	}
}

func TestScreenOrientationToKindleValue(t *testing.T) {
	testCases := []struct {
		val  ScreenOrientation
		want string
	}{
		{val: OrientationPortrait, want: "U"},
		{val: OrientationLandscape, want: "R"},
		{val: OrientationRPortrait, want: "D"},
		{val: OrientationRLandscape, want: "L"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("val:%d,want:%s", tc.val, tc.want), func(t *testing.T) {
			got, err := tc.val.toKindleValue()
			if err != nil {
				t.Error(err)
			}
			if got != tc.want {
				t.Errorf("Invalid. Want: %s, got: %s", tc.want, got)
			}
		})
	}
}
