package main

import (
	"flag"
	"os/exec"
	"strings"

	"net/http"

	"github.com/wass88/gameai/lib"
)

func main() {
	sendData := flag.String("send", "", "if set addr!id!token, send data")
	flag.Parse()
	args := flag.Args()
	game := args[0]
	send, err := lib.ParsePlayoutSender(*sendData, *http.NewClient())
	if err != nil {
		panic(err)
	}
	players := args[1:]
	cmds := []*exec.Cmd{}
	for _, player := range players {
		cmd := strings.Split(player, "")
		cmds = append(cmds, exec.Command(cmd[0], cmd[1:]...))
	}
	lib.StartPlayout(game, send, cmds)
}
