package human

import "testing"

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
