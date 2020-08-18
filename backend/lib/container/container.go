package container

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	unzip "github.com/artdarek/go-unzip"
	"github.com/pkg/errors"
)

/*
	Setup Container:
		Status found -> setup
		Kick Container Download and build image ai-#{id}
		<- Ready Container
		Status setup -> ready

	Exec Container:
		Exec ContainerId

	Purge Container:
		Remove out-dated files and image
		<- Purged Container
*/
type Dock struct {
	SaveDir string
	API     string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (d *Dock) Setup(c *Commit) error {
	err := c.Download(d.SaveDir)
	if err != nil {
		return errors.Wrap(err, "On Download")
	}
	err = c.Build()
	if err != nil {
		return errors.Wrapf(err, "On Build")
	}
	err = c.SendReady(d.API)
	if err != nil {
		return errors.Wrapf(err, "On Send")
	}

	return nil
}

func (d *Dock) Exec(containerID *int64) error {
	return nil
}
func (d *Dock) Purge(containerID *int64) error {
	return nil
}

type Commit struct {
	Github string
	Branch string
	Commit string
}

func (c *Commit) Dir(saveDir string) string {
	github := strings.Replace(c.Github, "/", "___", -1)
	d := fmt.Sprintf("%s___%s___%s", github, c.Branch, c.Commit)
	return path.Join(saveDir, d)
}

func (c *Commit) GithubZip() string {
	return fmt.Sprintf("https://github.com/%s/archive/%s.zip", c.Github, c.Commit)
}

const aiFormat = "ai-%d"

func (c *Commit) Image() string {
	github := strings.Replace(c.Github, "/", "___", -1)
	return fmt.Sprintf("ai___%s___%s___%s", github, c.Branch, c.Commit)
}

func DownloadZip(url string, dir string) (string, error) {
	fmt.Printf("Downloading %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "Error downloading")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.Errorf("Bad Status %d", resp.StatusCode)
	}

	zip := dir + ".zip"
	out, err := os.Create(zip)
	if err != nil {
		return "", errors.Wrapf(err, "Cannot Create Zip")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Error on Copy")
	}
	return zip, nil
}

func (c *Commit) Download(saveDir string) error {
	dir := c.Dir(saveDir)
	g := c.GithubZip()
	zip, err := DownloadZip(g, dir)
	if err != nil {
		return errors.Wrapf(err, "Error on Download")
	}
	uz := unzip.New(zip, zip[:len(zip)-4])
	err = uz.Extract()
	if err != nil {
		return errors.Wrapf(err, "Error on Extract")
	}
	return nil
}

func (c *Commit) Build() error {
	return nil
}

func (c *Commit) Exec() error {
	return nil
}

func (c *Commit) SendReady(API string) error {
	return nil
}
func (c *Commit) SendPurged(API string) error {
	return nil
}
