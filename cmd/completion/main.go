package main

import (
	"fmt"
	"os"

	"github.com/grimdork/climate/arg"
)

func main() {
	opt := arg.New("rep")
	opt.SetDefaultHelp(true)
	opt.SetOption("", "f", "float", "A floating point number.", nil, false, arg.VarFloat, nil)
	cmd := opt.SetCommand("one", "Command one.", "", one, nil)
	cmd.Options.SetOption("", "i", "int", "An integer.", nil, false, arg.VarInt, nil)
	comp, err := opt.Completions()
	if err != nil {
		panic(err)
	}

	err = opt.ParseEnvironment("", "")
	if err != nil {
		panic(err)
	}

	err = opt.Parse(os.Args)
	if err != nil && err != arg.ErrNoArgs {
		if err == arg.ErrRunCommand {
			return
		}

		panic(err)
	}

	fmt.Printf("%s\n", comp)
}

func one(opt *arg.Options) error {
	opt.PrintHelp()
	return nil
}
