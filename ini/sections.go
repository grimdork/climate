package ini

import (
	"bufio"
	"strings"
	"sync"
)

// Section holds one or more fields.
type Section struct {
	mu     sync.RWMutex
	Fields map[string]*Field
	// Order fields were loaded or added in.
	Order []string
}

// parse section properties until a new section or end of file.
func (s *Section) parse(r *bufio.Reader) {
	loop := true
	for loop {
		next, err := r.Peek(2)
		// EOF
		if err != nil {
			return
		}

		// Skip blank lines
		if next[0] == '\n' {
			return
		}

		// New section, so this one's done
		if next[0] == '[' || next[1] == '[' {
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
		if a == nil || len(a) != 2 {
			return
		}

		a[0] = strings.TrimSpace(a[0])
		a[1] = strings.TrimSpace(a[1])
		switch a[1] {
		case "yes", "true", "on", "no", "false", "off":
			s.AddBool(a[0], boolValue(a[1]))
			continue
		}

		// Try ints and floats next.
		s.AddString(a[0], a[1])
	}
}

// GetBool returns a field as a bool, or the alternative.
func (s *Section) GetBool(key string, alt bool) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[key]
	if !ok {
		return alt
	}

	return v.boolV
}

// AddBool adds a new bool field to the section.
func (s *Section) AddBool(key string, value bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := Field{}
	f.SetBool(value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// GetInt returns a field as an int64, or the alternative.
func (s *Section) GetInt(key string, alt int64) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[key]
	if !ok {
		return alt
	}

	return v.intV
}

// AddInt adds a new int64 field to the section.
func (s *Section) AddInt(key string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := Field{}
	f.SetInt(value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// GetFloat returns a field as a float64, or the alternative.
func (s *Section) GetFloat(key string, alt float64) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[key]
	if !ok {
		return alt
	}

	return v.floatV
}

// AddFloat adds a new float64 field to the section.
func (s *Section) AddFloat(key string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := Field{}
	f.SetFloat(value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// GetString returns a field as a string, or the alternative.
func (s *Section) GetString(key, alt string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Fields[key]
	if !ok {
		return alt
	}

	return v.Value
}

// AddString adds a new string field to the section.
func (s *Section) AddString(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := Field{}
	f.SetString(value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// boolValue from common strings.
func boolValue(s string) bool {
	switch s {
	case "yes", "true", "1", "on":
		return true
	}

	return false
}

// BoolValue is a convenience that counts anything but "true", "on" or "enabled" as false.
func BoolValue(s string) bool {
	switch s {
	case "true", "on", "enabled", "yes":
		return true
	}
	return false
}
