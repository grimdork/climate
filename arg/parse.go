package arg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ShowOptions shows the values of all options. Used for debugging.
func (opt *Options) ShowOptions() {
	for _, o := range opt.short {
		fmt.Printf("-%s: %v (%v)\n", o.ShortName, o.Value, o.Default)
	}

	for _, o := range opt.long {
		fmt.Printf("--%s: %v (%v)\n", o.LongName, o.Value, o.Default)
	}
}

// Parse command line arguments from a string slice.
// - If default help is defined, it will print the help message after parsing when "-h" or "--help" is supplied,
// then os.Exit(0). Returns ErrNoArgs if no arguments are supplied.
//
// Tool commands, short options (single dash and one letter), long options (double dash and one or more
// letters), and positional arguments are each paarsed in the order they are supplied. If a positional
// argument is of a slice type, it will swallow all remaining arguments, including long and short options.
//
// Single- and double-dash options found before any tool commands are parsed for the Options structure.
//
// Tool commands break the parsing off, and calls the command with the remaining arguments after running
//  any handlers for the pre-command options.
//
// Options criteria:
// - Short options start with a single dash ("-").
// - Short boolean options don't need to take a value.
// - Short boolean options require an equal sign ("=") after the option with a truthy or falsy value.
// - Truthy values are "true", "yes", "on", "1", and "t".
// - Falsy values are everything else.
// - Short options can be combined ("-a -b" can be written as "-ab").
// - Combined short options allow only the last one to take a value. The ones before must be booleans.
//
// - Long options start with a double dash ("--").
// - Long options are followed by either whitespace or an equal sign ("--foo bar" or "--foo=bar").
func (opt *Options) Parse(args []string) error {
	if len(os.Args) == 1 {
		return ErrNoArgs
	}

	err := opt.parseArgs(args)
	if err != nil {
		return err
	}

	if opt.hashelp && opt.GetBool("h") {
		opt.PrintHelp()
		os.Exit(0)
	}

	return nil
}

func (opt *Options) parseArgs(args []string) error {
	unknown := []string{}
	pos := opt.positional
	for i, arg := range args {
		if arg == "" {
			continue
		}

		cmd := opt.commands[arg]
		if cmd != nil {
			fn := cmd.Func
			if fn == nil {
				return fmt.Errorf("%s: %s", arg, ErrMissingFunc)
			}

			fn(cmd.Options, args[i+1:])
			return nil
		}

		if len(arg) < 2 && len(pos) < 0 {
			unknown = append(unknown, arg)
			continue
		}

		//
		// Long options
		//

		if arg[0] == '-' && arg[1] == '-' {
			arg = arg[2:]
			if arg == "" {
				return ErrEmptyLong
			}

			a := splitOption(arg)
			o, ok := opt.long[a[0]]
			if ok {
				switch o.Type {
				case VarBool:
					t, v := isTruthy(a[1])
					// We have the form "--option=value"
					if t {
						o.Value = v
						continue
					}

					if len(args) > i+1 {
						t, v = isTruthy(args[i+1])
						// We have the form "--option value"
						if t {
							o.Value = v
							args[i+1] = ""
							continue
						}
					}

					// It's a standalone boolean option, so just set it to teue. Phew!
					o.Value = true

				case VarString:
					if a[1] != "" {
						o.Value = a[1]
						continue
					}

					if len(args) > i+1 {
						o.Value = args[i+1]
						args[i+1] = ""
						continue
					}

					return fmt.Errorf("--%s: %w", o.LongName, ErrMissingParam)

				case VarInt:
					if a[1] != "" {
						v, err := strconv.Atoi(a[1])
						if err != nil {
							return err
						}

						o.Value = v
						continue
					}

					if len(args) > i+1 {
						v, err := strconv.Atoi(args[i+1])
						if err != nil {
							return err
						}

						o.Value = v
						args[i+1] = ""
						continue
					}

					return fmt.Errorf("--%s: %w", o.LongName, ErrMissingParam)

				case VarFloat:
					if a[1] != "" {
						v, err := strconv.ParseFloat(a[1], 64)
						if err != nil {
							return err
						}

						o.Value = v
						continue
					}

					if len(args) > i+1 {
						v, err := strconv.ParseFloat(args[i+1], 64)
						if err != nil {
							return err
						}

						o.Value = v
						args[i+1] = ""
						continue
					}

					return fmt.Errorf("--%s: %w", o.LongName, ErrMissingParam)

				default:
					return fmt.Errorf("--%s: %w", o.LongName, ErrUnknownType)
				} // switch o.Type
			} else {
				return fmt.Errorf("--%s: %w", a[0], ErrUnknownOption)
			} // if long option is defined
			continue
		} // if long option

		//
		// Short options
		//

		if arg[0] == '-' {
			s := arg[1:]
			a := splitOption(s)
			if a[1] != "" {
				s = a[0]
			}

			for _, c := range s {
				o, ok := opt.short[string(c)]
				if ok {
					switch o.Type {
					case VarBool:
						if a[0] == string(c) && a[1] != "" {
							_, v := isTruthy(a[1])
							o.Value = v
							continue
						}

						if len(args) > i+1 {
							t, v := isTruthy(args[i+1])
							if t {
								o.Value = v
								args[i+1] = ""
								continue
							}
						}

						o.Value = true

					case VarString:
						if a[0] == string(c) && a[1] != "" {
							o.Value = a[1]
							continue
						}

						if len(args) > i+1 {
							o.Value = args[i+1]
							args[i+1] = ""
							continue
						}

						return fmt.Errorf("-%c: %w", c, ErrMissingParam)

					case VarInt:
						if a[0] == string(c) && a[1] != "" {
							v, err := strconv.Atoi(a[1])
							if err != nil {
								return err
							}

							o.Value = v
							continue
						}

						if len(args) > i+1 {
							v, err := strconv.Atoi(args[i+1])
							if err != nil {
								return err
							}

							o.Value = v
							args[i+1] = ""
							continue
						}

						return fmt.Errorf("-%c: %w", c, ErrMissingParam)

					case VarFloat:
						if a[0] == string(c) && a[1] != "" {
							v, err := strconv.ParseFloat(a[1], 64)
							if err != nil {
								return err
							}

							o.Value = v
							continue
						}

						if len(args) > i+1 {
							v, err := strconv.ParseFloat(args[i+1], 64)
							if err != nil {
								return err
							}

							o.Value = v
							args[i+1] = ""
							continue
						}

						return fmt.Errorf("-%c: %w", c, ErrMissingParam)

					} // switch o.Type
				} else {
					return fmt.Errorf("-%c: %w", c, ErrUnknownOption)
				} // if short option is defined
			} // range s
			continue
		} // if short option

		//
		// Positional arguments
		//

		if len(pos) > 0 {
			switch pos[0].Type {
			case VarBool:
				t, v := isTruthy(arg)
				if t {
					pos[0].Value = v
				} else {
					pos[0].Value = false
				}

			case VarString:
				pos[0].Value = arg

			case VarStringSlice:
				if pos[0].Value == nil {
					pos[0].Value = []string{}
				}

				pos[0].Value = append(pos[0].Value.([]string), arg)
				continue

			case VarInt:
				v, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}

				pos[0].Value = v

			case VarIntSlice:
				if pos[0].Value == nil {
					pos[0].Value = []int{}
				}

				v, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}

				pos[0].Value = append(pos[0].Value.([]int), v)
				continue

			case VarFloat:
				v, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return err
				}

				pos[0].Value = v

			case VarFloatSlice:
				if pos[0].Value == nil {
					pos[0].Value = []float64{}
				}

				v, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return err
				}

				pos[0].Value = append(pos[0].Value.([]float64), v)
				continue
			}

			pos = pos[1:]
			continue
		}

		// The leftovers go here
		unknown = append(unknown, arg)
	}

	for _, arg := range unknown {
		if arg != "" {
		}
		unknown = unknown[1:]
	}
	opt.Remainder = unknown

	for _, o := range opt.short {
		if o.Required && o.Value == nil {
			return fmt.Errorf("-%s: %w", o.ShortName, ErrMissingRequired)
		}
	}

	for _, o := range opt.long {
		if o.Required && o.Value == nil {
			return fmt.Errorf("--%s: %w", o.LongName, ErrMissingRequired)
		}
	}

	return nil
}

func splitOption(arg string) []string {
	a := strings.SplitN(arg, "=", 2)
	if len(a) == 1 {
		return []string{arg, ""}
	}

	return a
}

// isTruthy returns whether the supplied string is a truthy value.
// The second value is the decoded value, if applicable, false otherwise.
func isTruthy(s string) (bool, bool) {
	switch strings.ToLower(s) {
	case "true", "yes", "on", "1", "t":
		return true, true
	case "false", "no", "off", "0", "f":
		return true, false
	}

	return false, false
}
