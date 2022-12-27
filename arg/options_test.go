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
		t.Fatalf("Expected 2 groups, but got %d", opt.GroupCount())
		t.FailNow()
	}

	g := opt.GetGroup(arg.GroupDefault)
	if g != nil {
		t.Fatalf("Group 'default' should not exist, but does.")
		t.FailNow()
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
		t.Fatalf("Expected 2 groups, but got %d", opt.GroupCount())
		t.FailNow()
	}

	g := opt.GetGroup("Two")
	if g != nil {
		t.Fatalf("Group 'Two' should not exist, but does.")
		t.FailNow()
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
		t.Fatalf("Expected 2 groups, but got %d", opt.GroupCount())
		t.FailNow()
	}

	g := opt.GetGroup("Three")
	if g != nil {
		t.Fatalf("Group 'Three' should not exist, but does.")
		t.FailNow()
	}
}

func TestSortGroup(t *testing.T) {
	opt := arg.New(appname)
	opt.SetOption(arg.GroupDefault, "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "f", "file", "Full file path.", "", false, arg.VarString, nil)
	opt.SetOption(arg.GroupDefault, "p", "port", "Port number.", 0, false, arg.VarInt, nil)
	opt.SetDefaultHelp(false)
	g := opt.GetGroup(arg.GroupDefault)
	if g == nil {
		t.Fatalf("Group 'default' should exist, but does not.")
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
	}
}

func TestGroupNames(t *testing.T) {
	opt := arg.New(appname)
	opt.SetOption("General", "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	list := opt.GetGroups()
	if list[1].Name != "General" {
		t.Fatalf("Expected 'General' group, but got %s", list[1].Name)
		t.FailNow()
	}
}

func TestLongShort(t *testing.T) {
	err := arg.New(appname).SetOption(arg.GroupDefault, "verbose", "", "", false, false, arg.VarBool, nil)
	if err == nil {
		t.Fatalf("Expected error, but long short worked.")
		t.FailNow()
	} else {
		t.Log("Long short failed as expected.")
	}
}

func TestShortLong(t *testing.T) {
	err := arg.New(appname).SetOption(arg.GroupDefault, "", "v", "", false, false, arg.VarBool, nil)
	if err == nil {
		t.Fatalf("Expected error, but short long worked.")
		t.FailNow()
	} else {
		t.Log("Short long failed as expected.")
	}
}

func TestBool(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"-v"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	if !opt.GetBool("v") {
		t.Fatalf("Expected verbose to be true, but got false.")
	} else {
		t.Log("Verbose is true as expected.")
	}
}

func TestString(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "f", "file", "Full file path.", nil, false, arg.VarString, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	opt.PrintHelp()
	args := []string{"-f", "test.txt"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	if opt.GetString("f") != "test.txt" {
		t.Fatalf("Expected 'test.txt', but got %s", opt.GetString("f"))
		t.FailNow()
	} else {
		t.Log("File path is as expected.")
	}
}

func TestInt(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "p", "port", "Port number.", 3000, false, arg.VarInt, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"-p", "4000"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	if opt.GetInt("p") != 4000 {
		t.Fatalf("Expected -p=4000, but got %d", opt.GetInt("p"))
		opt.ShowOptions()
		t.FailNow()
	} else {
		t.Log("Port number is as expected.")
	}
}

func TestFloat(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "p", "pi", "Your definition of pi.", 3.14, false, arg.VarFloat, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"-p", "3.14159"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	if opt.GetFloat("p") != 3.14159 {
		t.Fatalf("Expected -p=3.14159, but got %f", opt.GetFloat("p"))
		opt.ShowOptions()
		t.FailNow()
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
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}
}

func TestPositional(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("FILE", "Full file path.", nil, false, arg.VarString)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"test.txt"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	s := opt.GetPosString("FILE")
	if s != "test.txt" {
		t.Fatalf("Expected 'test.txt', but got %s", s)
		t.FailNow()
	} else {
		t.Log("File path is as expected.")
	}
}

func TestPositionalStringSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("FILE", "Full file path.", nil, false, arg.VarStringSlice)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"test.txt", "test2.txt"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	files := opt.GetPosStringSlice("FILE")
	if len(files) != 2 {
		t.Fatalf("Expected 2 files, but got %d", len(files))
		t.FailNow()
	} else {
		t.Logf("File paths are as expected: %+v", files)
	}
}

func TestPositionalIntSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("PORT", "Port number.", 3000, false, arg.VarIntSlice)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"4000", "5000"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	ports := opt.GetPosIntSlice("PORT")
	if len(ports) != 2 {
		t.Fatalf("Expected 2 ports, but got %d", len(ports))
		t.FailNow()
	} else {
		t.Logf("Ports are as expected: %+v", ports)
	}
}

func TestPositionalFloatSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("PI", "Your definition of pi.", 3.14, false, arg.VarFloatSlice)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	opt.PrintHelp()
	args := []string{"3.14159", "3.1415926535"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
		t.FailNow()
	}

	pis := opt.GetPosFloatSlice("PI")
	if len(pis) != 2 {
		t.Fatalf("Expected 2 pis, but got %d", len(pis))
		t.FailNow()
	} else {
		t.Logf("Pis are as expected: %+v", pis)
	}
}

func TestCompletions(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "p", "pi", "Your definition of pi.", 3.14, false, arg.VarFloat, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	cmd := opt.SetCommand("moo", "Have you mooed today?", "", moocmd, nil)
	if cmd == nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	comp, err := opt.Completions()
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	if comp == "" {
		t.Fatalf("Expected completions, but got nothing")
	}

	t.Logf("Completions:\n%s", comp)
}

func TestStringChoices(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption("", "w", "word", "A number in word form.", "", false, arg.VarString, []any{"one", "two", "three"})
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args := []string{"-w", "one"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args = []string{"--word", "four"}
	err = opt.Parse(args)
	if err == nil {
		t.Fatalf("Expected error for illegal choice, but got none")
	}

	t.Log("String choices fail where expected.")
}

func TestIntChoices(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "n", "number", "An integer.", "", false, arg.VarInt, []any{1, 1, 1})
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args := []string{"-n", "1"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args = []string{"--number", "4"}
	err = opt.Parse(args)
	if err == nil {
		t.Fatalf("Expected error for illegal choice, but got none")
	}

	t.Log("Integer choices fail where expected.")
}

func TestFloatChoices(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "p", "pi", "Your definition of pi.", "", false, arg.VarFloat, []any{3.14, 3.141, 3.145})
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args := []string{"-p", "3.14"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args = []string{"--pi", "3.0"}
	err = opt.Parse(args)
	if err == nil {
		t.Fatalf("Expected error for illegal choice, but got none")
	}

	t.Log("Float choices fail where expected.")
}

func TestPositionalStringAndSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("WORD", "Just a word - any word.", "", false, arg.VarString)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	err = opt.SetPositional("ARGS", "The rest of the arguments.", "", false, arg.VarStringSlice)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	args := []string{"one", "two", "three"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	p := opt.GetPosString("WORD")
	if p != args[0] {
		t.Fatalf("Expected '%s', but got '%s'", args[0], p)
	}

	s := opt.GetPosStringSlice("ARGS")
	if len(s) != len(args[1:]) {
		t.Fatalf("Expected %d args, but got %d", len(args[1:]), len(s))
	}

	t.Logf("First arg (%s) = '%s'", args[0], p)
	for i, v := range s {
		if v != args[i+1] {
			t.Fatalf("Expected '%s', but got '%s'", args[i+1], v)
		} else {
			t.Logf("Arg %d (%s) = '%s'", i+2, args[i+1], v)
		}
	}

	t.Log("Positional strings work as expected.")
}

func TestPositionalNoSlice(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetPositional("WORD", "Just a word - any word.", "", false, arg.VarString)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	err = opt.SetPositional("ARGS", "The rest of the arguments.", "", false, arg.VarString)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	args := []string{"one"}
	err = opt.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error parsing options, but got %s", err.Error())
	}

	p := opt.GetPosString("WORD")
	if p != args[0] {
		t.Fatalf("Expected '%s', but got '%s'", args[0], p)
	}

	s := opt.GetPosStringSlice("ARGS")
	if len(s) != 0 {
		t.Fatalf("Expected 0 args, but got %d", len(s))
	}

	t.Log("Positional string arguments with no supplied args work as expected.")
}

func TestDefault(t *testing.T) {
	opt := arg.New(appname)
	err := opt.SetOption(arg.GroupDefault, "n", "number", "An integer.", 0, false, arg.VarInt, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o := opt.GetOption("number")
	if o.ValidDefault() {
		t.Fatal("Expected default integer 0 to fail, but it's valid.")
	}

	err = opt.SetOption(arg.GroupDefault, "n", "numberlist", "A list of integers.", []int{}, false, arg.VarIntSlice, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o = opt.GetOption("numberlist")
	if o.ValidDefault() {
		t.Fatalf("Expected default int slice '%v' to fail, but it's valid.", o.Default)
	}

	err = opt.SetOption(arg.GroupDefault, "f", "float", "A float.", 0.0, false, arg.VarFloat, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o = opt.GetOption("float")
	if o.ValidDefault() {
		t.Fatal("Expected default float 0.0 to fail, but it's valid.")
	}

	err = opt.SetOption(arg.GroupDefault, "f", "floatlist", "A list of floats.", []float64{}, false, arg.VarFloatSlice, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o = opt.GetOption("floatlist")
	if o.ValidDefault() {
		t.Fatalf("Expected default float slice '%v' to fail, but it's valid.", o.Default)
	}

	err = opt.SetOption(arg.GroupDefault, "s", "string", "A string.", "", false, arg.VarString, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o = opt.GetOption("string")
	if o.ValidDefault() {
		t.Fatal("Expected default empty string to fail, but it's valid.")
	}

	err = opt.SetOption(arg.GroupDefault, "s", "stringlist", "A list of strings.", []string{}, false, arg.VarStringSlice, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o = opt.GetOption("stringlist")
	if o.ValidDefault() {
		t.Fatal("Expected default empty string slice to fail, but it's valid.")
	}

	err = opt.SetOption(arg.GroupDefault, "b", "bool", "A boolean.", false, false, arg.VarBool, nil)
	if err != nil {
		t.Fatalf("Expected no error creating option, but got %s", err.Error())
	}

	o = opt.GetOption("bool")
	if !o.ValidDefault() {
		t.Fatal("Expected default false to pass, but it's invalid.")
	}

	t.Log("Default value works as expected.")
}
