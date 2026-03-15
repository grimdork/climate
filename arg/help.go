package arg

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// SetDefaultHelp sets the default help option, optionally with a short "-h" flag.
func (opt *Options) SetDefaultHelp(full bool) {
	short := ""
	if full {
		short = "h"
	}
	opt.SetOption("", short, "help", "Print this help message.", nil, false, VarBool, nil)
	opt.hashelp = true
}

// PrintHelp builds and prints the help text based on available options.
func (opt *Options) PrintHelp() {
	// SENTINEL to confirm this PrintHelp is executed in the built binary.
	fmt.Fprintln(os.Stderr, "SENTINEL_PRINTHELP")
	// Dump opt.long and opt.short map keys + pointers for diagnosis
	for k, v := range opt.long {
		fmt.Fprintf(os.Stderr, "PRINTHELP_LONG: key=%s ptr=%p choices_len=%d\n", k, v, len(v.Choices))
	}
	for k, v := range opt.short {
		fmt.Fprintf(os.Stderr, "PRINTHELP_SHORT: key=%s ptr=%p choices_len=%d\n", k, v, len(v.Choices))
	}

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 8, 8, 1, '\t', 0)
	w.Write([]byte("Usage:\n  "))
	w.Write([]byte(opt.appname))

	count := 0
	for _, g := range opt.groups {
		count += len(g.options)
	}

	if count > 0 {
		w.Write([]byte(" [OPTIONS]"))
	}

	if len(opt.short)+len(opt.long) > 1 {
		w.Write([]byte("..."))
	}

	if len(opt.commands) > 0 {
		w.Write([]byte(" [COMMAND]"))
	}

	if len(opt.positional) > 0 {
		for _, o := range opt.positional {
			fmt.Fprintf(w, " [%s]", o.Placeholder)
			if o.Type == VarStringSlice {
				fmt.Fprintf(w, "...")
			}
		}
	}
	w.Write([]byte("\n\n"))

	for _, g := range opt.groups {
		if len(g.options) > 0 {
			if g.Name == GroupDefault {
				w.Write([]byte("Main options:\n"))
			} else {
				fmt.Fprintf(w, "%s options:\n", g.Name)
			}

			for _, o := range g.options {
				if o.ShortName != "" && o.LongName != "" {
					fmt.Fprintf(w, "\t-%s, --%s\t%s", o.ShortName, o.LongName, o.Help)
				}

				if o.ShortName != "" && o.LongName == "" {
					fmt.Fprintf(w, "\t-%s\t%s", o.ShortName, o.Help)
				}

				if o.LongName != "" && o.ShortName == "" {
					fmt.Fprintf(w, "\t--%s\t%s", o.LongName, o.Help)
				}

				// choices (if provided) show after the help text
				parts := make([]string, 0, len(o.Choices))
				for _, c := range o.Choices {
					parts = append(parts, fmt.Sprintf("%v", c))
				}
				// If the group's stored option has no choices, try the global maps (short/long) as a fallback.
				if len(parts) == 0 {
					if o.LongName != "" {
						if oo := opt.GetOption(o.LongName); oo != nil && len(oo.Choices) > 0 {
							for _, c := range oo.Choices {
								parts = append(parts, fmt.Sprintf("%v", c))
							}
						}
					} else if o.ShortName != "" {
						if oo := opt.GetOption(o.ShortName); oo != nil && len(oo.Choices) > 0 {
							for _, c := range oo.Choices {
								parts = append(parts, fmt.Sprintf("%v", c))
							}
						}
					}
				}
				if len(parts) > 0 {
					fmt.Fprintf(w, " (choices: %s)", strings.Join(parts, ","))
				}
				// DEBUG: always print option diagnostics to stderr
				if len(parts) > 0 {
					fmt.Fprintf(os.Stderr, "DEBUG_PRINTHELP: group=%s opt_ptr=%p long=%s short=%s choices_len=%d sample=%v\n", g.Name, o, o.LongName, o.ShortName, len(parts), func() interface{} { if len(parts) <= 3 { return parts } return parts[0:3] }())
				} else {
					fmt.Fprintf(os.Stderr, "DEBUG_PRINTHELP: group=%s opt_ptr=%p long=%s short=%s choices_len=0\n", g.Name, o, o.LongName, o.ShortName)
				}
				// Also print pointers for opt.long/opt.short entries (if present)
				if o.LongName != "" {
					if oo := opt.GetOption(o.LongName); oo != nil {
						fmt.Fprintf(os.Stderr, "DEBUG_PRINTHELP: opt.long[%s]=%p\n", o.LongName, oo)
					}
				}
				if o.ShortName != "" {
					if oo := opt.GetOption(o.ShortName); oo != nil {
						fmt.Fprintf(os.Stderr, "DEBUG_PRINTHELP: opt.short[%s]=%p\n", o.ShortName, oo)
					}
				}

				if o.Required {
					w.Write([]byte(" (required)"))
				}

				if o.ValidDefault() {
					fmt.Fprintf(w, " (default: %v)", o.Default)
				}

				w.Write([]byte("\n"))
			} // for range g.options
			w.Write([]byte("\n\n"))
		}

		if len(g.commands) > 0 {
			if g.Name == GroupDefault {
				w.Write([]byte("Main commands:\n"))
			} else {
				fmt.Fprintf(w, "%s commands:\n", g.Name)
			}

			for _, cmd := range g.commands {
				fmt.Fprintf(w, "\t%s\t%s", cmd, opt.commands[cmd].Help)
				if len(opt.commands[cmd].Aliases) > 0 {
					fmt.Fprintf(w, " (aliases: ")
					for i, alias := range opt.commands[cmd].Aliases {
						if i == 0 {
							fmt.Fprintf(w, "%s", alias)
						} else {
							fmt.Fprintf(w, ",%s", alias)
						}
					}
					fmt.Fprintf(w, ")")
				}
				w.Write([]byte("\n"))
			}
			w.Write([]byte("\n\n"))
		}
	}

	if len(opt.positional) > 0 {
		w.Write([]byte("Positional arguments:\n"))
		for _, o := range opt.positional {
			fmt.Fprintf(w, "\t%s\t%s\n", o.Placeholder, o.Help)
		}
		w.Write([]byte("\n"))
	}
	w.Flush()
}
