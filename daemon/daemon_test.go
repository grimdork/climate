package daemon

import "testing"

func TestErrNotRoot(t *testing.T) {
	if ErrNotRoot == nil {
		t.Fatal("expected non-nil ErrNotRoot")
	}
}

func TestBreakChannelReturns(t *testing.T) {
	ch := BreakChannel()
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
}
