package playout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	gi "github.com/wass88/gameai/lib/game"
	"github.com/wass88/gameai/lib/game/reversi"
	"github.com/wass88/gameai/lib/protocol"
	pr "github.com/wass88/gameai/lib/protocol"
)

func StartPlayout(gamename string, send gi.IPlayoutSender, cmds []*exec.Cmd) (*pr.Result, error) {
	gameSelector := NewGameSelector()
	gameSelector.Add("reversi", func() gi.Game { return reversi.NewReversi() })

	game := gameSelector.Get(gamename)
	ps := []*gi.CmdRW{}
	for _, cmd := range cmds {
		p, err := gi.RunWithReadWrite(cmd)
		if err != nil {
			return nil, fmt.Errorf("Failed Run: %v", cmd)
		}
		ps = append(ps, p)
	}
	result, err := game.Start(ps, send)
	if err != nil {
		return nil, errors.Wrapf(err, "On Start")
	}
	re := protocol.ResultA{Record: strings.Join(result.Record, "\n"), Exception: result.Exception}
	err = send.Update(re)
	if err != nil {
		return nil, errors.Wrapf(err, "On Update")
	}
	res := []protocol.ResultPlayerA{}
	for _, r := range result.Result {
		res = append(res, protocol.ResultPlayerA{Result: r.Result, Stderr: r.Stderr, Exception: r.Exception})
	}
	err = send.Complete(res)
	if err != nil {
		return nil, errors.Wrapf(err, "On Complete")
	}
	return result, err
}

type GameSelector struct {
	data map[string](func() gi.Game)
}

func NewGameSelector() GameSelector {
	return GameSelector{map[string]func() gi.Game{}}
}
func (g *GameSelector) Add(name string, game func() gi.Game) {
	g.data[name] = game
}

func (g *GameSelector) Get(name string) gi.Game {
	return g.data[name]()
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type PlayoutSender struct {
	API    string
	ID     string
	Token  string
	Client HttpClient
}

func (s *PlayoutSender) PostJson(url string, j interface{}) error {
	jsonBytes, err := json.Marshal(j)
	if err != nil {
		return errors.Wrapf(err, "Failed Marshal")
	}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return errors.Wrapf(err, "Failed Creating new request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Faild Post")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Faild ReadAll Body")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("Bad StatusCode: %d (%s)", resp.StatusCode, body)
	}
	return nil
}
func (s *PlayoutSender) URL(path string) string {
	return fmt.Sprintf("%s/results/%s/%s?token=%s", s.API, s.ID, path, s.Token)
}
func (s *PlayoutSender) Update(result protocol.ResultA) error {
	u := s.URL("update")
	err := s.PostJson(u, result)
	if err != nil {
		return errors.Wrapf(err, "Faild Update")
	}
	fmt.Printf("UPDATE COMPLETED\n")
	return nil
}
func (s *PlayoutSender) Complete(results []protocol.ResultPlayerA) error {
	u := s.URL("complete")
	err := s.PostJson(u, results)
	if err != nil {
		return errors.Wrapf(err, "Faild Completes")
	}
	return nil
}

type EmptySender struct{}

func (_ *EmptySender) Update(_ protocol.ResultA) error           { return nil }
func (_ *EmptySender) Complete(_ []protocol.ResultPlayerA) error { return nil }

func ParsePlayoutSender(s string, client HttpClient) (gi.IPlayoutSender, error) {
	if s == "" {
		return &EmptySender{}, nil
	}
	sp := strings.Split(s, "!")
	if len(sp) != 3 {
		return nil, errors.Errorf("Send format is invalied: %s", s)
	}
	return &PlayoutSender{API: sp[0], ID: sp[1], Token: sp[2], Client: client}, nil
}
