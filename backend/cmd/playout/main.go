package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"net/http"

	"github.com/wass88/gameai/lib/playout"
)

func main() {
	flag.Usage = func() {
		usageTxt := `Usage: playout [game name] [first player] [second player]

	Play the game in the players

	[game name]: "reversi" or "game27"
	[player]: "player command"
	-send [url]: send a playout to the server
	`
		fmt.Fprintf(os.Stderr, "%s\n", usageTxt)
	}
	sendData := flag.String("send", "", "if set addr@id@token, send data")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		panic("Select Game: reversi or game27")
	}
	game := args[0]
	send, err := playout.ParsePlayoutSender(*sendData, &http.Client{})
	if err != nil {
		panic(err)
	}
	players := args[1:]

	if len(players) != 2 {
		// TODO: more players
		panic("Need two players")
	}
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
	res, err := playout.StartPlayout(game, send, cmds)
	if err != nil {
		panic(err)
	}
	fmt.Printf("=============    Result   ===========\n")
	fmt.Printf("%v", res)
}
