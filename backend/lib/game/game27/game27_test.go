package game27

import "testing"

func TestGame27(t *testing.T) {
	r := NewGame27()
	p := r.Playable()
	if len(p) != 9 {
		t.Fatalf("Playable is 8 : %+#v", p)
	}
}
