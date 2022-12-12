package arg_test

import (
	"testing"

	"github.com/grimdork/climate/arg"
)

const appname = "testrun"

func TestRemoveFirstGroup(t *testing.T) {
	opt := arg.New(appname)
	opt.AddGroup("Two")
	opt.AddGroup("Three")
	t.Logf("%v", opt)
	opt.RemoveGroup(arg.GroupDefault)
	t.Logf("%v", opt)
	if opt.GroupCount() != 2 {
		t.Errorf("Expected 2 groups, but got %d", opt.GroupCount())
		t.Fail()
	}

	g := opt.GetGroup(arg.GroupDefault)
	if g != nil {
		t.Errorf("Group 'default' should not exist, but does.")
		t.Fail()
	}
}

func TestRemoveMiddleGroup(t *testing.T) {
	opt := arg.New(appname)
	opt.AddGroup("Two")
	opt.AddGroup("Three")
	t.Logf("%v", opt)
	opt.RemoveGroup("Two")
	t.Logf("%v", opt)
	if opt.GroupCount() != 2 {
		t.Errorf("Expected 2 groups, but got %d", opt.GroupCount())
		t.Fail()
	}

	g := opt.GetGroup("Two")
	if g != nil {
		t.Errorf("Group 'Two' should not exist, but does.")
		t.Fail()
	}
}

func TestRemoveLastGroup(t *testing.T) {
	opt := arg.New(appname)
	opt.AddGroup("Two")
	opt.AddGroup("Three")
	t.Logf("%v", opt)
	opt.RemoveGroup("Three")
	t.Logf("%v", opt)
	if opt.GroupCount() != 2 {
		t.Errorf("Expected 2 groups, but got %d", opt.GroupCount())
		t.Fail()
	}

	g := opt.GetGroup("Three")
	if g != nil {
		t.Errorf("Group 'Three' should not exist, but does.")
		t.Fail()
	}
}

func TestSortGroup(t *testing.T) {
	opt := arg.New(appname)
	opt.SetOption("", "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	opt.SetOption("", "f", "file", "Full file path.", "", false, arg.VarString, nil)
	opt.SetOption("", "p", "port", "Port number.", 0, false, arg.VarInt, nil)
	opt.SetDefaultHelp(false)
	g := opt.GetGroup(arg.GroupDefault)
	if g == nil {
		t.Errorf("Group 'default' should exist, but does not.")
		t.FailNow()
	}

	t.Log("Unsorted:")
	for _, o := range g.GetOptions() {
		t.Logf("ShortName: %s LongName: %s", o.ShortName, o.LongName)
	}

	g.Sort()
	t.Log("Sorted:")
	list := g.GetOptions()
	for _, o := range list {
		t.Logf("ShortName: %s LongName: %s", o.ShortName, o.LongName)
	}

	if list[0].LongName != "help" || list[1].ShortName != "f" || list[2].ShortName != "p" || list[3].ShortName != "v" {
		t.Log("Sort order not the same as expected.")
		t.FailNow()
	}
}

func TestAutoGroup(t *testing.T) {
	opt := arg.New(appname)
	opt.SetOption("General", "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	list := opt.GetGroups()
	if list[1].Name != "General" {
		t.Errorf("Expected 'General' group, but got %s", list[1].Name)
		t.Fail()
	}
}

func TestLongShort(t *testing.T) {
	err := arg.New(appname).SetOption("", "verbose", "", "", false, false, arg.VarBool, nil)
	if err == nil {
		t.Errorf("Expected error, but long short worked.")
		t.Fail()
	} else {
		t.Log("Long short failed as expected.")
	}
}

func TestShortLong(t *testing.T) {
	err := arg.New(appname).SetOption("", "", "v", "", false, false, arg.VarBool, nil)
	if err == nil {
		t.Errorf("Expected error, but short long worked.")
		t.Fail()
	} else {
		t.Log("Short long failed as expected.")
	}
}

func TestBool(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"-v"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	if !opt.GetBool("v") {
		t.Errorf("Expected verbose to be true, but got false.")
		t.FailNow()
	} else {
		t.Log("Verbose is true as expected.")
	}
}

func TestString(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "f", "file", "Full file path.", nil, false, arg.VarString, nil)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}

	opt.PrintHelp()
	args := []string{"-f", "test.txt"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}

	if opt.GetString("f") != "test.txt" {
		t.Errorf("Expected 'test.txt', but got %s", opt.GetString("f"))
		t.Fail()
	} else {
		t.Log("File path is as expected.")
	}
}

func TestInt(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "p", "port", "Port number.", 3000, false, arg.VarInt, nil)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"-p", "4000"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	if opt.GetInt("p") != 4000 {
		t.Errorf("Expected -p=4000, but got %d", opt.GetInt("p"))
		opt.ShowOptions()
		t.Fail()
	} else {
		t.Log("Port number is as expected.")
	}
}

func TestFloat(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "p", "pi", "Your definition of pi.", 3.14, false, arg.VarFloat, nil)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"-p", "3.14159"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	if opt.GetFloat("p") != 3.14159 {
		t.Errorf("Expected -p=3.14159, but got %f", opt.GetFloat("p"))
		opt.ShowOptions()
		t.Fail()
	} else {
		t.Log("Pi is as expected.")
	}
}

const moo = `                 (__)
                 (oo)
           /------\/
          / |    ||
         *  /\---/\
            ~~   ~~
..."Have you mooed today?"...
`

func moocmd(opt *arg.Options) error {
	println(moo)
	return nil
}

func TestCommand(t *testing.T) {
	opt := arg.New(appname)
	_ = opt.SetCommand("moo", "Have you mooed today?", "", moocmd, nil)
	opt.PrintHelp()

	args := []string{"moo", "--help"}
	err := opt.Parse(args)
	if err != nil && err != arg.ErrRunCommand {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}
}

func TestPositional(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("FILE", "Full file path.", nil, false, arg.VarString)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"test.txt"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}

	s := opt.GetPosString("FILE")
	if s != "test.txt" {
		t.Errorf("Expected 'test.txt', but got %s", s)
		t.Fail()
	} else {
		t.Log("File path is as expected.")
	}
}

func TestPositionalStringSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("FILE", "Full file path.", nil, false, arg.VarStringSlice)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"test.txt", "test2.txt"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}

	files := opt.GetPosStringSlice("FILE")
	if len(files) != 2 {
		t.Errorf("Expected 2 files, but got %d", len(files))
		t.Fail()
	} else {
		t.Logf("File paths are as expected: %+v", files)
	}
}

func TestPositionalIntSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("PORT", "Port number.", 3000, false, arg.VarIntSlice)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"4000", "5000"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}

	ports := opt.GetPosIntSlice("PORT")
	if len(ports) != 2 {
		t.Errorf("Expected 2 ports, but got %d", len(ports))
		t.Fail()
	} else {
		t.Logf("Ports are as expected: %+v", ports)
	}
}

func TestPositionalFloatSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("PI", "Your definition of pi.", 3.14, false, arg.VarFloatSlice)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"3.14159", "3.1415926535"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.Fail()
	}

	pis := opt.GetPosFloatSlice("PI")
	if len(pis) != 2 {
		t.Errorf("Expected 2 pis, but got %d", len(pis))
		t.Fail()
	} else {
		t.Logf("Pis are as expected: %+v", pis)
	}
}

func TestCompletions(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "p", "pi", "Your definition of pi.", 3.14, false, arg.VarFloat, nil)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	cmd := opt.SetCommand("moo", "Have you mooed today?", "", moocmd, nil)
	if cmd == nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	comp, err := opt.Completions()
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	if comp == "" {
		t.Errorf("Expected completions, but got nothing")
		t.FailNow()
	}

	t.Logf("Completions:\n%s", comp)
}

func TestStringChoices(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "w", "word", "A number in word form.", "", false, arg.VarString, []any{"one", "two", "three"})
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args := []string{"-w", "one"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args = []string{"--word", "four"}
	err = opt.Parse(args)
	if err == nil {
		t.Errorf("Expected error for illegal choice, but got none")
		t.FailNow()
	}

	t.Log("String choices fail where expected.")
}

func TestIntChoices(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "n", "number", "An integer.", "", false, arg.VarInt, []any{1, 1, 1})
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args := []string{"-n", "1"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args = []string{"--number", "4"}
	err = opt.Parse(args)
	if err == nil {
		t.Errorf("Expected error for illegal choice, but got none")
		t.FailNow()
	}

	t.Log("Integer choices fail where expected.")
}

func TestFloatChoices(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "p", "pi", "Your definition of pi.", "", false, arg.VarFloat, []any{3.14, 3.141, 3.145})
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args := []string{"-p", "3.14"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args = []string{"--pi", "3.0"}
	err = opt.Parse(args)
	if err == nil {
		t.Errorf("Expected error for illegal choice, but got none")
		t.FailNow()
	}

	t.Log("Float choices fail where expected.")
}

func TestPositionalStringAndSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("WORD", "Just a word - any word.", "", false, arg.VarString)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	err = opt.SetPositional("ARGS", "The rest of the arguments.", "", false, arg.VarStringSlice)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	args := []string{"one", "two", "three"}
	err = opt.Parse(args)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	p := opt.GetPosString("WORD")
	if p != args[0] {
		t.Errorf("Expected '%s', but got '%s'", args[0], p)
		t.FailNow()
	}

	s := opt.GetPosStringSlice("ARGS")
	if len(s) != len(args[1:]) {
		t.Errorf("Expected %d args, but got %d", len(args[1:]), len(s))
		t.FailNow()
	}

	t.Logf("First arg (%s) = '%s'", args[0], p)
	for i, v := range s {
		if v != args[i+1] {
			t.Errorf("Expected '%s', but got '%s'", args[i+1], v)
			t.FailNow()
		} else {
			t.Logf("Arg %d (%s) = '%s'", i+2, args[i+1], v)
		}
	}

	t.Log("Positional strings work as expected.")
}
