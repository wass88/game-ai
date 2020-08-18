package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

type Images struct {
	ID string
}

func CreateImages(images []types.ImageSummary) *Images {
	return nil
}

func (d *Docker) ImageList(c context.Context) (*Images, error) {
	images, err := d.Client.ImageList(c, types.ImageListOptions{All: true})
	if err != nil {
		return nil, errors.Wrapf(err, "On ImageList")
	}
	return CreateImages(images), nil
}

func (d *Docker) Build(c context.Context, dir string, image string) error {
	return nil
}

func (d *Docker) Purge(c context.Context, image string) error {

}

func (d *Docker) Run(c context.Context, image string) error {
	// https://qiita.com/Tsuzu/items/39e3996bbfffe1d492aa
	conf := cont.Config{Image: image, Tty: false,
		AttachStdin: true, AttachStdout: true, AttachStderr: true,
		OpenStdin: true, StdinOnce: true}
	container, err := d.Client.ContainerCreate(c, &conf, nil, nil, "")
	if err != nil {
		return errors.Wrapf(err, "On Creating container %s", image)
	}

	options := types.ContainerStartOptions{}
	err = d.Client.ContainerStart(c, container.ID, options)
	if err != nil {
		return errors.Wrapf(err, "On Starting")
	}

	hijk, err := d.Client.ContainerAttach(c, container.ID,
		types.ContainerAttachOptions{Stream: true, Stdin: true})
	if err != nil {
		return errors.Wrapf(err, "On Stdin")
	}
	defer hijk.Conn.Close()
	defer hijk.Close()
	go func(hijk types.HijackedResponse) {
		io.Copy(hijk.Conn, os.Stdin)
	}(hijk)

	hijk, err = d.Client.ContainerAttach(c, container.ID,
		types.ContainerAttachOptions{Stream: true, Stdout: true})
	if err != nil {
		return errors.Wrapf(err, "On Stdout")
	}
	defer hijk.Conn.Close()
	defer hijk.Close()
	go func(hijk types.HijackedResponse) {
		io.Copy(os.Stdout, hijk.Conn)
	}(hijk)

	hijk, err = d.Client.ContainerAttach(c, container.ID,
		types.ContainerAttachOptions{Stream: true, Stderr: true})
	if err != nil {
		return errors.Wrapf(err, "On Stderr")
	}
	defer hijk.Conn.Close()
	defer hijk.Close()
	go func(hijk types.HijackedResponse) {
		io.Copy(os.Stderr, hijk.Conn)
	}(hijk)

	_, err = d.Client.ContainerWait(c, container.ID)
	if err != nil {
		return errors.Wrapf(err, "On Waiting")
	}
	return nil
}
