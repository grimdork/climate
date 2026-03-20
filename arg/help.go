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
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 8, 8, 1, '\t', 0)
	w.Write([]byte("Usage:\n  "))
	w.Write([]byte(opt.appname))

	// Print inline application description if available (same line as appname).
	if opt.appdesc != "" {
		w.Write([]byte(" — "))
		w.Write([]byte(opt.appdesc))
	}

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
					// Pretty-print long choice lists on subsequent indented lines.
					// Aim for at most 8 items per line, and keep line length <= 80 if possible.
					maxPerLine := 8
					if maxPerLine > len(parts) {
						maxPerLine = len(parts)
					}
					chosen := 1
					// try from maxPerLine down to 1 to find the largest per-line count
					// that keeps every constructed line <= 80 chars (approx)
					for n := maxPerLine; n >= 1; n-- {
						// build sample lines
						ok := true
						for i := 0; i < len(parts); i += n {
							end := i + n
							if end > len(parts) {
								end = len(parts)
							}
							line := strings.Join(parts[i:end], ",")
							if len(line) > 80 {
								ok = false
								break
							}
						}
						if ok {
							chosen = n
							break
						}
					}
					// Print choices on following lines indented under the option text.
					first := true
					for i := 0; i < len(parts); i += chosen {
						end := i + chosen
						if end > len(parts) {
							end = len(parts)
						}
						if first {
							fmt.Fprintf(w, "\n\t\t  (choices: %s)", strings.Join(parts[i:end], ","))
							first = false
						} else {
							fmt.Fprintf(w, "\n\t\t    %s", strings.Join(parts[i:end], ","))
						}
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
