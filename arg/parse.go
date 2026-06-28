package arg

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/grimdork/climate/str"
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
//
// Tool commands, short options (single dash and one letter), long options (double dash and one or more
// letters), and positional arguments are each parsed in the order they are supplied. If a positional
// argument is of a slice type, it will swallow all remaining arguments, including long and short options.
//
// Single- and double-dash options found before any tool commands are parsed for the Options structure.
//
// Tool commands break parsing off, and call the command with the remaining arguments after running
//
//	any handlers for the pre-command options.
//
// When a subcommand is dispatched, CommandRun is set to the command name and Parse returns whatever
// the command function returned (nil on success).
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
	if len(args) == 0 {
		return nil
	}

	return opt.parseArgs(args)
}

// ParseAndRun is a helper method that calls Parse and exits the program on error.
// It does not return any error, but is instead intended as the final error handler in main().
func (opt *Options) ParseAndRun(args []string) {
	err := opt.Parse(args)
	if err != nil {
		if err == ErrNoArgs || len(opt.Args) == 0 {
			if opt.hashelp {
				opt.PrintHelp()
			}
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	if opt.CommandRun != "" {
		os.Exit(0)
	}

	if len(args) == 0 && opt.hashelp {
		opt.PrintHelp()
		os.Exit(0)
	}
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
				return fmt.Errorf("%s: %w", arg, ErrMissingFunc)
			}

			cmd.Options.Args = args[i+1:]
			opt.CommandRun = cmd.Name
			return fn(cmd.Options)
		}

		if len(arg) < 2 && len(pos) == 0 {
			unknown = append(unknown, arg)
			continue
		}

		//
		// Long options
		//

		if len(arg) > 1 && arg[0] == '-' && arg[1] == '-' {
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
					if t {
						o.Value = v
						continue
					}
					// Unrecognised explicit value like --option=xyz
					if a[1] != "" {
						return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalValue)
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

					// It's a standalone boolean option, so just set it to true. Phew!
					o.Value = true

				case VarString:
					if a[1] != "" {
						o.Value = a[1]
						if !hasChoice(o.Value.(string), o.Choices) {
							return fmt.Errorf("%s=%v: %w", arg, o.Value, ErrIllegalChoice)
						}
						continue
					}

					if len(args) > i+1 {
						o.Value = args[i+1]
						if !hasChoice(o.Value.(string), o.Choices) {
							return fmt.Errorf("%s=%v: %w", arg, o.Value, ErrIllegalChoice)
						}
						args[i+1] = ""
						continue
					}

					return fmt.Errorf("--%s: %w", o.LongName, ErrMissingParam)

				case VarInt:
					if a[1] != "" {
						v, err := strconv.Atoi(a[1])
						if err != nil {
							return fmt.Errorf("%s=%v: %w", arg, a[1], err)
						}
						if !hasChoice(v, o.Choices) {
							return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalChoice)
						}

						o.Value = v
						continue
					}

					if len(args) > i+1 {
						v, err := strconv.Atoi(args[i+1])
						if err != nil {
							return fmt.Errorf("--%s=%v: %w", o.LongName, args[i+1], err)
						}
						if !hasChoice(v, o.Choices) {
							return fmt.Errorf("--%s=%v: %w", o.LongName, args[i+1], ErrIllegalChoice)
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
							return fmt.Errorf("%s=%v: %w", arg, a[1], err)
						}
						if !hasChoice(v, o.Choices) {
							return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalChoice)
						}

						o.Value = v
						continue
					}

					if len(args) > i+1 {
						v, err := strconv.ParseFloat(args[i+1], 64)
						if err != nil {
							return fmt.Errorf("--%s=%v: %w", o.LongName, args[i+1], err)
						}
						if !hasChoice(v, o.Choices) {
							return fmt.Errorf("--%s=%v: %w", o.LongName, args[i+1], ErrIllegalChoice)
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
							t, v := isTruthy(a[1])
							if !t {
								return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalValue)
							}
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
							if !hasChoice(a[1], o.Choices) {
								return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalChoice)
							}
							o.Value = a[1]
							continue
						}

						if len(args) > i+1 {
							if !hasChoice(args[i+1], o.Choices) {
								return fmt.Errorf("%s=%v: %w", arg, args[i+1], ErrIllegalChoice)
							}
							o.Value = args[i+1]
							args[i+1] = ""
							continue
						}

						return fmt.Errorf("-%c: %w", c, ErrMissingParam)

					case VarInt:
						if a[0] == string(c) && a[1] != "" {
							v, err := strconv.Atoi(a[1])
							if err != nil {
								return fmt.Errorf("%s=%v: %w", arg, a[1], err)
							}
							if !hasChoice(v, o.Choices) {
								return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalChoice)
							}

							o.Value = v
							continue
						}

						if len(args) > i+1 {
							v, err := strconv.Atoi(args[i+1])
							if err != nil {
								return fmt.Errorf("-%c=%v: %w", c, args[i+1], err)
							}
							if !hasChoice(v, o.Choices) {
								return fmt.Errorf("-%c=%v: %w", c, args[i+1], ErrIllegalChoice)
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
								return fmt.Errorf("%s=%v: %w", arg, a[1], err)
							}
							if !hasChoice(v, o.Choices) {
								return fmt.Errorf("%s=%v: %w", arg, a[1], ErrIllegalChoice)
							}

							o.Value = v
							continue
						}

						if len(args) > i+1 {
							v, err := strconv.ParseFloat(args[i+1], 64)
							if err != nil {
								return fmt.Errorf("-%c=%v: %w", c, args[i+1], err)
							}
							if !hasChoice(v, o.Choices) {
								return fmt.Errorf("-%c=%v: %w", c, args[i+1], ErrIllegalChoice)
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

	opt.Args = unknown
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

	for _, o := range opt.positional {
		if o.Required && o.Value == nil {
			return fmt.Errorf("%s: %w", o.Placeholder, ErrMissingRequired)
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
	v, ok := str.BoolFromString(s)
	if !ok {
		return false, false
	}
	return true, v
}

func hasChoice[C comparable](c C, list []any) bool {
	// If the exclusion list is empty, every choice is valid
	if len(list) == 0 {
		return true
	}

	for _, x := range list {
		if x == c {
			return true
		}
	}

	return false
}
