package arg

import "fmt"

// Option definition.
type Option struct {
	// Placeholder is the placeholder variable shown in help text.
	Placeholder string
	// ShortName of the option.
	ShortName string
	// LongName of the option.
	LongName string
	// Help text of the option.
	Help string

	// Value of the option.
	Value any
	// Default value if unspecified.
	Default any
	// Choices allowed for the option.
	Choices []any

	// Type of value.
	Type uint8
	// Required is true if this must be defined. A default would satisfy this.
	Required bool
}

// Variable types
const (
	// 	VarBool option.
	VarBool uint8 = iota
	// 	VarInt option.
	VarInt
	// VarIntSlice option.
	VarIntSlice
	// 	VarFloat option.
	VarFloat
	// 	VarFloatSlice option.
	VarFloatSlice
	// 	VarString option.
	VarString
	// VarStringSlice option.
	VarStringSlice
)

// ValidDefault returns true if the default value is usable (not the zero-value).
func (o *Option) ValidDefault() bool {
	switch o.Value.(type) {
	case bool:
		if o.Type != VarBool {
			return false
		}

	case int:
		if o.Type != VarInt && o.Default.(int) != 0 {
			return false
		}

	case []int:
		if o.Type != VarIntSlice && len(o.Default.([]int)) > 0 {
			return false
		}

	case float64:
		if o.Type != VarFloat && o.Default.(float64) != 0.0 {
			return false
		}

	case []float64:
		if o.Type != VarFloatSlice && len(o.Default.([]float64)) > 0 {
			return false
		}

	case string:
		if o.Type != VarString && o.Default.(string) != "" {
			return false
		}

	case []string:
		if o.Type != VarStringSlice && len(o.Default.([]string)) > 0 {
			return false
		}
	}

	return true
}

// SetOption sets an option.
func (opt *Options) SetOption(group, short, long, help string, defaultvalue any, required bool, t uint8, choices []any) error {
	if len(short) > 1 {
		return fmt.Errorf("-%s: %w", short, ErrLongShort)
	}

	if len(long) == 1 {
		return fmt.Errorf("--%s: %w", long, ErrShortLong)
	}

	g := opt.GetGroup(group)
	if g == nil {
		g = opt.AddGroup(group)
	}

	o := &Option{
		ShortName: short,
		LongName:  long,
		Help:      help,
		Default:   defaultvalue,
		Choices:   choices,
		Type:      t,
		Required:  required,
	}

	g.options = append(g.options, o)
	if short != "" {
		opt.short[short] = o
	}

	if long != "" {
		opt.long[long] = o
	}

	return nil
}
