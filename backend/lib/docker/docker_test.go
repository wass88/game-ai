package docker

import (
	"archive/tar"
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	d, err := NewDocker()
	if err != nil {
		t.Fatal(err)
	}
	c := context.Background()
	c, cancel := context.WithTimeout(c, time.Second/2)
	defer cancel()
	err = d.Run(c, "reversi-random", 100, 256)
	if strings.Contains(err.Error(), " On Waiting") {
		t.Fatal(err)
	}
}

func TestBuild(t *testing.T) {
	d, err := NewDocker()
	if err != nil {
		t.Fatal(err)
	}
	c := context.Background()
	err = d.Build(c, "../../../reversi-random/", "reversi-random")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMakeTar(t *testing.T) {
	dir := "../../../reversi-random/"
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	buf, err := makeTar(dir)
	if err != nil {
		t.Fatal(err)
	}
	reader := tar.NewReader(buf)
	ok := false
	for {
		h, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if h.Name == "Dockerfile" {
			ok = true
		}
	}
	if !ok {
		t.Fatal("Missing Docker")
	}
}

func TestHasImage(t *testing.T) {
	d, err := NewDocker()
	if err != nil {
		t.Fatal(err)
	}
	c := context.Background()
	ok, err := d.HasImage(c, "reversi-random")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Missig image")
	}
	ok2, err := d.HasImage(c, "reversi-wwwwww")
	if err != nil {
		t.Fatal(err)
	}
	if ok2 {
		t.Fatal("Having image")
	}
}