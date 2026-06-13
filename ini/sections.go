package ini

import (
	"bufio"
	"strings"
	"sync"

	"github.com/grimdork/climate/str"
)

// Section holds one or more fields.
type Section struct {
	mu     sync.RWMutex
	Fields map[string][]*Field
	// Order fields were loaded or added in.
	Order []string
}

// parse section properties until a new section or end of file.
func (s *Section) parse(r *bufio.Reader) {
	for {
		next, err := r.Peek(2)
		if err != nil || len(next) == 0 {
			return
		}

		// Blank line — skip
		if next[0] == '\n' || next[0] == '\r' {
			r.ReadString('\n')
			continue
		}

		// New section — done
		if next[0] == '[' || (len(next) > 1 && next[1] == '[') {
			return
		}

		p, err := r.ReadString('\n')
		if err != nil {
			return
		}

		// Skip comments
		if strings.HasPrefix(p, "#") || strings.HasPrefix(p, ";") {
			continue
		}

		a := strings.SplitN(p, "=", 2)
		if len(a) != 2 {
			continue
		}

		a[0] = strings.TrimSpace(a[0])
		a[1] = strings.TrimSpace(a[1])
		a[0] = strings.ToLower(a[0])
		if _, ok := str.BoolFromString(a[1]); ok {
			_ = s.AddBool(a[0], boolValue(a[1]))
			continue
		}

		// Try ints and floats next.
		_ = s.AddString(a[0], a[1])
	}
}

// GetBool returns the first field as a bool, or the alternative.
func (s *Section) GetBool(key string, alt bool) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[strings.ToLower(key)]
	if !ok || len(v) == 0 {
		return alt
	}
	return v[0].boolV
}

// AddBool appends a new bool field to the section.
func (s *Section) AddBool(key string, value bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	f := Field{}
	f.SetBool(value)
	s.Fields[key] = append(s.Fields[key], &f)
	if len(s.Fields[key]) == 1 {
		s.Order = append(s.Order, key)
	}
	return nil
}

// SetBool replaces the first entry's value, or adds a new entry if the key does not exist.
func (s *Section) SetBool(key string, value bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	if len(s.Fields[key]) > 0 {
		s.Fields[key][0].SetBool(value)
		return nil
	}
	f := Field{}
	f.SetBool(value)
	s.Fields[key] = []*Field{&f}
	s.Order = append(s.Order, key)
	return nil
}

// GetInt returns the first field as an int64, or the alternative.
func (s *Section) GetInt(key string, alt int64) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[strings.ToLower(key)]
	if !ok || len(v) == 0 {
		return alt
	}
	return v[0].intV
}

// AddInt appends a new int64 field to the section.
func (s *Section) AddInt(key string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	f := Field{}
	f.SetInt(value)
	s.Fields[key] = append(s.Fields[key], &f)
	if len(s.Fields[key]) == 1 {
		s.Order = append(s.Order, key)
	}
	return nil
}

// SetInt replaces the first entry's value, or adds a new entry if the key does not exist.
func (s *Section) SetInt(key string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	if len(s.Fields[key]) > 0 {
		s.Fields[key][0].SetInt(value)
		return nil
	}
	f := Field{}
	f.SetInt(value)
	s.Fields[key] = []*Field{&f}
	s.Order = append(s.Order, key)
	return nil
}

// GetFloat returns the first field as a float64, or the alternative.
func (s *Section) GetFloat(key string, alt float64) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[strings.ToLower(key)]
	if !ok || len(v) == 0 {
		return alt
	}
	return v[0].floatV
}

// AddFloat appends a new float64 field to the section.
func (s *Section) AddFloat(key string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	f := Field{}
	f.SetFloat(value)
	s.Fields[key] = append(s.Fields[key], &f)
	if len(s.Fields[key]) == 1 {
		s.Order = append(s.Order, key)
	}
	return nil
}

// SetFloat replaces the first entry's value, or adds a new entry if the key does not exist.
func (s *Section) SetFloat(key string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	if len(s.Fields[key]) > 0 {
		s.Fields[key][0].SetFloat(value)
		return nil
	}
	f := Field{}
	f.SetFloat(value)
	s.Fields[key] = []*Field{&f}
	s.Order = append(s.Order, key)
	return nil
}

// GetString returns the first field as a string, or the alternative.
func (s *Section) GetString(key, alt string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[strings.ToLower(key)]
	if !ok || len(v) == 0 {
		return alt
	}
	return v[0].Value
}

// AddString appends a new string field to the section.
func (s *Section) AddString(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	f := Field{}
	f.SetString(value)
	s.Fields[key] = append(s.Fields[key], &f)
	if len(s.Fields[key]) == 1 {
		s.Order = append(s.Order, key)
	}
	return nil
}

// SetString replaces the first entry's value, or adds a new entry if the key does not exist.
func (s *Section) SetString(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	key = strings.ToLower(key)
	if len(s.Fields[key]) > 0 {
		s.Fields[key][0].SetString(value)
		return nil
	}
	f := Field{}
	f.SetString(value)
	s.Fields[key] = []*Field{&f}
	s.Order = append(s.Order, key)
	return nil
}

// GetAll returns all fields for a key, or nil if the key does not exist.
func (s *Section) GetAll(key string) []*Field {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[strings.ToLower(key)]
	if !ok {
		return nil
	}
	return v
}

// boolValue from common strings.
func boolValue(s string) bool {
	v, ok := str.BoolFromString(s)
	return ok && v
}

// BoolFromString returns true for truthy values.
func BoolFromString(s string) bool {
	return boolValue(s)
}
