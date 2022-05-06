package arg

// Options base definition.
type Options struct {
	short      map[string]*Option
	long       map[string]*Option
	positional []*Option
	posmap     map[string]*Option
	groups     map[string]*Group
	commands   map[string]*Command
	// Order of groups.
	order []string
	// Remainder contains args not parsed as options, commands or positional args.
	Remainder []string
	appname   string
	// hashelp is true if default help is defined.
	hashelp bool
}

// New options instance.
// name - Name of the program.
func New(name string) *Options {
	opt := &Options{
		short:    make(map[string]*Option),
		long:     make(map[string]*Option),
		posmap:   make(map[string]*Option),
		groups:   make(map[string]*Group),
		commands: make(map[string]*Command),
		appname:  name,
	}

	opt.AddGroup("default")
	return opt
}

// GroupCount returns the number of groups.
func (opt *Options) GroupCount() int {
	return len(opt.order)
}

// AddGroup adds a new group. This ensures the order for help listing.
func (opt *Options) AddGroup(group string) *Group {
	opt.groups[group] = &Group{Name: group}
	opt.order = append(opt.order, group)
	return opt.groups[group]
}

// GetGroup returns a pointer to a group.
func (opt *Options) GetGroup(name string) *Group {
	if name == "" {
		return opt.groups["default"]
	}

	g := opt.groups[name]
	return g
}

// GetGroups returns a slice of groups.
func (opt *Options) GetGroups() []*Group {
	list := [](*Group){}
	for _, g := range opt.order {
		list = append(list, opt.groups[g])
	}
	return list
}

// RemoveGroup from map and order.
func (opt *Options) RemoveGroup(name string) {
	delete(opt.groups, name)
	for i, v := range opt.order {
		if v == name {
			opt.order = append(opt.order[:i], opt.order[i+1:]...)
			return
		}
	}
}

// GetOption returns a pointer to an option.
func (opt *Options) GetOption(name string) *Option {
	var o *Option
	var ok bool
	if len(name) > 1 {
		o, ok = opt.long[name]
	} else {
		o, ok = opt.short[name]
	}

	if !ok {
		return nil
	}

	return o
}

// GetBool returns a bool option's value.
func (opt *Options) GetBool(name string) bool {
	o := opt.GetOption(name)
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

// GetString returns a string option's value.
func (opt *Options) GetString(name string) string {
	o := opt.GetOption(name)
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

// GetStringSlice returns a string slice option's value.
func (opt *Options) GetStringSlice(name string) []string {
	o := opt.GetOption(name)
	if o == nil {
		return []string{}
	}

	if o.Value == nil {
		if o.Default != nil {
			return o.Default.([]string)
		}

		return []string{}
	}

	return o.Value.([]string)
}

// GetInt returns an int option's value.
func (opt *Options) GetInt(name string) int {
	o := opt.GetOption(name)
	if o == nil {
		return 0
	}

	if o.Value == nil {
		if o.Default != nil {
			return o.Default.(int)
		}

		return 0
	}

	return o.Value.(int)
}

// GetFloat returns a float option's value.
func (opt *Options) GetFloat(name string) float64 {
	o := opt.GetOption(name)
	if o == nil {
		return 0.0
	}

	if o.Value == nil {
		if o.Default != nil {
			return o.Default.(float64)
		}

		return 0.0
	}

	return o.Value.(float64)
}
