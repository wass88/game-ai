package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type Result struct {
	Game      string
	Record    []string
	Exception string
	Result    []ResultPlayer
}

type ResultPlayer struct {
	Result    int
	Stderr    string
	Exception string
}

func StartPlayout(gamename string, send IPlayoutSender, cmds []*exec.Cmd) (*Result, error) {
	gameSelector := NewGameSelector()
	gameSelector.Add("reversi", func() Game { return NewReversi() })
	game := gameSelector.Get(gamename)
	ps := []*CmdRW{}
	for _, cmd := range cmds {
		p, err := RunWithReadWrite(cmd)
		if err != nil {
			return nil, fmt.Errorf("Failed Run: %v", cmd)
		}
		ps = append(ps, p)
	}
	result, err := game.Start(ps, send)
	if err != nil {
		return nil, errors.Wrapf(err, "On Start")
	}
	re := ResultA{Record: strings.Join(result.Record, "\n"), Exception: result.Exception}
	err = send.Update(re)
	if err != nil {
		return nil, errors.Wrapf(err, "On Update")
	}
	res := []ResultPlayerA{}
	for _, r := range result.Result {
		res = append(res, ResultPlayerA{Result: r.Result, Stderr: r.Stderr, Exception: r.Exception})
	}
	err = send.Complete(res)
	if err != nil {
		return nil, errors.Wrapf(err, "On Complete")
	}
	return result, err
}

type GameSelector struct {
	data map[string](func() Game)
}

func NewGameSelector() GameSelector {
	return GameSelector{map[string]func() Game{}}
}
func (g *GameSelector) Add(name string, game func() Game) {
	g.data[name] = game
}

func (g *GameSelector) Get(name string) Game {
	return g.data[name]()
}

type Game interface {
	Start(players []*CmdRW, sender IPlayoutSender) (*Result, error)
}

type IPlayoutSender interface {
	Update(result ResultA) error
	Complete(results []ResultPlayerA) error
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
func (s *PlayoutSender) Update(result ResultA) error {
	u := s.URL("update")
	err := s.PostJson(u, result)
	if err != nil {
		return errors.Wrapf(err, "Faild Update")
	}
	fmt.Printf("UPDATE COMPLETED\n")
	return nil
}
func (s *PlayoutSender) Complete(results []ResultPlayerA) error {
	u := s.URL("complete")
	err := s.PostJson(u, results)
	if err != nil {
		return errors.Wrapf(err, "Faild Completes")
	}
	return nil
}

type EmptySender struct{}

func (_ *EmptySender) Update(_ ResultA) error           { return nil }
func (_ *EmptySender) Complete(_ []ResultPlayerA) error { return nil }

func ParsePlayoutSender(s string, client HttpClient) (IPlayoutSender, error) {
	if s == "" {
		return &EmptySender{}, nil
	}
	sp := strings.Split(s, "!")
	if len(sp) != 3 {
		return nil, errors.Errorf("Send format is invalied: %s", s)
	}
	return &PlayoutSender{API: sp[0], ID: sp[1], Token: sp[2], Client: client}, nil
}
