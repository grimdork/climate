package ini

import (
	"math"
	"testing"
)

func TestNumericParsing(t *testing.T) {
	data := `

i=42
f=3.14
exp=1e6
neg=-2.5
zero=0
one=1
notnum=1.2.3

[sec]
si=100
sf=2.5
sci=6e3
snot=12x
`

	ini, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Dump properties for debugging
	t.Logf("properties: %+v", ini.Properties)
	for k, v := range ini.Properties {
		t.Logf("prop %s: Type=%d Value=%q intV=%d floatV=%v boolV=%v", k, v.Type, v.Value, v.intV, v.floatV, v.boolV)
	}

	if got := ini.GetInt("", "i"); got != 42 {
		t.Fatalf("expected i=42, got %d", got)
	}

	if got := ini.GetFloat("", "f"); math.Abs(got-3.14) > 1e-9 {
		t.Fatalf("expected f=3.14, got %v", got)
	}

	if got := ini.GetFloat("", "exp"); math.Abs(got-1e6) > 1e-6 {
		t.Fatalf("expected exp=1e6, got %v", got)
	}

	if got := ini.GetFloat("", "neg"); math.Abs(got+2.5) > 1e-9 {
		t.Fatalf("expected neg=-2.5, got %v", got)
	}

	if got := ini.GetInt("", "zero"); got != 0 {
		t.Fatalf("expected zero=0, got %d", got)
	}

	if got := ini.GetInt("", "one"); got != 1 {
		t.Fatalf("expected one=1, got %d", got)
	}

	if got := ini.GetString("", "notnum"); got != "1.2.3" {
		t.Fatalf("expected notnum to remain string, got %q", got)
	}

	// section values
	if got := ini.GetInt("sec", "si"); got != 100 {
		t.Fatalf("expected sec.si=100, got %d", got)
	}
	if got := ini.GetFloat("sec", "sf"); math.Abs(got-2.5) > 1e-9 {
		t.Fatalf("expected sec.sf=2.5, got %v", got)
	}
	if got := ini.GetFloat("sec", "sci"); math.Abs(got-6000) > 1e-6 {
		t.Fatalf("expected sec.sci=6e3, got %v", got)
	}
	if got := ini.GetString("sec", "snot"); got != "12x" {
		t.Fatalf("expected sec.snot=12x, got %q", got)
	}
}

func TestUnmarshalEmptyInputReturnsError(t *testing.T) {
	if _, err := Unmarshal(""); err == nil {
		t.Fatalf("expected error for empty input, got nil")
	}
}
