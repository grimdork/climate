package arg_test

import (
	"testing"

	"github.com/grimdork/climate/arg"
)

func TestPrintHelpShowsChoices(t *testing.T) {
	opt := arg.New("testapp", "")
	opt.SetDefaultHelp(true)
	op := []any{"dev", "prod", "test"}
	opt.SetOption(arg.GroupDefault, "m", "mode", "Run mode.", "dev", false, arg.VarString, op)
	opt.PrintHelp()
}
