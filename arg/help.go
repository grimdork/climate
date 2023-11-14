package arg

import (
	"fmt"
	"os"
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
					fmt.Fprintf(w, ")\n")
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
