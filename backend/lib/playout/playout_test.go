package playout

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	gi "github.com/wass88/gameai/lib/game"
	"github.com/wass88/gameai/lib/game/reversi"
	"github.com/wass88/gameai/lib/protocol"
)

func TestPlayout(t *testing.T) {
	prog := "../../../reversi-random/target/release/reversi_random"
	cmd0 := exec.Command(prog)
	cmd1 := exec.Command(prog)
	sender := &EmptySender{}
	p0, err := gi.RunWithReadWrite(cmd0)
	if err != nil {
		t.Fatal(err)
	}
	p1, err := gi.RunWithReadWrite(cmd1)
	if err != nil {
		t.Fatal(err)
	}
	r, err := reversi.NewReversi().Start([]*gi.CmdRW{p0, p1}, sender)
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
	err = sender.Update(protocol.ResultA{"put 1 2", ""})
	if err != nil {
		t.Fatal(err)
	}
	err = sender.Complete([]protocol.ResultPlayerA{{1, "", ""}, {-1, "", ""}})
	if err != nil {
		t.Fatal(err)
	}
}
