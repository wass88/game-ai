package container

import (
	"io/ioutil"
	"os"
	"testing"
)

var dock = Dock{
	SaveDir: ".data/ai-docker/",
	API:     "http://api",
}

// https://github.com/wass88/reversi-random/archive/9703d869c9ce2245812ef657cbfe3210a7c33a86.zip
var commit = Commit{
	Github: "wass88/reversi-random",
	Branch: "master",
	Commit: "9703d869c9ce2245812ef657cbfe3210a7c33a86",
	ID:     1,
}

func TestSetup(t *testing.T) {
	id, err := dock.Setup(&commit)
	if err != nil {
		t.Fatal(err)
	}
	if id.ID != 1 {
		t.Fatalf("id.ID %d != 1", id.ID)
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
	e := "dir/wass88___reversi-random___master___9703d869c9ce2245812ef657cbfe3210a7c33a86"
	if p != e {
		t.Fatalf("%s is not %s", p, e)
	}
}

func TestGithub(t *testing.T) {
	p := commit.GithubZip()
	e := "https://github.com/wass88/reversi-random/archive/9703d869c9ce2245812ef657cbfe3210a7c33a86.zip"
	if p != e {
		t.Fatalf("\n%s is not \n%s", p, e)
	}
}

func TestImage(t *testing.T) {
	p := commit.Image()
	e := "ai-1"
	if p != e {
		t.Fatalf("%s is not %s", p, e)
	}
}
