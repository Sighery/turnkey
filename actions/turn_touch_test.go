package actions

import (
	"fmt"
	"testing"
)

func TestCalculatePoint(t *testing.T) {
	testCases := []struct {
		start      int32
		end        int32
		percentage float32
		want       int32
	}{
		{start: 0, end: 1695, percentage: 50, want: 847},
		{start: 0, end: 1271, percentage: 90, want: 1143},
		{start: 200, end: 1000, percentage: 50, want: 600},
		{start: 200, end: 1000, percentage: 55.5, want: 644},
	}
	for _, tc := range testCases {
		testString := fmt.Sprintf(
			"start: %d, end: %d, percentage: %f, want: %d",
			tc.start, tc.end, tc.percentage, tc.want,
		)
		t.Run(testString, func(t *testing.T) {
			got := calculatePoint(tc.start, tc.end, tc.percentage)
			if got != tc.want {
				t.Errorf("Invalid. Want: %d, got: %d", tc.want, got)
			}
		})
	}
}

func TestCalculateCoordinate(t *testing.T) {
	kindleX := inputDeviceRange{start: 0, end: 1271}
	kindleY := inputDeviceRange{start: 0, end: 1695}

	testCases := []struct {
		xRange inputDeviceRange
		yRange inputDeviceRange
		xPer   float32
		yPer   float32
		ori    ScreenOrientation
		want   coordinate
	}{
		// Forward direction
		{
			xRange: kindleX, yRange: kindleY, xPer: 95, yPer: 50, ori: OrientationPortrait,
			want: coordinate{x: 1207, y: 847},
		},
		{
			xRange: kindleX, yRange: kindleY, xPer: 95, yPer: 50, ori: OrientationLandscape,
			want: coordinate{x: 635, y: 84},
		},
		{
			xRange: kindleX, yRange: kindleY, xPer: 95, yPer: 50, ori: OrientationRPortrait,
			want: coordinate{x: 63, y: 847},
		},
		{
			xRange: kindleX, yRange: kindleY, xPer: 95, yPer: 50, ori: OrientationRLandscape,
			want: coordinate{x: 635, y: 1610},
		},
		// Backward direction
		{
			xRange: kindleX, yRange: kindleY, xPer: 5, yPer: 50, ori: OrientationPortrait,
			want: coordinate{x: 63, y: 847},
		},
		{
			xRange: kindleX, yRange: kindleY, xPer: 5, yPer: 50, ori: OrientationLandscape,
			want: coordinate{x: 635, y: 1610},
		},
		{
			xRange: kindleX, yRange: kindleY, xPer: 5, yPer: 50, ori: OrientationRPortrait,
			want: coordinate{x: 1207, y: 847},
		},
		{
			xRange: kindleX, yRange: kindleY, xPer: 5, yPer: 50, ori: OrientationRLandscape,
			want: coordinate{x: 635, y: 84},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := calculateCoordinate(tc.xRange, tc.yRange, tc.xPer, tc.yPer, tc.ori)
			if got != tc.want {
				t.Errorf("Invalid. Wanted %+v, got %+v", tc.want, got)
			}
		})
	}
}
