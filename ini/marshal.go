package ini

import (
	"errors"
	"strconv"
	"strings"

	"github.com/grimdork/climate/str"
)

// Marshal INI into a string. Global properties first (in insertion order), then sections (in insertion order).
func (ini *INI) Marshal() string {
	ini.mu.RLock()
	defer ini.mu.RUnlock()

	b := &strings.Builder{}

	// Properties
	if len(ini.Properties) > 0 {
		for _, k := range ini.PropOrder {
			for _, p := range ini.Properties[k] {
				b.WriteString(k)
				b.WriteString("=")
				b.WriteString(p.Value)
				b.WriteString("\n")
			}
		}
		if len(ini.Sections) > 0 {
			b.WriteString("\n")
		}
	}

	// Sections
	first := true
	for _, name := range ini.Order {
		sec := ini.Sections[name]
		if sec == nil {
			continue
		}
		sec.mu.RLock()
		if len(sec.Order) == 0 {
			sec.mu.RUnlock()
			continue
		}
		if !first {
			b.WriteString("\n")
		}
		first = false
		b.WriteString("[")
		b.WriteString(name)
		b.WriteString("]\n")
		if len(sec.Order) > 0 {
			for _, key := range sec.Order {
				for _, f := range sec.Fields[key] {
					b.WriteString(key)
					b.WriteString("=")
					b.WriteString(f.Value)
					b.WriteString("\n")
				}
			}
		}
		sec.mu.RUnlock()
	}

	return b.String()
}

// NormaliseNumeric prepares a number string for parsing: remove underscores, handle localised decimal comma.
func NormaliseNumeric(raw string) string {
	return str.NormaliseNumeric(raw)
}

// Parse parses an INI document string into the existing INI structure, honouring any
// ExpectedTypes declared on the receiver. This is handy when you want to declare
// types before parsing (e.g. hex integers or locale-specific decimals).
func (ini *INI) Parse(s string) error {
	if s == "" {
		return errors.New("empty input")
	}

	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return errors.New("no lines")
	}

	ini.mu.Lock()
	defer ini.mu.Unlock()

	for len(lines) > 0 {
		line := strings.TrimSpace(lines[0])
		lines = lines[1:]
		if line == "" {
			continue
		}
		if line[0] == ';' || line[0] == '#' {
			continue
		}
		if line[0] == '[' {
			// Section header
			if len(line) < 2 || line[len(line)-1] != ']' {
				continue
			}
			name := line[1 : len(line)-1]
			name = strings.ToLower(name)
			sec := ini.addSection(name)
			// Consume following lines until next section or EOF
			for len(lines) > 0 {
				ln := strings.TrimSpace(lines[0])
				if ln == "" {
					lines = lines[1:]
					continue
				}
				if ln[0] == ';' || ln[0] == '#' {
					lines = lines[1:]
					continue
				}
				if ln[0] == '[' {
					break
				}
				// key=value
				a := splitProp(ln)
				key := strings.ToLower(a[0])
				val := strings.TrimSpace(a[1])
				if key == "" {
					lines = lines[1:]
					continue
				}
				// Check expected type declarations first
				if ini.ExpectedTypes != nil {
					if m, ok := ini.ExpectedTypes[name]; ok {
						if t, ok2 := m[key]; ok2 {
							switch t {
							case Bool:
								_ = sec.AddBool(key, boolValue(val))
								lines = lines[1:]
								continue
							case Int:
								n := NormaliseNumeric(val)
								if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
									if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
										_ = sec.AddInt(key, iv)
										lines = lines[1:]
										continue
									}
								} else {
									if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
										_ = sec.AddInt(key, iv)
										lines = lines[1:]
										continue
									}
								}
							case Float:
								n := NormaliseNumeric(val)
								if fv, err := strconv.ParseFloat(n, 64); err == nil {
									_ = sec.AddFloat(key, fv)
									lines = lines[1:]
									continue
								}
							case String:
								_ = sec.AddString(key, val)
								lines = lines[1:]
								continue
							}
						}
					}
				}
				// try bool
				if _, ok := str.BoolFromString(val); ok {
					_ = sec.AddBool(key, boolValue(val))
					lines = lines[1:]
					continue
				}
				// try int (support hex and '_' separators)
				n := NormaliseNumeric(val)
				if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
					if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
						_ = sec.AddInt(key, iv)
						lines = lines[1:]
						continue
					}
				} else {
					if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
						_ = sec.AddInt(key, iv)
						lines = lines[1:]
						continue
					}
				}
				// try float
				if fv, err := strconv.ParseFloat(n, 64); err == nil {
					_ = sec.AddFloat(key, fv)
					lines = lines[1:]
					continue
				}
				_ = sec.AddString(key, val)
				lines = lines[1:]
			}
		} else {
			a := splitProp(line)
			key := strings.ToLower(a[0])
			val := strings.TrimSpace(a[1])
			if key == "" {
				continue
			}
			// Check expected type declarations for top-level properties
			if ini.ExpectedTypes != nil {
				if m, ok := ini.ExpectedTypes[""]; ok {
					if t, ok2 := m[key]; ok2 {
						switch t {
						case Bool:
							ini.addProp(key, val)
							fields := ini.Properties[key]
							fields[len(fields)-1].SetBool(boolValue(val))
							continue
						case Int:
							n := NormaliseNumeric(val)
							if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
								if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
									ini.addProp(key, val)
									fields := ini.Properties[key]
									fields[len(fields)-1].SetInt(iv)
									continue
								}
							} else {
								if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
									ini.addProp(key, val)
									fields := ini.Properties[key]
									fields[len(fields)-1].SetInt(iv)
									continue
								}
							}
						case Float:
							n := NormaliseNumeric(val)
							if fv, err := strconv.ParseFloat(n, 64); err == nil {
								ini.addProp(key, val)
								fields := ini.Properties[key]
								fields[len(fields)-1].SetFloat(fv)
								continue
							}
						case String:
							ini.addProp(key, val)
							continue
						}
					}
				}
			}
			// try bool
			if _, ok := str.BoolFromString(val); ok {
				ini.addProp(key, val)
				fields := ini.Properties[key]
				fields[len(fields)-1].SetBool(boolValue(val))
				continue
			}
			// try int
			n := NormaliseNumeric(val)
			if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
				if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
					ini.addProp(key, val)
					fields := ini.Properties[key]
					fields[len(fields)-1].SetInt(iv)
					continue
				}
			} else {
				if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
					ini.addProp(key, val)
					fields := ini.Properties[key]
					fields[len(fields)-1].SetInt(iv)
					continue
				}
			}
			// try float
			if fv, err := strconv.ParseFloat(n, 64); err == nil {
				ini.addProp(key, val)
				fields := ini.Properties[key]
				fields[len(fields)-1].SetFloat(fv)
				continue
			}
			ini.addProp(key, val)
		}
	}

	return nil
}

// Unmarshal parses an INI document string into a new INI structure.
func Unmarshal(s string) (*INI, error) {
	ini, _ := New()
	if err := ini.Parse(s); err != nil {
		return nil, err
	}
	return ini, nil
}

// addProp adds a top-level property without locking. Caller must hold ini.mu.
func (ini *INI) addProp(key, val string) error {
	k := strings.ToLower(key)
	if k == "" {
		return ErrEmptyKey
	}
	f := Field{}
	f.SetString(val)
	ini.Properties[k] = append(ini.Properties[k], &f)
	if len(ini.Properties[k]) == 1 {
		ini.PropOrder = append(ini.PropOrder, k)
	}
	return nil
}

func splitProp(s string) []string {
	a := strings.SplitN(s, "=", 2)
	if len(a) == 1 {
		return []string{strings.TrimSpace(a[0]), ""}
	}
	a[0] = strings.TrimSpace(a[0])
	a[1] = strings.TrimSpace(a[1])
	return a
}
