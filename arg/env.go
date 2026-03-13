package arg

import (
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
			_, t := isTruthy(vs)
			o.Value = t

		case VarInt:
			i, _ := strconv.Atoi(vs)
			o.Value = i

		case VarIntSlice:
			var list []int
			for _, v := range strings.Split(vs, delimiter) {
				i, _ := strconv.Atoi(v)
				list = append(list, i)
			}
			o.Value = list

		case VarFloat:
			f, _ := strconv.ParseFloat(vs, 64)
			o.Value = f

		case VarFloatSlice:
			var list []float64
			for _, v := range strings.Split(vs, delimiter) {
				f, _ := strconv.ParseFloat(v, 64)
				list = append(list, f)
			}
			o.Value = list

		case VarString:
			o.Value = vs

		case VarStringSlice:
			o.Value = strings.Split(vs, delimiter)
		}
	}

	return nil
}
