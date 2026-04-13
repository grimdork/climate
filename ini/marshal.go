package ini

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// Marshal INI into a string. Global properties first (sorted), then sections (sorted).
func (ini *INI) Marshal() string {
	ini.mu.RLock()
	defer ini.mu.RUnlock()

	b := &strings.Builder{}

	// Properties
	if len(ini.Properties) > 0 {
		keys := make([]string, 0, len(ini.Properties))
		for k := range ini.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			p := ini.Properties[k]
			b.WriteString(k)
			b.WriteString(" = ")
			b.WriteString(p.Value)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Sections
	if len(ini.Sections) > 0 {
		secNames := make([]string, 0, len(ini.Sections))
		for name := range ini.Sections {
			secNames = append(secNames, name)
		}
		sort.Strings(secNames)
		for _, name := range secNames {
			b.WriteString("[")
			b.WriteString(name)
			b.WriteString("]\n")
			sec := ini.Sections[name]
			// Preserve field order if available
			if len(sec.Order) > 0 {
				for _, key := range sec.Order {
					f := sec.Fields[key]
					b.WriteString(key)
					b.WriteString(" = ")
					b.WriteString(f.Value)
					b.WriteString("\n")
				}
			} else {
				// Fallback: sort keys
				keys := make([]string, 0, len(sec.Fields))
				for k := range sec.Fields {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					f := sec.Fields[k]
					b.WriteString(k)
					b.WriteString(" = ")
					b.WriteString(f.Value)
					b.WriteString("\n")
				}
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// normaliseNumeric prepares a number string for parsing: remove underscores, handle localised decimal comma.
func normaliseNumeric(raw string) string {
	r := strings.TrimSpace(raw)
	// remove common thousands separators (underscores)
	r = strings.ReplaceAll(r, "_", "")
	// If it contains a comma but no dot, assume comma is decimal separator and convert to dot.
	if strings.Contains(r, ",") && !strings.Contains(r, ".") {
		r = strings.ReplaceAll(r, ",", ".")
	}
	return r
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
			sec := ini.AddSection(name)
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
				// Check expected type declarations first
				if ini.ExpectedTypes != nil {
					if m, ok := ini.ExpectedTypes[name]; ok {
						if t, ok2 := m[key]; ok2 {
							switch t {
							case Bool:
								sec.AddBool(key, boolValue(val))
								lines = lines[1:]
								continue
							case Int:
								n := normaliseNumeric(val)
								if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
									if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
										sec.AddInt(key, iv)
										lines = lines[1:]
										continue
									}
								} else {
									if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
										sec.AddInt(key, iv)
										lines = lines[1:]
										continue
									}
								}
							case Float:
								n := normaliseNumeric(val)
								if fv, err := strconv.ParseFloat(n, 64); err == nil {
									sec.AddFloat(key, fv)
									lines = lines[1:]
									continue
								}
							case String:
								sec.AddString(key, val)
								lines = lines[1:]
								continue
							}
						}
					}
				}
				// try bool
				switch strings.ToLower(val) {
				case "yes", "true", "on", "no", "false", "off":
					sec.AddBool(key, boolValue(val))
					lines = lines[1:]
					continue
				}
				// try int (support hex and '_' separators)
				n := normaliseNumeric(val)
				if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
					if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
						sec.AddInt(key, iv)
						lines = lines[1:]
						continue
					}
				} else {
					if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
						sec.AddInt(key, iv)
						lines = lines[1:]
						continue
					}
				}
				// try float
				if fv, err := strconv.ParseFloat(n, 64); err == nil {
					sec.AddFloat(key, fv)
					lines = lines[1:]
					continue
				}
				sec.AddString(key, val)
				lines = lines[1:]
			}
		} else {
			a := splitProp(line)
			key := strings.ToLower(a[0])
			val := strings.TrimSpace(a[1])
			// Check expected type declarations for top-level properties
			if ini.ExpectedTypes != nil {
				if m, ok := ini.ExpectedTypes[""]; ok {
					if t, ok2 := m[key]; ok2 {
						switch t {
						case Bool:
							ini.Set("", key, val)
							p := ini.Properties[key]
							p.SetBool(boolValue(val))
							continue
						case Int:
							n := normaliseNumeric(val)
							if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
								if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
									ini.Set("", key, val)
									p := ini.Properties[key]
									p.SetInt(iv)
									continue
								}
							} else {
								if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
									ini.Set("", key, val)
									p := ini.Properties[key]
									p.SetInt(iv)
									continue
								}
							}
						case Float:
							n := normaliseNumeric(val)
							if fv, err := strconv.ParseFloat(n, 64); err == nil {
								ini.Set("", key, val)
								p := ini.Properties[key]
								p.SetFloat(fv)
								continue
							}
						case String:
							ini.Set("", key, val)
							continue
						}
					}
				}
			}
			// try bool
			switch strings.ToLower(val) {
			case "yes", "true", "on", "no", "false", "off":
				ini.Set("", key, val)
				continue
			}
			// try int
			n := normaliseNumeric(val)
			if strings.HasPrefix(n, "0x") || strings.HasPrefix(n, "0X") {
				if iv, err := strconv.ParseInt(n, 0, 64); err == nil {
					ini.Set("", key, val)
					// store typed value
					p := ini.Properties[key]
					p.SetInt(iv)
					continue
				}
			} else {
				if iv, err := strconv.ParseInt(n, 10, 64); err == nil {
					ini.Set("", key, val)
					// store typed value
					p := ini.Properties[key]
					p.SetInt(iv)
					continue
				}
			}
			// try float
			if fv, err := strconv.ParseFloat(n, 64); err == nil {
				ini.Set("", key, val)
				p := ini.Properties[key]
				p.SetFloat(fv)
				continue
			}
			ini.Set("", key, val)
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

func splitProp(s string) []string {
	a := strings.SplitN(s, "=", 2)
	if len(a) == 1 {
		return []string{strings.TrimSpace(a[0]), ""}
	}
	a[0] = strings.TrimSpace(a[0])
	a[1] = strings.TrimSpace(a[1])
	return a
}
