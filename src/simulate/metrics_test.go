package simulate

import (
	"testing"
)

func TestOutbreakDuration(t *testing.T) {

	params := Parameters{N: 1000}

	for _, test := range []struct {
		Is   []float64
		want float64
	}{
		{Is: []float64{}, want: 0},
		// no outbreak (5% threshold)
		{Is: []float64{1, 5, 10, 49, 11, 5, 0}, want: 0},
		{Is: []float64{1, 3, 100, 0}, want: 1},
		{Is: []float64{1, 3, 100, 100, 99, 50, 3}, want: 4},
		// Test that if we end before dropping below the threshold we properly terminate
		{Is: []float64{100}, want: 1},
	} {
		got := computeOutbreakDuration(test.Is, params)
		if got != test.want {
			t.Fatalf("computeOutbreakDuration(%v) = %v; want %v",
				test.Is, got, test.want)
		}

	}
}

func TestPeakTime(t *testing.T) {

	params := Parameters{N: 1000}

	for _, test := range []struct {
		Is   []float64
		want float64
	}{
		{Is: []float64{}, want: 0},
		// no outbreak (5% threshold)
		{Is: []float64{1, 5, 10, 49, 11, 5, 0}, want: 3},
		{Is: []float64{1, 3, 100, 0}, want: 2},
		{Is: []float64{1, 3, 100, 100, 99, 50, 3}, want: 2},
		// Test that if we end before dropping below the threshold we properly terminate
		{Is: []float64{1, 100}, want: 1},
	} {
		got := computePeakTime(test.Is, params)
		if got != test.want {
			t.Fatalf("computePeakTime(%v) = %v; want %v",
				test.Is, got, test.want)
		}

	}
}
