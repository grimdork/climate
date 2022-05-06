package str

import (
	"reflect"
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
// - Comma separator ) ","
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
func (s *Stringer) WriteI(v ...interface{}) (int, error) {
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

// writeX can recurse deeply.
func (s *Stringer) writeX(x interface{}) (int, error) {
	var err error
	var size, c int
	switch reflect.TypeOf(x).Kind() {
	case reflect.Bool:
		t := reflect.ValueOf(x).Interface().(bool)
		if t {
			c, err = s.WriteString("true")
		} else {
			c, err = s.WriteString("false")
		}
	case reflect.String:
		c, err = s.WriteString(x.(string))
	case reflect.Int:
		c, err = s.WriteString(strconv.FormatInt(int64(x.(int)), 10))
	case reflect.Int64:
		c, err = s.WriteString(strconv.FormatInt(x.(int64), 10))
	case reflect.Float64:
		c, err = s.WriteString(strconv.FormatFloat(x.(float64), 'f', -1, 64))
	case reflect.Slice:
		v := reflect.ValueOf(x)
		for i := 0; i < v.Len(); i++ {
			c, err = s.writeX(v.Index(i).Interface())
			if err != nil {
				return size, err
			}
			size += c
			c = 0
			if s.sliceComma && i < (v.Len()-1) {
				err = s.WriteByte(s.comma)
				if err != nil {
					return size, err
				}
				size++
			}
		}
	case reflect.Map:
		v := reflect.ValueOf(x)
		m := v.MapRange()
		mapsize := len(v.MapKeys())
		i := 0
		for m.Next() {
			c, err = s.writeX(m.Key().Interface())
			if err != nil {
				return size, err
			}
			size += c
			err = s.WriteByte(s.equals)
			if err != nil {
				return size, err
			}
			size++
			c, err = s.writeX(m.Value().Interface())
			if err != nil {
				return size, err
			}
			size += c
			c = 0
			if s.mapComma && i < (mapsize-1) {
				err = s.WriteByte(s.comma)
				if err != nil {
					return size, err
				}
			}
			size++
			i++
		}
	default:
	}
	size += c
	return size, err
}
