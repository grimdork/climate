package ini

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMarshalUnmarshalRoundtrip(t *testing.T) {
	ini, _ := New()
	ini.Set("", "global", "gval")
	ini.Set("main", "host", "localhost")
	ini.Set("main", "port", "8080")

	s := ini.Marshal()
	if strings.TrimSpace(s) == "" {
		t.Fatalf("Marshal produced empty output")
	}

	ini2, err := Unmarshal(s)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if ini2 == nil {
		t.Fatal("Unmarshal returned nil")
	}

	if got := ini2.GetString("main", "host"); got != "localhost" {
		t.Fatalf("expected host localhost, got %q", got)
	}
	if got := ini2.GetString("", "global"); got != "gval" {
		t.Fatalf("expected global gval, got %q", got)
	}
}

func TestGetEnvOverrideForceUpper(t *testing.T) {
	key := "HOST"
	os.Setenv(key, "envhost")
	defer os.Unsetenv(key)

	ini, _ := New()
	ini.SetSecure(false)
	ini.Set("main", "host", "filehost")

	// default: envFirst is false -> env key not used
	if got := ini.GetString("main", "host"); got != "filehost" {
		t.Fatalf("expected filehost without ForceUpper, got %q", got)
	}

	// Enable env-first lookup and force upper-case env keys
	ini.SetEnvFirst(true)
	ini.ForceUpper()
	if got := ini.GetString("main", "host"); got != "envhost" {
		t.Fatalf("expected envhost with ForceUpper and envFirst, got %q", got)
	}
}

func TestLoadNoTrailingNewline(t *testing.T) {
	d := t.TempDir()
	fn := filepath.Join(d, "no_nl.ini")
	_ = os.WriteFile(fn, []byte("key=value"), 0644)
	ini, err := Load(fn)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got := ini.GetString("", "key"); got != "value" {
		t.Fatalf("expected value, got %q", got)
	}
}

func TestSaveFilePerms(t *testing.T) {
	ini, _ := New()
	ini.Set("main", "a", "b")
	d := t.TempDir()
	fn := filepath.Join(d, "out.ini")
	// default (non-secure) -> 0644
	if err := ini.Save(fn, false); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	st, err := os.Stat(fn)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	// The OS umask may reduce permissions; ensure the actual perms are a subset of FilePerm.
	if st.Mode().Perm()&ini.FilePerm != st.Mode().Perm() {
		t.Fatalf("file permissions %o are not a subset of expected %o", st.Mode().Perm(), ini.FilePerm)
	}

	// secure -> 0600
	ini.SetSecure(true)
	fn2 := filepath.Join(d, "out2.ini")
	if err := ini.Save(fn2, false); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	st2, err := os.Stat(fn2)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if st2.Mode().Perm()&ini.FilePerm != st2.Mode().Perm() {
		t.Fatalf("file permissions %o are not a subset of expected %o", st2.Mode().Perm(), ini.FilePerm)
	}
}
