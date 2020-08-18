package container

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type mockClient struct {
}

var cid = "b3cd1a475dded156758005866761de51ee690607"

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Printf("%s\n", req.URL.String())
	body := ioutil.NopCloser(&bytes.Buffer{})
	return &http.Response{StatusCode: 200, Body: body}, nil
}

var cont = Cont{
	SaveDir: "../../.data/ai-docker/",
	API:     "http://api",
	Client:  &mockClient{},
}

var commit = Commit{
	Github: "wass88/reversi-random",
	Branch: "master",
	Commit: cid,
}

func TestSetup(t *testing.T) {
	err := cont.Setup(commit)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExec(t *testing.T) {
	ch := make(chan error, 1)
	go func() {
		err := cont.Exec(commit, Resource{CPUPersent: 100, MemoryMB: 256})
		ch <- err
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatal(err)
		}
		t.Fatal("Not Timeout")
	case <-time.After(3 * time.Second):
	}
}

func TestDownload(t *testing.T) {
	dir, err := ioutil.TempDir("", "download")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = commit.Download(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDir(t *testing.T) {
	p := commit.Dir("./dir")
	e := "dir/wass88.reversi-random.master." + cid
	if p != e {
		t.Fatalf("%s is not %s", p, e)
	}
}

func TestGithub(t *testing.T) {
	p := commit.GithubZip()
	e := "https://github.com/wass88/reversi-random/archive/" + cid + ".zip"
	if p != e {
		t.Fatalf("\n%s is not \n%s", p, e)
	}
}

func TestImage(t *testing.T) {
	p := commit.Image()
	e := "ai.wass88.reversi.random.master:" + cid
	if p != e {
		t.Fatalf("%s is not %s", p, e)
	}
}
