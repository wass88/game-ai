package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
)

func TestPlayout(t *testing.T) {
	cmd0 := exec.Command("/Users/admin/Documents/reversi-random/target/release/reversi_random")
	cmd1 := exec.Command("/Users/admin/Documents/reversi-random/target/release/reversi_random")
	sender := &EmptySender{}
	p0, err := RunWithReadWrite(cmd0)
	if err != nil {
		t.Fatal(err)
	}
	p1, err := RunWithReadWrite(cmd1)
	if err != nil {
		t.Fatal(err)
	}
	r, err := NewReversi().Start([]*CmdRW{p0, p1}, sender)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+#v", r)
	if r.Exception != "" {
		t.Fatal(r.Exception)
	}
	if r.Result[0].Exception != "" {
		t.Fatal(r.Result[0].Exception)
	}
	if r.Result[1].Exception != "" {
		t.Fatal(r.Result[1].Exception)
	}
	if r.Result[0].Result != -r.Result[1].Result {
		t.Fatalf("%d != - %d", r.Result[0].Result, r.Result[1].Result)
	}
}

func TestReversi(t *testing.T) {
	r := NewReversi()
	p := r.playable()
	if len(p) != 4 {
		t.Fatalf("Playable is 4 : %+#v", p)
	}
}

type mockClient struct{}

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Printf("%v\n", req)
	empty := ioutil.NopCloser(strings.NewReader(""))
	return &http.Response{StatusCode: 200, Body: empty}, nil
}

func TestParseSender(t *testing.T) {
	sender, err := ParsePlayoutSender("http://xxx!1!TOKEN", &mockClient{})
	if err != nil {
		t.Fatal(err)
	}
	err = sender.Update(ResultA{"put 1 2", ""})
	if err != nil {
		t.Fatal(err)
	}
	err = sender.Complete([]ResultPlayerA{{1, "", ""}, {-1, "", ""}})
	if err != nil {
		t.Fatal(err)
	}
}
