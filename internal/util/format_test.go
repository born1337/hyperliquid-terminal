package util

import "testing"

func TestFormatUSD(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "$0.00"},
		{1.5, "$1.50"},
		{999.99, "$999.99"},
		{1000, "$1,000.00"},
		{1234567.89, "$1,234,567.89"},
		{-500, "-$500.00"},
		{-1234567.89, "-$1,234,567.89"},
		{0.01, "$0.01"},
		{99999999.99, "$99,999,999.99"},
	}
	for _, tt := range tests {
		got := FormatUSD(tt.input)
		if got != tt.want {
			t.Errorf("FormatUSD(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatSignedUSD(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{100, "+$100.00"},
		{-100, "-$100.00"},
		{0, "+$0.00"},
		{1234.56, "+$1,234.56"},
	}
	for _, tt := range tests {
		got := FormatSignedUSD(tt.input)
		if got != tt.want {
			t.Errorf("FormatSignedUSD(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatPercent(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{10.5, "+10.50%"},
		{-5.25, "-5.25%"},
		{0, "+0.00%"},
		{100, "+100.00%"},
	}
	for _, tt := range tests {
		got := FormatPercent(tt.input)
		if got != tt.want {
			t.Errorf("FormatPercent(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatFundingRate(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0.0001, "0.2400%"},
		{-0.00005, "-0.1200%"},
		{0, "0.0000%"},
	}
	for _, tt := range tests {
		got := FormatFundingRate(tt.input)
		if got != tt.want {
			t.Errorf("FormatFundingRate(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatLeverage(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{10, "10x"},
		{4.16, "4.16x"},
		{1, "1x"},
		{50, "50x"},
		{2.5, "2.50x"},
	}
	for _, tt := range tests {
		got := FormatLeverage(tt.input)
		if got != tt.want {
			t.Errorf("FormatLeverage(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{91000, "$91,000.00"},
		{100.50, "$100.50"},
		{0.0045, "$0.0045"},
		{1234.56, "$1,234.56"},
		{0.5, "$0.5000"},
	}
	for _, tt := range tests {
		got := FormatPrice(tt.input)
		if got != tt.want {
			t.Errorf("FormatPrice(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{1.0982, "1.0982"},
		{-1.0982, "1.0982"}, // should strip negative
		{0.00123, "0.001230"},
		{100, "100.0000"},
	}
	for _, tt := range tests {
		got := FormatSize(tt.input)
		if got != tt.want {
			t.Errorf("FormatSize(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatTime(t *testing.T) {
	// 2025-01-15 10:30:00 UTC = 1736935800000 ms
	got := FormatTime(1736935800000)
	// Just verify it doesn't panic and produces a reasonable format
	if len(got) == 0 {
		t.Error("FormatTime returned empty string")
	}
}
