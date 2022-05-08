package arg

// Command definition.
type Command struct {
	// Name of the command.
	Name string
	// Help text of the command.
	Help string
	// Func to execute the command.
	Func ToolCommand
	// Options for this command.
	Options *Options
	// Aliases for this command.
	Aliases []string
}

// ToolCommand function signature.
type ToolCommand func(*Options) error

// SetCommand to a group.
func (opt *Options) SetCommand(name, help, group string, fn ToolCommand, aliases []string) *Command {
	cmd := &Command{
		Name:    name,
		Help:    help,
		Func:    fn,
		Options: New(opt.appname + " " + name),
		Aliases: aliases,
	}

	opt.commands[name] = cmd
	g := opt.GetGroup(group)
	if g == nil {
		g = opt.AddGroup(group)
	}
	g.commands = append(g.commands, cmd.Name)
	return cmd
}
