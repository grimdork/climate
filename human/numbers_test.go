package human

import (
	"math"
	"testing"
)

func TestUIntBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		si       bool
		expected string
	}{
		{0, false, "0 B"},
		{1, false, "1 B"},
		{512, false, "512 B"},
		{999, true, "999 B"},
	}
	for _, tt := range tests {
		result := UInt(tt.input, tt.si)
		if result != tt.expected {
			t.Errorf("UInt(%d, %v) = %q, want %q", tt.input, tt.si, result, tt.expected)
		}
	}
}

func TestUIntBinary(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{1024, "1 kiB"},
		{1536, "1.5 kiB"},
		{1048576, "1 MiB"},
		{1572864, "1.5 MiB"},
		{1073741824, "1 GiB"},
		{1099511627776, "1 TiB"},
	}
	for _, tt := range tests {
		result := UInt(tt.input, false)
		if result != tt.expected {
			t.Errorf("UInt(%d, false) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestUIntSI(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{1000, "1 kB"},
		{1500, "1.5 kB"},
		{1000000, "1 MB"},
		{1500000, "1.5 MB"},
		{1000000000, "1 GB"},
		{1000000000000, "1 TB"},
	}
	for _, tt := range tests {
		result := UInt(tt.input, true)
		if result != tt.expected {
			t.Errorf("UInt(%d, true) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestUIntWholeNumbers(t *testing.T) {
	// Whole numbers should not have .0
	result := UInt(1048576, false)
	if result != "1 MiB" {
		t.Errorf("expected '1 MiB', got %q", result)
	}

	result = UInt(1000000, true)
	if result != "1 MB" {
		t.Errorf("expected '1 MB', got %q", result)
	}
}

func TestFloatBinary(t *testing.T) {
	tests := []struct {
		input float64
		prec  int
		want  string
	}{
		{0, 0, "0 B"},
		{512, 0, "512 B"},
		{1024, 0, "1 kiB"},
		{1536, 1, "1.5 kiB"},
		{1048576, 0, "1 MiB"},
		{1572864, 1, "1.5 MiB"},
		{1073741824, 0, "1 GiB"},
	}
	for _, tc := range tests {
		got := Float(tc.input, tc.prec, false)
		if got != tc.want {
			t.Errorf("Float(%v, %d, false) = %q, want %q", tc.input, tc.prec, got, tc.want)
		}
	}
}

func TestFloatSI(t *testing.T) {
	tests := []struct {
		input float64
		prec  int
		want  string
	}{
		{0, 0, "0 B"},
		{500, 0, "500 B"},
		{1000, 0, "1 kB"},
		{1500, 1, "1.5 kB"},
		{1000000, 0, "1 MB"},
		{1500000, 1, "1.5 MB"},
		{1000000000, 0, "1 GB"},
	}
	for _, tc := range tests {
		got := Float(tc.input, tc.prec, true)
		if got != tc.want {
			t.Errorf("Float(%v, %d, true) = %q, want %q", tc.input, tc.prec, got, tc.want)
		}
	}
}

func TestFloatNaNInf(t *testing.T) {
	if got := Float(math.NaN(), 2, false); got != "NaN B" {
		t.Errorf("expected 'NaN B', got %q", got)
	}
	if got := Float(math.Inf(1), 2, false); got != "+Inf B" {
		t.Errorf("expected '+Inf B', got %q", got)
	}
	if got := Float(math.Inf(-1), 2, false); got != "-Inf B" {
		t.Errorf("expected '-Inf B', got %q", got)
	}
}

func TestFloatWholeNumberNoDecimal(t *testing.T) {
	if got := Float(2048, 2, false); got != "2 kiB" {
		t.Errorf("expected '2 kiB', got %q", got)
	}
	if got := Float(2000, 2, true); got != "2 kB" {
		t.Errorf("expected '2 kB', got %q", got)
	}
}

func TestFloatEdgeCases(t *testing.T) {
	// EiB threshold
	if got := Float(1152921504606846976, 0, false); got != "1 EiB" {
		t.Errorf("expected '1 EiB', got %q", got)
	}
	if got := Float(1e18, 0, true); got != "1 EB" {
		t.Errorf("expected '1 EB', got %q", got)
	}
}
