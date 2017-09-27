package linda

import "testing"

func TestProcess_String(t *testing.T) {
	if p, err := newProcess("test"); err != nil {
		t.Error(err)
	} else {
		t.Log(p)
	}
}
