package util

import (
	"math"
	"testing"
)

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"314096.11", 314096.11},
		{"-148.25", -148.25},
		{"0", 0},
		{"", 0},
		{"not-a-number", 0},
		{"1.0982", 1.0982},
		{"-0.046", -0.046},
		{"0.0000693858", 0.0000693858},
	}
	for _, tt := range tests {
		got := ParseFloat(tt.input)
		if math.Abs(got-tt.want) > 1e-10 {
			t.Errorf("ParseFloat(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
