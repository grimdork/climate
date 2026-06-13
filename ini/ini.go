package ini

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/grimdork/climate/env"
	"github.com/grimdork/climate/str"
)

// INI file base structure.
type INI struct {
	mu sync.RWMutex
	// Sections with settings.
	Sections map[string]*Section
	// Order sections were loaded or added in.
	Order []string
	// Properties outside of sections (top-level)
	Properties map[string][]*Field
	// PropOrder preserves insertion order of top-level properties.
	PropOrder []string
	// ExpectedTypes lets callers declare types for specific section/key pairs.
	// Map: section -> key -> type (Bool/Int/Float/String). Use section == "" for top-level properties.
	ExpectedTypes map[string]map[string]byte
	// When true, environment variable lookup prefers upper-case keys
	upper bool
	// If true, all Get* helpers will consult the environment variable first.
	// The environment key's case is controlled by the `upper` flag.
	envFirst bool
	// FilePerm controls permissions used when saving (default 0644). Set via SetSecure.
	FilePerm os.FileMode
}

const (
	// Bool type
	Bool = iota
	// Int type
	Int
	// Float type
	Float
	// String type
	String
)

// New returns an empty INI structure.
// Returns (*INI, error) to allow future initialization errors. Currently always succeeds.
func New() (*INI, error) {
	return &INI{
		Sections:      make(map[string]*Section),
		Properties:    make(map[string][]*Field),
		ExpectedTypes: make(map[string]map[string]byte),
		FilePerm:      0644,
	}, nil
}

// SetSecure sets file permissions used by Save(). If secure is true, Save uses 0600.
func (ini *INI) SetSecure(secure bool) {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	if secure {
		ini.FilePerm = 0600
	} else {
		ini.FilePerm = 0644
	}
}

// ForceUpper checks the environment for upper-case versions of the supplied INI variables.
func (ini *INI) ForceUpper() {
	ini.mu.Lock()
	ini.upper = true
	ini.mu.Unlock()
}

// SetEnvFirst enables or disables consulting environment variables before falling back
// to values in the INI. When enabled, all Get* helpers will check the environment first.
func (ini *INI) SetEnvFirst(on bool) {
	ini.mu.Lock()
	ini.envFirst = on
	ini.mu.Unlock()
}

// Set sets a property value. If the key already exists, the first entry is
// replaced; otherwise a new entry is added. Returns an error if the key is empty.
func (ini *INI) Set(s, k, v string) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		if len(ini.Properties[k]) > 0 {
			ini.Properties[k][0].SetString(v)
			return nil
		}
		f := Field{}
		f.SetString(v)
		ini.Properties[k] = []*Field{&f}
		ini.PropOrder = append(ini.PropOrder, k)
		return nil
	}

	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.SetString(k, v)
}

// Add appends a new property value. The key may appear multiple times in the
// section after repeated calls. Returns an error if the key is empty.
func (ini *INI) Add(s, k, v string) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		f := Field{}
		f.SetString(v)
		ini.Properties[k] = append(ini.Properties[k], &f)
		if len(ini.Properties[k]) == 1 {
			ini.PropOrder = append(ini.PropOrder, k)
		}
		return nil
	}

	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.AddString(k, v)
}

// SetBool sets a boolean property. Replaces the first entry or adds a new one.
func (ini *INI) SetBool(s, k string, v bool) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		if len(ini.Properties[k]) > 0 {
			ini.Properties[k][0].SetBool(v)
			return nil
		}
		f := Field{}
		f.SetBool(v)
		ini.Properties[k] = []*Field{&f}
		ini.PropOrder = append(ini.PropOrder, k)
		return nil
	}
	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.SetBool(k, v)
}

// AddBool appends a boolean property.
func (ini *INI) AddBool(s, k string, v bool) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		f := Field{}
		f.SetBool(v)
		ini.Properties[k] = append(ini.Properties[k], &f)
		if len(ini.Properties[k]) == 1 {
			ini.PropOrder = append(ini.PropOrder, k)
		}
		return nil
	}
	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.AddBool(k, v)
}

// SetInt sets an integer property. Replaces the first entry or adds a new one.
func (ini *INI) SetInt(s, k string, v int64) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		if len(ini.Properties[k]) > 0 {
			ini.Properties[k][0].SetInt(v)
			return nil
		}
		f := Field{}
		f.SetInt(v)
		ini.Properties[k] = []*Field{&f}
		ini.PropOrder = append(ini.PropOrder, k)
		return nil
	}
	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.SetInt(k, v)
}

// AddInt appends an integer property.
func (ini *INI) AddInt(s, k string, v int64) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		f := Field{}
		f.SetInt(v)
		ini.Properties[k] = append(ini.Properties[k], &f)
		if len(ini.Properties[k]) == 1 {
			ini.PropOrder = append(ini.PropOrder, k)
		}
		return nil
	}
	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.AddInt(k, v)
}

// SetFloat sets a float property. Replaces the first entry or adds a new one.
func (ini *INI) SetFloat(s, k string, v float64) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		if len(ini.Properties[k]) > 0 {
			ini.Properties[k][0].SetFloat(v)
			return nil
		}
		f := Field{}
		f.SetFloat(v)
		ini.Properties[k] = []*Field{&f}
		ini.PropOrder = append(ini.PropOrder, k)
		return nil
	}
	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.SetFloat(k, v)
}

// AddFloat appends a float property.
func (ini *INI) AddFloat(s, k string, v float64) error {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	k = strings.ToLower(k)
	if k == "" {
		return errors.New("empty key")
	}
	if s == "" {
		f := Field{}
		f.SetFloat(v)
		ini.Properties[k] = append(ini.Properties[k], &f)
		if len(ini.Properties[k]) == 1 {
			ini.PropOrder = append(ini.PropOrder, k)
		}
		return nil
	}
	sec, ok := ini.Sections[s]
	if !ok {
		sec = ini.AddSection(s)
	}
	return sec.AddFloat(k, v)
}

// DeclareType declares the expected type for a given section/key. Use section=="" for top-level properties.
func (ini *INI) DeclareType(section, key string, t byte) {
	ini.mu.Lock()
	defer ini.mu.Unlock()
	if ini.ExpectedTypes == nil {
		ini.ExpectedTypes = make(map[string]map[string]byte)
	}
	m, ok := ini.ExpectedTypes[section]
	if !ok {
		m = make(map[string]byte)
		ini.ExpectedTypes[section] = m
	}
	m[strings.ToLower(key)] = t
}

// GetString returns a variable from specified section.
// If envFirst is enabled on the INI, the environment is consulted first. The
// `upper` flag controls whether env keys are upper-cased before lookup.
func (ini *INI) GetString(s, k string) string {
	ini.mu.RLock()
	envFirst := ini.envFirst
	upper := ini.upper
	ini.mu.RUnlock()

	l := strings.ToLower(k)
	var fallback string
	if s == "" {
		ini.mu.RLock()
		p, ok := ini.Properties[l]
		ini.mu.RUnlock()
		if ok && len(p) > 0 {
			fallback = p[0].Value
		}
	} else {
		ini.mu.RLock()
		sec, ok := ini.Sections[s]
		ini.mu.RUnlock()
		if ok {
			fallback = sec.GetString(l, "")
		}
	}

	if envFirst {
		kk := k
		if upper {
			kk = strings.ToUpper(kk)
		}
		return env.Get(kk, fallback)
	}

	return fallback
}

// GetBool returns the boolean value for a variable in the INI file.
// If envFirst is enabled, the environment is consulted first.
func (ini *INI) GetBool(s, k string) bool {
	ini.mu.RLock()
	envFirst := ini.envFirst
	upper := ini.upper
	ini.mu.RUnlock()

	var fallback bool
	if s == "" {
		ini.mu.RLock()
		p, ok := ini.Properties[k]
		ini.mu.RUnlock()
		if ok && len(p) > 0 {
			fallback = p[0].GetBool()
		}
	} else {
		fallback = ini.Sections[s].GetBool(k, false)
	}

	if envFirst {
		kk := k
		if upper {
			kk = strings.ToUpper(kk)
		}
		return env.GetBool(kk, fallback)
	}

	return fallback
}

// GetInt returns the integer value for a variable in the INI file.
// If envFirst is enabled, the environment is consulted first.
func (ini *INI) GetInt(s, k string) int64 {
	ini.mu.RLock()
	envFirst := ini.envFirst
	upper := ini.upper
	ini.mu.RUnlock()

	var fallback int64
	if s == "" {
		ini.mu.RLock()
		p, ok := ini.Properties[k]
		ini.mu.RUnlock()
		if ok && len(p) > 0 {
			fallback = p[0].GetInt()
		}
	} else {
		fallback = ini.Sections[s].GetInt(k, 0)
	}

	if envFirst {
		kk := k
		if upper {
			kk = strings.ToUpper(kk)
		}
		return env.GetInt(kk, fallback)
	}

	return fallback
}

// GetFloat returns the float value for a variable in the INI file.
// If envFirst is enabled, the environment is consulted first.
func (ini *INI) GetFloat(s, k string) float64 {
	ini.mu.RLock()
	envFirst := ini.envFirst
	upper := ini.upper
	ini.mu.RUnlock()

	var fallback float64
	if s == "" {
		ini.mu.RLock()
		p, ok := ini.Properties[k]
		ini.mu.RUnlock()
		if ok && len(p) > 0 {
			fallback = p[0].GetFloat()
		}
	} else {
		fallback = ini.Sections[s].GetFloat(k, 0.0)
	}

	if envFirst {
		kk := k
		if upper {
			kk = strings.ToUpper(kk)
		}
		return env.GetFloat(kk, fallback)
	}

	return fallback
}

// GetMatch returns all fields for a key. For section-scoped keys it returns
// the section's fields; for top-level (section == "") it returns property fields.
func (ini *INI) GetMatch(s, k string) []*Field {
	ini.mu.RLock()
	defer ini.mu.RUnlock()
	k = strings.ToLower(k)
	if s == "" {
		v, ok := ini.Properties[k]
		if !ok {
			return nil
		}
		return v
	}
	sec, ok := ini.Sections[s]
	if !ok {
		return nil
	}
	return sec.GetAll(k)
}

// Load INI from file and take a guess at the types of each value.
func Load(filename string) (*INI, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ini, _ := New()
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil && l == "" {
			if err == io.EOF {
				break
			}
			return ini, err
		}

		// Trim trailing newline characters safely.
		l = strings.TrimRight(l, "\r\n")
		if l == "" {
			if err == io.EOF {
				break
			}
			// continue reading
			if err != nil {
				// non-EOF error but we have content handled above
				return ini, err
			}
			continue
		}

		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			name := l[1 : len(l)-1]
			s := ini.AddSection(name)
			s.parse(r)
			continue
		}

		// Handle top-level properties (key=value) if present.
		a := splitProp(l)
		// key is non-empty for valid input so the error is unused
		_ = ini.Add("", a[0], a[1])

		if err == io.EOF {
			break
		}
	}

	return ini, nil
}

// Save outputs the INI to a file.
// If tabbed is true, the fields will be saved with a tab character prepended.
func (ini *INI) Save(filename string, tabbed bool) error {
	ini.mu.RLock()
	filePerm := ini.FilePerm
	defer ini.mu.RUnlock()

	b := str.NewStringer()

	// Top-level properties
	for _, key := range ini.PropOrder {
		for _, f := range ini.Properties[key] {
			if tabbed {
				b.WriteRune('\t')
			}
			b.WriteStrings(key, "=", f.Value, "\n")
		}
	}

	// Sections
	first := true
	for _, secname := range ini.Order {
		sec := ini.Sections[secname]
		sec.mu.RLock()
		if len(sec.Order) == 0 {
			sec.mu.RUnlock()
			continue
		}
		if !first {
			b.WriteString("\n")
		}
		first = false
		b.WriteStrings("[", secname, "]\n")
		for _, key := range sec.Order {
			for _, f := range sec.Fields[key] {
				if tabbed {
					b.WriteRune('\t')
				}
				b.WriteStrings(key, "=", f.Value, "\n")
			}
		}
		sec.mu.RUnlock()
	}
	return os.WriteFile(filename, []byte(b.String()), filePerm)
}

// AddSection adds a named section to the INI or returns an existing one.
// It is NOT safe for concurrent use unless the caller holds ini.mu.
func (ini *INI) AddSection(name string) *Section {
	if sec := ini.Sections[name]; sec != nil {
		return sec
	}
	sec := &Section{
		Fields: make(map[string][]*Field),
	}
	ini.Sections[name] = sec
	ini.Order = append(ini.Order, name)
	return sec
}
