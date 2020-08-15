package main

import (
	"github.com/wass88/gameai/lib"
	"os"
	"os/exec"
)

func main() {
	p0 := os.Args[1]
	p1 := os.Args[2]
	c0 := exec.Command(p0)
	c1 := exec.Command(p1)
	lib.Playout(c0, c1)
}