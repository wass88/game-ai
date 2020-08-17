package reversi

import "testing"

func TestReversi(t *testing.T) {
	r := NewReversi()
	p := r.Playable()
	if len(p) != 4 {
		t.Fatalf("Playable is 4 : %+#v", p)
	}
}
