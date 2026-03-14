package env

import (
	"os"
	"testing"
)

func TestGetAndFallback(t *testing.T) {
	os.Unsetenv("TEST_ENV_GET")
	if v := Get("TEST_ENV_GET", "fallback"); v != "fallback" {
		t.Fatalf("expected fallback, got %q", v)
	}

	os.Setenv("TEST_ENV_GET", "value")
	defer os.Unsetenv("TEST_ENV_GET")
	if v := Get("TEST_ENV_GET", "fallback"); v != "value" {
		t.Fatalf("expected value, got %q", v)
	}
	// Empty string is preserved when the variable is set.
	os.Setenv("TEST_ENV_GET", "")
	if v := Get("TEST_ENV_GET", "fallback"); v != "" {
		t.Fatalf("expected empty string, got %q", v)
	}
}

func TestGetBool(t *testing.T) {
	cases := []struct {
		val string
		alt bool
		exp bool
	}{
		{"", true, true},   // unset -> alt
		{"true", false, true},
		{"TRUE", false, true},
		{"yes", false, true},
		{"on", false, true},
		{"1", false, true},
		{"false", true, false},
		{"0", true, false},
		{"no", true, false},
		{"invalid", true, true}, // parse fails -> alt when set
		{"\"true\"", false, true}, // quoted true
	}

	for _, c := range cases {
		if c.val == "" {
			os.Unsetenv("TEST_ENV_BOOL")
		} else {
			os.Setenv("TEST_ENV_BOOL", c.val)
		}
		res := GetBool("TEST_ENV_BOOL", c.alt)
		if res != c.exp {
			t.Fatalf("GetBool(%q, %v) = %v; want %v", c.val, c.alt, res, c.exp)
		}
	}
	os.Unsetenv("TEST_ENV_BOOL")
}

func TestGetInt(t *testing.T) {
	cases := []struct {
		val string
		alt int64
		exp int64
	}{
		{"", 7, 7},
		{"42", 0, 42},
		{"0x2A", 0, 42},
		{"1_000", 0, 1000},
		{"-5", 0, -5},
		{"bad", 9, 9}, // parse fails -> alt when set
	}

	for _, c := range cases {
		if c.val == "" {
			os.Unsetenv("TEST_ENV_INT")
		} else {
			os.Setenv("TEST_ENV_INT", c.val)
		}
		res := GetInt("TEST_ENV_INT", c.alt)
		if res != c.exp {
			t.Fatalf("GetInt(%q, %d) = %d; want %d", c.val, c.alt, res, c.exp)
		}
	}
	os.Unsetenv("TEST_ENV_INT")
}

func TestGetFloat(t *testing.T) {
	cases := []struct {
		val string
		alt float64
		exp float64
	}{
		{"", 1.5, 1.5},
		{"3.14", 0, 3.14},
		{"3,14", 0, 3.14},
		{"1_000.5", 0, 1000.5},
		{"-2.5", 0, -2.5},
		{"bad", 2.25, 2.25},
	}

	for _, c := range cases {
		if c.val == "" {
			os.Unsetenv("TEST_ENV_FLOAT")
		} else {
			os.Setenv("TEST_ENV_FLOAT", c.val)
		}
		res := GetFloat("TEST_ENV_FLOAT", c.alt)
		if res != c.exp {
			t.Fatalf("GetFloat(%q, %v) = %v; want %v", c.val, c.alt, res, c.exp)
		}
	}
	os.Unsetenv("TEST_ENV_FLOAT")
}
