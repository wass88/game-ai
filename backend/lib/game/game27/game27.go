package game27

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/wass88/gameai/lib/game"
	"github.com/wass88/gameai/lib/protocol"
	pr "github.com/wass88/gameai/lib/protocol"
)

type Game27 struct {
	board  [][]int
	first  bool
	record []string
}

const size = 9;

func NewGame27() *Game27 {
	board := [][]int{}
	for i := 0; i < size; i++ {
		board = append(board, []int{})
	}
	for i := 0; i < size; i++ {
		board[0] = append(board[0], 1)
		board[size-1] = append(board[size-1], -1)
	}
	return &Game27{board: board, first: true, record: []string{}}
}

const maxScore = 18

func (g *Game27) Start(players []*game.CmdRW, sender game.IPlayoutSender) (*pr.Result, error) {
	r0 := pr.ResultPlayer{0, "", ""}
	r1 := pr.ResultPlayer{0, "", ""}
	result := &pr.Result{Result: []pr.ResultPlayer{r0, r1}, Record: []string{}, Game: "Game27", Exception: ""}

	p0 := players[0]
	p1 := players[1]

	// TODO: Timeout Config...
	p0.WriteLn(fmt.Sprintf("init 0 %d", 5))
	p1.WriteLn(fmt.Sprintf("init 1 %d", 5))

	cp := p0
	op := p1
	cn := 0
	on := 1

	defer func() {
		result.Result[0].Stderr = string(p0.Stderr)
		result.Result[1].Stderr = string(p1.Stderr)
	}()
	for {
		s, err := cp.Wait()
		if err != nil {
			result.Exception = fmt.Sprintf("Failed By Player #%d", cn)
			result.Result[cn].Exception = fmt.Sprintf("P%d: Unexpected EOF: %s", cn, err.Error())
			result.Result[cn].Result = -maxScore
			result.Result[on].Result = maxScore
			return result, nil
		}
		result.Record = append(result.Record, s)
		err = sender.Update(protocol.ResultA{Record: strings.Join(result.Record, "\n"), Exception: result.Exception})
		if err != nil {
			return nil, errors.Wrapf(err, "On Update")
		}

		a, err := parseAction(s)
		if err != nil {
			result.Exception = fmt.Sprintf("Failed By Player #%d", cn)
			result.Result[cn].Exception = fmt.Sprintf("P%d: Unexpected Action: %s", cn, err.Error())
			result.Result[cn].Result = -maxScore
			result.Result[on].Result = maxScore
			return result, nil
		}
		err = g.act(*a)
		if err != nil {
			result.Exception = fmt.Sprintf("Failed By Player #%d", cn)
			result.Result[cn].Exception = fmt.Sprintf("P%d: Wrong Action: %s", cn, err.Error())
			result.Result[cn].Result = -maxScore
			result.Result[on].Result = maxScore
			return result, nil
		}
		op.WriteLn(fmt.Sprintf("played %s", s))
		fmt.Printf("%s", g.boardStr())
		if g.isEnd() {
			break
		}

		op, cp = cp, op
		on, cn = cn, on
	}
	res := g.result()
	result.Result[0].Result = res
	result.Result[1].Result = -res
	p0.WriteLn(fmt.Sprintf("result %d", res))
	p1.WriteLn(fmt.Sprintf("result %d", -res))
	return result, nil
}


type move struct {
	pass bool
	c int
	i int
}

func (g *Game27) Active() int {
	if g.first { return 1 } else { return -1 }
}
func isIn(c int) bool {
	return 0 <= c && c < size;
}
func (g *Game27) myTowers () int {
	res := 0;
	for c := 0; c < size; c++ {
		if len(g.board[c]) > 0 && g.board[c][0] == g.Active() {
			res++;
		}
	}
	return res;
}
func (g *Game27) moved(c int) int {
	return c + g.myTowers() * g.Active()
}
func (g *Game27) Playable() []move {
	res := []move{}
	for c := 0; c < size; c++ {
		if len(g.board[c]) > 0 && g.board[c][0] == g.Active() {
			for i := 1; i <= len(g.board[c]); i++ {
				j := g.moved(c)
				if isIn(j) {
					res = append(res, move{false, c, i})
				}
			}
		}
	}
	if len(res) == 0 {
		res = append(res, move{true, 0, 0})
	}
	return res
}
func (g *Game27) move(c, i int) error {
	if !isIn(c) {
		return fmt.Errorf("Out of board at %d", c)
	}
	if len(g.board[c]) == 0 {
		return fmt.Errorf("There is no tower at %d", c)
	}
	if g.board[c][0] != g.Active() {
		return fmt.Errorf("The tower is not yours")
	}
	if i <= 0 {
		return fmt.Errorf("The tower to move is zero (i = %d)", i)
	}
	if i > len(g.board[c]) {
		return fmt.Errorf("The tower to move is too high (i=%d, tower=%d)", i, len(g.board[c]))
	}
	j := g.moved(c)
	if !isIn(j) {
		return fmt.Errorf("The place the tower moves to is out of board (at %d)", j)
	}
	moved := make([]int, i)
	copy(moved, g.board[c][0:i])
	remain := g.board[c][i:len(g.board[c])]
	g.board[j] = append(moved, g.board[j]...)
	g.board[c] = remain
	g.first = !g.first
	return nil
}
func (g *Game27) pass() error {
	if !g.canPass() {
		return fmt.Errorf("It has playable moves")
	}
	g.first = !g.first
	return nil
}
func (g *Game27) canPass() bool {
	return len(g.Playable()) == 1 && g.Playable()[0].pass 
}
func (g *Game27) isEnd() bool {
	if g.canPass() {
		g.first = !g.first
		res := g.canPass()
		g.first = !g.first
		return res
	}
	return false
}

type action struct {
	pass bool
	c   int
	i   int
}

func parseAction(s string) (*action, error) {
	a := strings.Split(s, " ")
	if a[0] == "move" {
		a1, err := strconv.Atoi(a[1])
		if err != nil {
			return nil, err
		}
		a2, err := strconv.Atoi(a[2])
		if err != nil {
			return nil, err
		}
		return &action{false, a1, a2}, nil
	}
	if a[0] == "pass" {
		return &action{true, 0, 0,}, nil
	}
	return nil, fmt.Errorf("Unknown command: %s", s)
}

func (g *Game27) act(a action) error {
	if a.pass {
		return g.pass()
	} else {
		return g.move(a.c, a.i)
	}
}
func (g *Game27) result() int {
	f := len(g.board[size - 1])
	s := len(g.board[0])
	return f - s
}
func (g *Game27) boardStr() string {
	res := ""
	for c := 0; c < size; c++ {
		res += strconv.Itoa(c) + ":"
		for x := 0; x < len(g.board[c]); x++ {
			c := g.board[c][x]
			if c == 1 {
				res += "O"
			} else if c == -1 {
				res += "X"
			}
		}
		res += "\n"
	}
	return res
}
