package main

import (
	"flag"
	"fmt"
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
	send, err := lib.ParsePlayoutSender(*sendData, &http.Client{})
	if err != nil {
		panic(err)
	}
	players := args[1:]
	cmds := []*exec.Cmd{}
	fmt.Printf("Game: %s\n", game)
	fmt.Printf("Send: %+#v\n", send)
	for i, player := range players {
		cmd := strings.Split(player, " ")
		cmds = append(cmds, exec.Command(cmd[0], cmd[1:]...))
		fmt.Printf("Player #%d: %s\n", i, player)
	}
	_, err = lib.StartPlayout(game, send, cmds)
	if err != nil {
		panic(err)
	}
}
