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
			return o.Default.(bool)
		}

		return false
	}

	return o.Value.(bool)
}

// GetPosString returns a positional string's value.
func (opt *Options) GetPosString(placeholder string) string {
	o := opt.posmap[placeholder]
	if o == nil {
		return ""
	}

	if o.Value == nil {
		if o.Default != nil {
			return o.Default.(string)
		}

		return ""
	}

	return o.Value.(string)
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

	return o.Value.([]string)
}

// GetPosInt returns a positional int's value.
func (opt *Options) GetPosInt(placeholder string) int {
	o := opt.posmap[placeholder]
	if o == nil {
		return 0
	}

	if o.Value == nil {
		if o.Default != nil {
			return o.Default.(int)
		}

		return 0
	}

	return o.Value.(int)
}

// GetPosIntSlice returns a positional int slice's values.
func (opt *Options) GetPosIntSlice(placeholder string) []int {
	o := opt.posmap[placeholder]
	if o == nil {
		return nil
	}

	if o.Value == nil {
		return o.Default.([]int)
	}

	return o.Value.([]int)
}

// GetPosFloat returns a positional float's value.
func (opt *Options) GetPosFloat(placeholder string) float64 {
	o := opt.posmap[placeholder]
	if o == nil {
		return 0
	}

	if o.Value == nil {
		if o.Default != nil {
			return o.Default.(float64)
		}

		return 0
	}

	return o.Value.(float64)
}

// GetPosFloatSlice returns a positional float slice's values.
func (opt *Options) GetPosFloatSlice(placeholder string) []float64 {
	o := opt.posmap[placeholder]
	if o == nil {
		return nil
	}

	if o.Value == nil {
		return o.Default.([]float64)
	}

	return o.Value.([]float64)
}
