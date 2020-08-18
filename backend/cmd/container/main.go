package main

import (
	"context"

	"github.com/wass88/gameai/lib/docker"
)

func main() {
	d, err := docker.NewDocker()
	if err != nil {
		panic(err)
	}

	c := context.Background()
	err = d.Run(c, "random")
	if err != nil {
		panic(err)
	}
}
