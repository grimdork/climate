package arg

import (
	"fmt"
	"os"
)

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
	switch o.Default.(type) {
	case bool:
		if o.Type != VarBool {
			return false
		}

		return true

	case int:
		if o.Type != VarInt || o.Default.(int) == 0 {
			return false
		}

		return true

	case []int:
		if o.Type != VarIntSlice || len(o.Default.([]int)) == 0 {
			return false
		}

		return true

	case float64:
		if o.Type != VarFloat || o.Default.(float64) == 0.0 {
			return false
		}

		return true

	case []float64:
		if o.Type != VarFloatSlice || len(o.Default.([]float64)) == 0 {
			return false
		}

		return true

	case string:
		if o.Type != VarString || o.Default.(string) == "" {
			return false
		}

		return true

	case []string:
		if o.Type != VarStringSlice || len(o.Default.([]string)) == 0 {
			return false
		}

		return true
	}

	return false
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

// HelpOrFail parses the CLI arguments, then prints the help text and exits if the -h flag is set,
// or prints an error and exits with exit code 2 if something went wrong.
func (opt *Options) HelpOrFail() {
	opt.HelpOrFailArgs(os.Args[1:])
}

// HelpOrFailArgs works like HelpOrFail(), but takes a string slice to parse. Use os.Args[1:] to mimic HelpOrFail(),
// and increase the number for sub-commands.
func (opt *Options) HelpOrFailArgs(args []string) {
	err := opt.Parse(args)
	if err != nil {
		// -h was supplied somewhere on the command line, so exit cleanly after printing help.
		if err == ErrNoArgs {
			opt.PrintHelp()
			os.Exit(0)
		}

		// Some other error occurred, probably not an issue with input arguments, so just print the error and exit.
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(2)
	}
	// If we got this far, everything is fine and the program can proceed.
}
