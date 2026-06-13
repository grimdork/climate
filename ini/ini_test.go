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

func TestAddAppendsSetReplacesFirst(t *testing.T) {
	ini, _ := New()

	// Add two values for the same key in a section
	ini.Add("sec", "k", "first")
	ini.Add("sec", "k", "second")

	if got := ini.GetString("sec", "k"); got != "first" {
		t.Fatalf("GetString should return the first value, got %q", got)
	}

	// Set should replace the first entry
	ini.Set("sec", "k", "replaced")
	if got := ini.GetString("sec", "k"); got != "replaced" {
		t.Fatalf("Set should replace first entry, got %q", got)
	}

	// Second entry should still be there
	fields := ini.GetMatch("sec", "k")
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
	if fields[1].Value != "second" {
		t.Fatalf("expected second entry to remain, got %q", fields[1].Value)
	}

	// Add another and verify three entries
	ini.Add("sec", "k", "third")
	fields = ini.GetMatch("sec", "k")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
}

func TestAddAppendsSetReplacesFirstTopLevel(t *testing.T) {
	ini, _ := New()

	ini.Add("", "k", "a")
	ini.Add("", "k", "b")
	ini.Add("", "k", "c")

	if got := ini.GetString("", "k"); got != "a" {
		t.Fatalf("GetString should return first value, got %q", got)
	}

	ini.Set("", "k", "x")
	if got := ini.GetString("", "k"); got != "x" {
		t.Fatalf("Set should replace first, got %q", got)
	}

	fields := ini.GetMatch("", "k")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
}

func TestGetMatchReturnsAll(t *testing.T) {
	ini, _ := New()
	ini.Add("", "top", "v1")
	ini.Add("", "top", "v2")
	ini.Add("db", "host", "primary")
	ini.Add("db", "host", "secondary")

	top := ini.GetMatch("", "top")
	if len(top) != 2 {
		t.Fatalf("expected 2 top-level matches, got %d", len(top))
	}
	if top[0].Value != "v1" || top[1].Value != "v2" {
		t.Fatalf("unexpected top values: %q, %q", top[0].Value, top[1].Value)
	}

	db := ini.GetMatch("db", "host")
	if len(db) != 2 {
		t.Fatalf("expected 2 db host matches, got %d", len(db))
	}
	if db[0].Value != "primary" || db[1].Value != "secondary" {
		t.Fatalf("unexpected db values: %q, %q", db[0].Value, db[1].Value)
	}

	// Non-existent key
	if ini.GetMatch("", "nope") != nil {
		t.Fatal("expected nil for missing top-level key")
	}
	if ini.GetMatch("nope", "x") != nil {
		t.Fatal("expected nil for missing section")
	}
}

func TestDuplicateKeyRoundtrip(t *testing.T) {
	ini, _ := New()
	ini.Add("", "top", "a")
	ini.Add("", "top", "b")
	ini.Add("sec", "x", "1")
	ini.Add("sec", "x", "2")
	ini.Add("sec", "y", "single")

	d := t.TempDir()
	fn := filepath.Join(d, "dup.ini")
	if err := ini.Save(fn, false); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(fn)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Top-level duplicates should be preserved
	top := loaded.GetMatch("", "top")
	if len(top) != 2 {
		t.Fatalf("expected 2 top matches after load, got %d", len(top))
	}
	if top[0].Value != "a" || top[1].Value != "b" {
		t.Fatalf("unexpected top values after load: %q, %q", top[0].Value, top[1].Value)
	}

	// Section duplicates should be preserved
	xs := loaded.GetMatch("sec", "x")
	if len(xs) != 2 {
		t.Fatalf("expected 2 sec.x matches after load, got %d", len(xs))
	}

	// GetString still returns first
	if got := loaded.GetString("sec", "x"); got != "1" {
		t.Fatalf("GetString should return first, got %q", got)
	}
}

func TestSetAndAddErrorOnEmptyKey(t *testing.T) {
	ini, _ := New()
	if err := ini.Set("", "", "val"); err == nil {
		t.Fatal("expected error for empty key in Set")
	}
	if err := ini.Add("", "", "val"); err == nil {
		t.Fatal("expected error for empty key in Add")
	}
	if err := ini.Set("sec", "", "val"); err == nil {
		t.Fatal("expected error for empty section key in Set")
	}
	if err := ini.Add("sec", "", "val"); err == nil {
		t.Fatal("expected error for empty section key in Add")
	}
}

func TestSaveTopLevelProperties(t *testing.T) {
	ini, _ := New()
	ini.Set("", "global", "gval")
	ini.Set("main", "host", "localhost")

	d := t.TempDir()
	fn := filepath.Join(d, "top.ini")
	if err := ini.Save(fn, false); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(fn)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got := loaded.GetString("", "global"); got != "gval" {
		t.Fatalf("expected global=gval, got %q", got)
	}
	if got := loaded.GetString("main", "host"); got != "localhost" {
		t.Fatalf("expected main.host=localhost, got %q", got)
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

func TestMarshalWithDuplicateKeys(t *testing.T) {
	ini, _ := New()
	ini.Add("", "top", "a")
	ini.Add("", "top", "b")
	ini.Add("sec", "x", "1")
	ini.Add("sec", "x", "2")
	ini.Add("sec", "y", "single")

	s := ini.Marshal()
	if !strings.Contains(s, "top = a\n") {
		t.Fatal("Marshal missing top=a")
	}
	if !strings.Contains(s, "top = b\n") {
		t.Fatal("Marshal missing top=b")
	}
	if !strings.Contains(s, "x = 1\n") {
		t.Fatal("Marshal missing x=1")
	}
	if !strings.Contains(s, "x = 2\n") {
		t.Fatal("Marshal missing x=2")
	}
	if !strings.Contains(s, "y = single\n") {
		t.Fatal("Marshal missing y=single")
	}

	// round-trip
	ini2, err := Unmarshal(s)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	top := ini2.GetMatch("", "top")
	if len(top) != 2 {
		t.Fatalf("expected 2 top matches, got %d", len(top))
	}
	xs := ini2.GetMatch("sec", "x")
	if len(xs) != 2 {
		t.Fatalf("expected 2 sec.x matches, got %d", len(xs))
	}
}

func TestSectionGetAll(t *testing.T) {
	ini, _ := New()
	ini.Add("sec", "k", "one")
	ini.Add("sec", "k", "two")
	ini.Add("sec", "k", "three")

	sec := ini.Sections["sec"]
	if sec == nil {
		t.Fatal("section not found")
	}

	fields := sec.GetAll("k")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[0].Value != "one" || fields[1].Value != "two" || fields[2].Value != "three" {
		t.Fatalf("unexpected field values: %q, %q, %q", fields[0].Value, fields[1].Value, fields[2].Value)
	}

	// Non-existent key
	if sec.GetAll("none") != nil {
		t.Fatal("expected nil for missing key")
	}
}

func TestINISetBoolIntFloat(t *testing.T) {
	ini, _ := New()
	if err := ini.SetBool("", "flag", true); err != nil {
		t.Fatalf("SetBool failed: %v", err)
	}
	if err := ini.SetInt("", "count", 99); err != nil {
		t.Fatalf("SetInt failed: %v", err)
	}
	if err := ini.SetFloat("", "ratio", 1.5); err != nil {
		t.Fatalf("SetFloat failed: %v", err)
	}
	if err := ini.SetBool("sec", "flag", false); err != nil {
		t.Fatalf("SetBool sec failed: %v", err)
	}

	if got := ini.GetBool("", "flag"); got != true {
		t.Fatalf("expected flag=true, got %v", got)
	}
	if got := ini.GetInt("", "count"); got != 99 {
		t.Fatalf("expected count=99, got %d", got)
	}
	if got := ini.GetFloat("", "ratio"); got != 1.5 {
		t.Fatalf("expected ratio=1.5, got %v", got)
	}
	if got := ini.GetBool("sec", "flag"); got != false {
		t.Fatalf("expected sec.flag=false, got %v", got)
	}
}

func TestINIAddBoolIntFloat(t *testing.T) {
	ini, _ := New()
	ini.AddBool("", "f", true)
	ini.AddBool("", "f", false)
	ini.AddInt("", "n", 10)
	ini.AddInt("", "n", 20)
	ini.AddFloat("", "x", 1.0)
	ini.AddFloat("", "x", 2.0)

	if got := ini.GetBool("", "f"); got != true {
		t.Fatalf("expected first f=true, got %v", got)
	}
	fields := ini.GetMatch("", "f")
	if len(fields) != 2 {
		t.Fatalf("expected 2 f entries, got %d", len(fields))
	}

	if got := ini.GetInt("", "n"); got != 10 {
		t.Fatalf("expected first n=10, got %d", got)
	}
	if got := ini.GetFloat("", "x"); got != 1.0 {
		t.Fatalf("expected first x=1.0, got %v", got)
	}

	// SetBool should replace first
	ini.SetBool("", "f", true)
	if got := ini.GetBool("", "f"); got != true {
		t.Fatalf("SetBool should keep first as true, got %v", got)
	}
}

func TestBoolFromString(t *testing.T) {
	if !BoolFromString("true") {
		t.Fatal("true should be truthy")
	}
	if !BoolFromString("yes") {
		t.Fatal("yes should be truthy")
	}
	if !BoolFromString("on") {
		t.Fatal("on should be truthy")
	}
	if !BoolFromString("enabled") {
		t.Fatal("enabled should be truthy")
	}
	if BoolFromString("false") {
		t.Fatal("false should be falsy")
	}
	if BoolFromString("off") {
		t.Fatal("off should be falsy")
	}
	if BoolFromString("0") {
		t.Fatal("0 should be falsy")
	}
	if BoolFromString("random") {
		t.Fatal("random should be falsy")
	}
}

func TestSectionAddSetErrorOnEmptyKey(t *testing.T) {
	ini, _ := New()
	sec := ini.AddSection("test")

	if err := sec.AddString("", "v"); err == nil {
		t.Fatal("expected error for empty key in AddString")
	}
	if err := sec.SetString("", "v"); err == nil {
		t.Fatal("expected error for empty key in SetString")
	}
	if err := sec.AddBool("", true); err == nil {
		t.Fatal("expected error for empty key in AddBool")
	}
	if err := sec.SetBool("", true); err == nil {
		t.Fatal("expected error for empty key in SetBool")
	}
	if err := sec.AddInt("", 1); err == nil {
		t.Fatal("expected error for empty key in AddInt")
	}
	if err := sec.SetInt("", 1); err == nil {
		t.Fatal("expected error for empty key in SetInt")
	}
	if err := sec.AddFloat("", 1.0); err == nil {
		t.Fatal("expected error for empty key in AddFloat")
	}
	if err := sec.SetFloat("", 1.0); err == nil {
		t.Fatal("expected error for empty key in SetFloat")
	}
}
