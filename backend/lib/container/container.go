package container

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/wass88/gameai/lib/docker"
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
type Cont struct {
	SaveDir string
	API     string
	Client  HttpClient
}

type Commit struct {
	Github string
	Branch string
	Commit string
}

type Resource struct {
	CPUPersent int64
	MemoryMB   int64
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (d *Cont) Setup(c Commit) error {
	fmt.Printf("Start Downlaod\n")
	err := c.Download(d.SaveDir)
	if err != nil {
		return errors.Wrap(err, "On Download")
	}
	fmt.Printf("Start Build\n")
	err = c.Build(d)
	if err != nil {
		return errors.Wrapf(err, "On Build")
	}
	fmt.Printf("Start SendReady\n")
	err = c.SendReady(d)
	if err != nil {
		return errors.Wrapf(err, "On Send")
	}
	fmt.Printf("Setup Completed\n")
	return nil
}

func (d *Cont) Exec(commit Commit, resource Resource) error {
	doc, err := docker.NewDocker()
	if err != nil {
		return err
	}
	image := commit.Image()
	c := context.Background()
	err = doc.Run(c, image, resource.CPUPersent, resource.MemoryMB)
	if err != nil {
		return err
	}
	return nil
}

func (d *Cont) Purge(containerID *int64) error {
	panic("Not implemented")
}

func (c *Commit) Dir(saveDir string) string {
	github := strings.Replace(c.Github, "/", ".", -1)
	d := fmt.Sprintf("%s.%s.%s", github, c.Branch, c.Commit)
	return path.Join(saveDir, d)
}

func (c *Commit) GithubZip() string {
	return fmt.Sprintf("https://github.com/%s/archive/%s.zip", c.Github, c.Commit)
}

func (c *Commit) Image() string {
	image := fmt.Sprintf("ai.%s.%s:%s", c.Github, c.Branch, c.Commit)
	image = strings.Replace(image, "/", ".", -1)
	image = strings.Replace(image, "-", ".", -1)
	image = strings.Replace(image, "_", ".", -1)
	return image
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func DownloadZip(url string, dir string) (string, bool, error) {
	zip := dir + ".zip"
	if Exists(zip) {
		fmt.Printf("Already downloaded %s\n", zip)
		return zip, true, nil
	}
	fmt.Printf("Downloading %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", false, errors.Wrapf(err, "Error downloading")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", false, errors.Errorf("Bad Status %d", resp.StatusCode)
	}

	out, err := os.Create(zip)
	if err != nil {
		return "", false, errors.Wrapf(err, "Cannot Create Zip")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", false, errors.Wrapf(err, "Error on Copy")
	}
	return zip, false, nil
}

func (c *Commit) Download(saveDir string) error {
	dir := c.Dir(saveDir)
	g := c.GithubZip()
	zip, already, err := DownloadZip(g, dir)
	if already {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "Error on Download")
	}
	err = Unzip(zip, saveDir)
	if err != nil {
		return errors.Wrapf(err, "Error on Extract")
	}
	return nil
}

func (c *Commit) Build(d *Cont) error {
	ctx := context.Background()
	doc, err := docker.NewDocker()
	if err != nil {
		return err
	}
	fmt.Printf("Build %s %s\n", c.Dir(d.SaveDir), c.Image())
	err = doc.Build(ctx, c.Dir(d.SaveDir), c.Image())
	if err != nil {
		return err
	}
	return nil
}

func (d *Cont) PostJson(u string, j interface{}) error {
	jsonBytes, err := json.Marshal(j)
	if err != nil {
		return errors.Wrapf(err, "Failed Marshal")
	}
	req, err := http.NewRequest(
		"POST",
		d.API+u,
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return errors.Wrapf(err, "Failed Creating new request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Faild Post")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Faild ReadAll Body")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("Bad StatusCode: %d (%s)", resp.StatusCode, body)
	}
	return nil
}

func (c *Commit) SendReady(d *Cont) error {
	err := d.PostJson("/container/ready", c)
	if err != nil {
		return errors.Wrapf(err, "Error sending")
	}
	return nil
}
func (c *Commit) SendPurged(d *Cont) error {
	panic("not implemented")
}

func Unzip(src, dest string) error {
	//https://github.com/artdarek/go-unzip/blob/master/unzip.go

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	zipName := filepath.Base(src)
	dirName := zipName[:len(zipName)-4]

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		fs := strings.Split(f.Name, string(os.PathSeparator))
		renamef := append([]string{dirName}, fs[1:]...)
		renamed := filepath.Join(renamef...)
		path := filepath.Join(dest, renamed)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: Illegal file path", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
