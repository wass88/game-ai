package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/wass88/gameai/lib/server"
)

func main() {
	confFile := os.Getenv("CONF_FILE")
	conf, err := server.LoadConfig(confFile)
	if err != nil {
		panic(err)
	}
	db := conf.NewDB()

	cmds := os.Args
	var cmd string
	if len(cmds) >= 2 {
		cmd = cmds[1]

		if cmd == "setupai" {
			err := db.KickSetupAI()
			if err != nil {
				panic(err)
			}
		} else if cmd == "playout" {
			err := db.KickPlayout()
			if err != nil {
				panic(err)
			}
		} else if cmd == "autoplayout" { 
			err := server.NewAutoPlayout(db).Kick()
			if err != nil {
				panic(err)
			}
		} else {
			panic(fmt.Sprintf("Unknown command %s", cmd))
		}
	}
	if len(cmds) == 1 {
		fmt.Printf("Start Kick Tick")
		kicker := server.NewTaskKicker(db)
		kicker.Start()

		// Wait forever
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	}
}
