package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
)

type Docker struct {
	Client *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "On Init Docker-CLI")
	}
	return &Docker{cli}, nil
}

func (d *Docker) Build(c context.Context, dir string, image string) error {
	tar, err := makeTar(dir)
	if err != nil {
		return err
	}
	option := types.ImageBuildOptions{
		Tags: []string{image},
		// Wait to complete building
		SuppressOutput: true,
	}
	_, err = d.Client.ImageBuild(c, tar, option)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docker) Purge(c context.Context, image string) error {
	panic("not implemented")
}

func (d *Docker) Run(c context.Context, image string, cpuPersent, memoryMB int64) error {
	// https://qiita.com/Tsuzu/items/39e3996bbfffe1d492aa
	conf := cont.Config{
		Image:        image,
		Tty:          false,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
	}

	hostConf := cont.HostConfig{
		Resources: cont.Resources{
			// Nano CPU Second
			NanoCPUs: int64(10000000 * cpuPersent),
			// Memory Bytes
			Memory: int64(1024 * 1024 * memoryMB),
		},
	}
	container, err := d.Client.ContainerCreate(c, &conf, &hostConf, nil, "")
	if err != nil {
		return errors.Wrapf(err, "On Creating container %s", image)
	}

	options := types.ContainerStartOptions{}
	err = d.Client.ContainerStart(c, container.ID, options)
	if err != nil {
		return errors.Wrapf(err, "On Starting")
	}

	hijk, err := d.Client.ContainerAttach(c, container.ID,
		types.ContainerAttachOptions{Stream: true, Stdin: true, Stdout: true, Stderr: true})
	if err != nil {
		return errors.Wrapf(err, "On Stdin")
	}
	defer hijk.Conn.Close()
	defer hijk.Close()
	go func() {
		_, err := io.Copy(hijk.Conn, os.Stdin)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stdin, hijk.Conn)
		if err != nil {
			panic(err)
		}
	}()

	_, err = d.Client.ContainerWait(c, container.ID)
	if err != nil {
		return errors.Wrapf(err, "On Waiting")
	}
	return nil
}

func makeTar(path string) (io.Reader, error) {
	// https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
	buf := &bytes.Buffer{}
	tw := tar.NewWriter(buf)

	// TODO .dockerignore (now skip .git)

	fmt.Printf("make %s\n", path)
	err := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		// Skip symlink
		if fi.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		relfile, err := filepath.Rel(path, file)
		if err != nil {
			return errors.Wrapf(err, "Failed Rel")
		}
		if strings.HasPrefix(relfile, ".git/") {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return errors.Wrapf(err, "Failed info")
		}

		header.Name = filepath.ToSlash(relfile)
		if err := tw.WriteHeader(header); err != nil {
			return errors.Wrap(err, "Failed write header")
		}
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return errors.Wrapf(err, "Failed Open")
			}
			if _, err := io.Copy(tw, data); err != nil {
				return errors.Wrapf(err, "Failed Copy")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = tw.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "On Close")
	}
	return buf, err
}
