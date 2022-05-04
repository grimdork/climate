package paths_test

import (
	"testing"

	"github.com/grimdork/cliapp/paths"
)

func TestPaths(t *testing.T) {
	cp, err := paths.New("test")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}

	t.Logf("Base: %s", cp.UserBase)
	t.Logf("ServerBase: %s", cp.ServerBase)
}
