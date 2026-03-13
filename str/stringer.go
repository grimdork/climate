package str

import (
	"strconv"
	"strings"
)

// Stringer extends strings.Builder with varargs-based write methods.
type Stringer struct {
	strings.Builder
	sliceComma bool
	mapComma   bool
	comma      byte
	equals     byte
}

// NewStringer returns a new Stringer with the default settings.
// - Comma separator ","
// - Equals symbol "="
func NewStringer() *Stringer {
	s := Stringer{
		comma:  ',',
		equals: '=',
	}
	return &s
}

// SetSliceComma enables adding a comma between elements in supplied slices.
func (s *Stringer) SetSliceComma(b bool) *Stringer {
	s.sliceComma = b
	return s
}

// SetMapComma enables adding a comma between elements in supplied maps.
func (s *Stringer) SetMapComma(b bool) *Stringer {
	s.mapComma = b
	return s
}

// SetComma sets the symbol to use for joining slices and map k-v pairs.
func (s *Stringer) SetComma(c byte) *Stringer {
	s.comma = c
	return s
}

// SetEquals sets the symbol to join keys and values in maps.
func (s *Stringer) SetEquals(e byte) *Stringer {
	s.equals = e
	return s
}

// WriteStrings writes any number of strings in one go.
func (s *Stringer) WriteStrings(v ...string) (int, error) {
	var err error
	var size, c int

	for _, x := range v {
		c, err = s.WriteString(x)
		if err != nil {
			return size, err
		}
		size += c
	}
	return size, nil
}

// WriteI writes any number of different types at once.
// Special notes:
// Integers - int and int64 are the only approved types, and untyped whole numbers will be parsed as int
// Floating point numbers - all numbers with a decimal point are interpreted as float64, with the fewest necessary decimal places
// Maps & slices - commas are not on by default, and maps will have "=" between each key-value pair
func (s *Stringer) WriteI(v ...any) (int, error) {
	var err error
	var size, c int

	for _, x := range v {
		c, err = s.writeX(x)
		size += c
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// writeX handles concrete types without reflect.
func (s *Stringer) writeX(x any) (int, error) {
	switch v := x.(type) {
	case bool:
		if v {
			return s.WriteString("true")
		}
		return s.WriteString("false")

	case string:
		return s.WriteString(v)

	case int:
		return s.WriteString(strconv.FormatInt(int64(v), 10))

	case int64:
		return s.WriteString(strconv.FormatInt(v, 10))

	case float64:
		return s.WriteString(strconv.FormatFloat(v, 'f', -1, 64))

	case []bool:
		return s.writeSlice(len(v), func(i int) any { return v[i] })
	case []string:
		return s.writeSlice(len(v), func(i int) any { return v[i] })
	case []int:
		return s.writeSlice(len(v), func(i int) any { return v[i] })
	case []int64:
		return s.writeSlice(len(v), func(i int) any { return v[i] })
	case []float64:
		return s.writeSlice(len(v), func(i int) any { return v[i] })
	case []any:
		return s.writeSlice(len(v), func(i int) any { return v[i] })

	case map[string]string:
		return s.writeMap(len(v), mapIter(v))
	case map[string]int:
		return s.writeMap(len(v), mapIter(v))
	case map[string]any:
		return s.writeMap(len(v), mapIter(v))
	case map[int]string:
		return s.writeMap(len(v), mapIter(v))
	case map[int]int:
		return s.writeMap(len(v), mapIter(v))
	case map[int]any:
		return s.writeMap(len(v), mapIter(v))

	default:
		return 0, nil
	}
}

func (s *Stringer) writeSlice(length int, get func(int) any) (int, error) {
	var size, c int
	var err error
	for i := 0; i < length; i++ {
		c, err = s.writeX(get(i))
		if err != nil {
			return size, err
		}
		size += c
		if s.sliceComma && i < length-1 {
			err = s.WriteByte(s.comma)
			if err != nil {
				return size, err
			}
			size++
		}
	}
	return size, nil
}

type mapEntry struct {
	key, val any
}

func mapIter[K comparable, V any](m map[K]V) []mapEntry {
	entries := make([]mapEntry, 0, len(m))
	for k, v := range m {
		entries = append(entries, mapEntry{k, v})
	}
	return entries
}

func (s *Stringer) writeMap(length int, entries []mapEntry) (int, error) {
	var size, c int
	var err error
	for i, e := range entries {
		c, err = s.writeX(e.key)
		if err != nil {
			return size, err
		}
		size += c

		err = s.WriteByte(s.equals)
		if err != nil {
			return size, err
		}
		size++

		c, err = s.writeX(e.val)
		if err != nil {
			return size, err
		}
		size += c

		if s.mapComma && i < length-1 {
			err = s.WriteByte(s.comma)
			if err != nil {
				return size, err
			}
			size++
		}
	}
	return size, nil
}