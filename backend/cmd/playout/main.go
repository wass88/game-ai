package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"net/http"

	"github.com/wass88/gameai/lib/playout"
)

func main() {
	sendData := flag.String("send", "", "if set addr@id@token, send data")
	flag.Parse()
	args := flag.Args()
	game := args[0]
	send, err := playout.ParsePlayoutSender(*sendData, &http.Client{})
	if err != nil {
		panic(err)
	}
	players := args[1:]
	cmds := []*exec.Cmd{}
	fmt.Printf("============= Game Setting ===========\n")
	fmt.Printf("Game: %s\n", game)
	fmt.Printf("Send: %+#v\n", send)
	for i, player := range players {
		cmd := strings.Split(player, " ")
		cmds = append(cmds, exec.Command(cmd[0], cmd[1:]...))
		fmt.Printf("Player #%d: %s\n", i, player)
	}
	fmt.Printf("=============    Start     ===========\n")
	_, err = playout.StartPlayout(game, send, cmds)
	if err != nil {
		panic(err)
	}
}
