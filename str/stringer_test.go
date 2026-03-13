package str

import (
	"testing"
)

func TestWriteStrings(t *testing.T) {
	s := NewStringer()
	s.WriteStrings("hello", " ", "world")
	if s.String() != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", s.String())
	}
}

func TestWriteIScalars(t *testing.T) {
	s := NewStringer()
	s.WriteI("count:", 42, " pi:", 3.14, " ok:", true, " big:", int64(9999999999))
	expected := "count:42 pi:3.14 ok:true big:9999999999"
	if s.String() != expected {
		t.Fatalf("expected '%s', got '%s'", expected, s.String())
	}
}

func TestWriteIStringSlice(t *testing.T) {
	s := NewStringer().SetSliceComma(true)
	s.WriteI([]string{"a", "b", "c"})
	if s.String() != "a,b,c" {
		t.Fatalf("expected 'a,b,c', got '%s'", s.String())
	}
}

func TestWriteIStringSliceNoComma(t *testing.T) {
	s := NewStringer()
	s.WriteI([]string{"a", "b", "c"})
	if s.String() != "abc" {
		t.Fatalf("expected 'abc', got '%s'", s.String())
	}
}

func TestWriteIIntSlice(t *testing.T) {
	s := NewStringer().SetSliceComma(true)
	s.WriteI([]int{1, 2, 3})
	if s.String() != "1,2,3" {
		t.Fatalf("expected '1,2,3', got '%s'", s.String())
	}
}

func TestWriteIFloatSlice(t *testing.T) {
	s := NewStringer().SetSliceComma(true).SetComma(';')
	s.WriteI([]float64{1.5, 2.5})
	if s.String() != "1.5;2.5" {
		t.Fatalf("expected '1.5;2.5', got '%s'", s.String())
	}
}

func TestWriteIBoolSlice(t *testing.T) {
	s := NewStringer().SetSliceComma(true)
	s.WriteI([]bool{true, false, true})
	if s.String() != "true,false,true" {
		t.Fatalf("expected 'true,false,true', got '%s'", s.String())
	}
}

func TestWriteIInt64Slice(t *testing.T) {
	s := NewStringer().SetSliceComma(true)
	s.WriteI([]int64{100, 200})
	if s.String() != "100,200" {
		t.Fatalf("expected '100,200', got '%s'", s.String())
	}
}

func TestWriteIAnySlice(t *testing.T) {
	s := NewStringer().SetSliceComma(true)
	s.WriteI([]any{"x", 1, true})
	if s.String() != "x,1,true" {
		t.Fatalf("expected 'x,1,true', got '%s'", s.String())
	}
}

func TestWriteIMapStringString(t *testing.T) {
	s := NewStringer()
	s.WriteI(map[string]string{"key": "val"})
	if s.String() != "key=val" {
		t.Fatalf("expected 'key=val', got '%s'", s.String())
	}
}

func TestWriteIMapStringInt(t *testing.T) {
	s := NewStringer()
	s.WriteI(map[string]int{"cpu": 4})
	if s.String() != "cpu=4" {
		t.Fatalf("expected 'cpu=4', got '%s'", s.String())
	}
}

func TestWriteIMapCustomDelimiters(t *testing.T) {
	s := NewStringer().SetMapComma(true).SetEquals(':').SetComma('|')
	// Single entry map to avoid order issues
	s.WriteI(map[string]int{"ram": 16})
	if s.String() != "ram:16" {
		t.Fatalf("expected 'ram:16', got '%s'", s.String())
	}
}

func TestWriteIMapComma(t *testing.T) {
	s := NewStringer().SetMapComma(true)
	m := map[string]string{"a": "1"}
	s.WriteI(m, ",")
	m2 := map[string]string{"b": "2"}
	s.WriteI(m2)
	// With single-entry maps we can predict output
	if s.String() != "a=1,b=2" {
		t.Fatalf("expected 'a=1,b=2', got '%s'", s.String())
	}
}

func TestWriteIUnsupportedType(t *testing.T) {
	s := NewStringer()
	// Unsupported types should be silently skipped
	s.WriteI(struct{}{})
	if s.String() != "" {
		t.Fatalf("expected empty string for unsupported type, got '%s'", s.String())
	}
}

func TestChaining(t *testing.T) {
	s := NewStringer().SetSliceComma(true).SetMapComma(true).SetComma(';').SetEquals(':')
	s.WriteI([]string{"a", "b"})
	if s.String() != "a;b" {
		t.Fatalf("expected 'a;b', got '%s'", s.String())
	}
}