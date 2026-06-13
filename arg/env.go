package arg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseEnvironment parses environment variables into options.
// By default, this parses uppercase versions of long options, but an optional prefix can be provided.
// So if you provide the prefix "moo" (case doesn't matter) and have a long option "host", this will
// be set if the envvar "MOO_HOST" is set.
// Delimiter is used when parsing slices. If unset, it defaults to a comma.
func (opt *Options) ParseEnvironment(prefix, delimiter string) error {
	for _, o := range opt.long {
		s := strings.ToUpper(o.LongName)
		if prefix != "" {
			s = strings.ToUpper(prefix) + "_" + s
		}
		vs := os.Getenv(s)
		if vs == "" {
			continue
		}

		if delimiter == "" {
			delimiter = ","
		}

		switch o.Type {
		case VarBool:
			t, v := isTruthy(vs)
			if !t {
				return fmt.Errorf("env %s=%s: %w", s, vs, ErrIllegalValue)
			}
			o.Value = v

		case VarInt:
			i, err := strconv.Atoi(vs)
			if err != nil {
				return fmt.Errorf("env %s: %w", s, err)
			}
			if !hasChoice(i, o.Choices) {
				return fmt.Errorf("env %s=%d: %w", s, i, ErrIllegalChoice)
			}
			o.Value = i

		case VarIntSlice:
			var list []int
			for _, v := range strings.Split(vs, delimiter) {
				i, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("env %s: %w", s, err)
				}
				if !hasChoice(i, o.Choices) {
					return fmt.Errorf("env %s=%d: %w", s, i, ErrIllegalChoice)
				}
				list = append(list, i)
			}
			o.Value = list

		case VarFloat:
			f, err := strconv.ParseFloat(vs, 64)
			if err != nil {
				return fmt.Errorf("env %s: %w", s, err)
			}
			if !hasChoice(f, o.Choices) {
				return fmt.Errorf("env %s=%v: %w", s, f, ErrIllegalChoice)
			}
			o.Value = f

		case VarFloatSlice:
			var list []float64
			for _, v := range strings.Split(vs, delimiter) {
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("env %s: %w", s, err)
				}
				if !hasChoice(f, o.Choices) {
					return fmt.Errorf("env %s=%v: %w", s, f, ErrIllegalChoice)
				}
				list = append(list, f)
			}
			o.Value = list

		case VarString:
			if !hasChoice(vs, o.Choices) {
				return fmt.Errorf("env %s=%s: %w", s, vs, ErrIllegalChoice)
			}
			o.Value = vs

		case VarStringSlice:
			parts := strings.Split(vs, delimiter)
			for _, p := range parts {
				if !hasChoice(p, o.Choices) {
					return fmt.Errorf("env %s=%s: %w", s, p, ErrIllegalChoice)
				}
			}
			o.Value = parts
		}
	}

	return nil
}
