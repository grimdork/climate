package arg

import "fmt"

// SetPositional sets a positional argument.
// Arguments which aren't long or short options or tool commands are considered positional.
func (opt *Options) SetPositional(placeholder, help string, defaultvalue any, required bool, t uint8) error {
	if len(placeholder) == 0 {
		return fmt.Errorf("%w", ErrNoPlaceholder)
	}

	o := &Option{
		Placeholder: placeholder,
		Help:        help,
		Default:     defaultvalue,
		Type:        t,
		Required:    required,
	}

	opt.positional = append(opt.positional, o)
	opt.posmap[placeholder] = o
	return nil
}

// GetPosBool returns a positional boolean's value.
func (opt *Options) GetPosBool(placeholder string) bool {
	o := opt.posmap[placeholder]
	if o == nil {
		return false
	}

	if o.Value == nil {
		if o.Default != nil {
			if v, ok := o.Default.(bool); ok {
				return v
			}
		}

		return false
	}

	if v, ok := o.Value.(bool); ok {
		return v
	}
	return false
}

// GetPosString returns a positional string's value.
func (opt *Options) GetPosString(placeholder string) string {
	o := opt.posmap[placeholder]
	if o == nil {
		return ""
	}

	if o.Value == nil {
		if o.Default != nil {
			if v, ok := o.Default.(string); ok {
				return v
			}
		}

		return ""
	}

	if v, ok := o.Value.(string); ok {
		return v
	}
	return ""
}

// GetPosStringSlice returns a positional string slice's values.
func (opt *Options) GetPosStringSlice(placeholder string) []string {
	o := opt.posmap[placeholder]
	if o == nil {
		return nil
	}

	if o.Value == nil {
		return []string{}
	}

	if v, ok := o.Value.([]string); ok {
		return v
	}
	return []string{}
}

// GetPosInt returns a positional int's value.
func (opt *Options) GetPosInt(placeholder string) int {
	o := opt.posmap[placeholder]
	if o == nil {
		return 0
	}

	if o.Value == nil {
		if o.Default != nil {
			if v, ok := o.Default.(int); ok {
				return v
			}
		}

		return 0
	}

	if v, ok := o.Value.(int); ok {
		return v
	}
	return 0
}

// GetPosIntSlice returns a positional int slice's values.
func (opt *Options) GetPosIntSlice(placeholder string) []int {
	o := opt.posmap[placeholder]
	if o == nil {
		return nil
	}

	if o.Value == nil {
		if o.Default != nil {
			if v, ok := o.Default.([]int); ok {
				return v
			}
		}

		return nil
	}

	if v, ok := o.Value.([]int); ok {
		return v
	}
	return nil
}

// GetPosFloat returns a positional float's value.
func (opt *Options) GetPosFloat(placeholder string) float64 {
	o := opt.posmap[placeholder]
	if o == nil {
		return 0
	}

	if o.Value == nil {
		if o.Default != nil {
			if v, ok := o.Default.(float64); ok {
				return v
			}
		}

		return 0
	}

	if v, ok := o.Value.(float64); ok {
		return v
	}
	return 0
}

// GetPosFloatSlice returns a positional float slice's values.
func (opt *Options) GetPosFloatSlice(placeholder string) []float64 {
	o := opt.posmap[placeholder]
	if o == nil {
		return nil
	}

	if o.Value == nil {
		if o.Default != nil {
			if v, ok := o.Default.([]float64); ok {
				return v
			}
		}

		return nil
	}

	if v, ok := o.Value.([]float64); ok {
		return v
	}
	return nil
}
