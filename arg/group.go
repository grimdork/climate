package arg

import (
	"sort"
)

// Group definition.
type Group struct {
	// Name of the group.
	Name     string
	options  []*Option
	commands []string
}

// Len returns the number of options in the group.
func (g *Group) Len() int {
	return len(g.options)
}

// Less returns true if the first option comes before the second.
func (g *Group) Less(i, j int) bool {
	return g.options[i].ShortName < g.options[j].ShortName
}

// Swap two options in the group.
func (g *Group) Swap(i, j int) {
	g.options[i], g.options[j] = g.options[j], g.options[i]
}

// Sort options in the group alphabetically.
func (g *Group) Sort() {
	sort.Sort(g)
}

// GetOptions returns a slice of options.
func (g *Group) GetOptions() []*Option {
	return g.options
}
