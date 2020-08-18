package main

import (
	"flag"
	"net/http"

	"github.com/wass88/gameai/lib/container"
)

func main() {
	github := flag.String("github", "", "Github user/repo")
	branch := flag.String("branch", "", "Branch name")
	commit := flag.String("commit", "", "commit")
	api := flag.String("api", "", "api + /container/")
	dir := flag.String("dir", "", "Save Dir")
	cpu := flag.Int64("cpu", 0, "cpu persent")
	mem := flag.Int64("mem", 0, "memory MB")
	flag.Parse()
	if *github == "" || *branch == "" || *commit == "" {
		panic("Need --github --branch --commit")
	}
	c := container.Commit{Github: *github, Branch: *branch, Commit: *commit}
	args := flag.Args()
	cmd := args[0]
	cont := container.Cont{
		API:     *api,
		SaveDir: *dir,
		Client:  &http.Client{},
	}
	if cmd == "setup" {
		if *api == "" || *dir == "" {
			panic("Need --api --dir")
		}
		err := cont.Setup(c)
		if err != nil {
			panic(err)
		}
	} else if cmd == "run" {
		if *cpu == 0 || *mem == 0 {
			panic("Need --cpu --mem")
		}
		err := cont.Exec(c, container.Resource{
			CPUPersent: *cpu, MemoryMB: *mem,
		})
		if err != nil {
			panic(err)
		}
	} else {
		panic("unknown command: " + cmd)
	}
}
