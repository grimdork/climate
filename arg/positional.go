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
		if o.Default != nil {
			return o.Default.([]string)
		}

		return nil
	}

	return o.Value.([]string)
}
