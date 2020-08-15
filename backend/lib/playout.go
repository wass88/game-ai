package lib

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Result    int
	Record    []string
	Game      string
	Exception string
}

type CmdRW struct {
	buf    *bufio.ReadWriter
	stderr []byte
}

func (r *CmdRW) WriteLn(s string) error {
	_, err := r.buf.WriteString(s + "\n")
	if err != nil {
		return err
	}
	err = r.buf.Flush()
	return err
}
func (r *CmdRW) ReadLn() (string, error) {
	l := []byte{}
	c := []byte{}
	p := true
	var err error
	for p {
		c, p, err = r.buf.ReadLine()
		if err != nil {
			return "", err
		}
		l = append(l, c...)
	}
	return string(l), nil
}

// RunWithReadWrite runs cmd and returns pipe
func RunWithReadWrite(c *exec.Cmd) (*CmdRW, error) {
	in, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	serr, err := c.StderrPipe()
	if err != nil {
		return nil, err
	}
	b := bufio.NewReadWriter(bufio.NewReader(out), bufio.NewWriter(in))
	res := &CmdRW{buf: b, stderr: []byte{}}
	go func() {
		b := make([]byte, 1024)
		for {
			k, err := serr.Read(b)
			if err != nil {
				panic(err)
			}
			res.stderr = append(res.stderr, b[:k]...)
		}
	}()
	err = c.Start()
	if err != nil {
		return nil, err
	}

	return res, nil

}

// Playout runs two players
func Playout(player0, player1 *exec.Cmd) (*Result, error) {
	fmt.Printf("Start...\n")
	p0, err := RunWithReadWrite(player0)
	if err != nil {
		return nil, err
	}
	p1, err := RunWithReadWrite(player1)
	if err != nil {
		return nil, err
	}
	defer func() {
		fmt.Printf("P0 stderr: %s\n", p0.stderr)
		fmt.Printf("P1 stderr: %s\n", p1.stderr)
	}()

	reversi := NewReversi()
	res, err := reversi.Start([]*CmdRW{p0, p1})

	return res, err
}

type Game interface {
	Start(players []*CmdRW) (*Result, error)
}

type Reversi struct {
	board  [][]int
	first  bool
	record []string
}

func NewReversi() *Reversi {
	board := [][]int{}
	for i := 0; i < 8; i++ {
		board = append(board, []int{0, 0, 0, 0, 0, 0, 0, 0})
	}
	board[4][3] = 1
	board[3][4] = 1
	board[3][3] = 2
	board[4][4] = 2
	return &Reversi{board: board, first: true, record: []string{}}
}

var d8y = [8]int{0, 1, 1, 1, 0, -1, -1, -1}
var d8x = [8]int{1, 1, 0, -1, -1, -1, 0, 1}

type point struct {
	y int
	x int
}

func (r *Reversi) playable() []point {
	res := []point{}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			rev := r.reversal(y, x)
			ok := false
			for _, v := range rev {
				if v > 0 {
					ok = true
					break
				}
			}
			if ok {
				res = append(res, point{y: y, x: x})
			}
		}
	}
	return res
}
func isIn(y, x int) bool {
	return 0 <= y && y < 8 && 0 <= x && x < 8
}
func (r *Reversi) reversal(y, x int) []int {
	res := []int{0, 0, 0, 0, 0, 0, 0, 0}
	if r.board[y][x] != 0 {
		return res
	}
	for d := 0; d < 8; d++ {
		ny := y + d8y[d]
		nx := x + d8x[d]
		if !isIn(ny, nx) {
			continue
		}
		c := r.board[ny][nx]
		if c == 0 {
			continue
		}
		m := 1
		if !r.first {
			m = 2
		}
		if c == m {
			continue
		}
		for t := 2; t < 8; t++ {
			ny = y + d8y[d]*t
			nx = x + d8x[d]*t
			if !isIn(ny, nx) {
				break
			}
			c := r.board[ny][nx]
			if c == 0 {
				break
			}
			if c == m {
				res[d] = t - 1
				break
			}
		}
	}
	return res
}
func (r *Reversi) put(y, x int) error {
	if !isIn(y, x) {
		return fmt.Errorf("Out of board %d, %d", y, x)
	}
	if r.board[y][x] != 0 {
		return fmt.Errorf("There is ocupied %d, %d", y, x)
	}
	m := 1
	if !r.first {
		m = 2
	}
	rev := r.reversal(y, x)
	ok := false
	for d, v := range rev {
		for i := 1; i <= v; i++ {
			ok = true
			ny := y + d8y[d]*i
			nx := x + d8x[d]*i
			r.board[ny][nx] = m
		}
	}
	if !ok {
		fmt.Printf("%+v\n", r.playable())
		return fmt.Errorf("No Revesible Piece %d, %d", y, x)
	}
	r.board[y][x] = m
	r.first = !r.first
	return nil
}
func (r *Reversi) pass() error {
	if len(r.playable()) != 0 {
		return fmt.Errorf("You have places which can reverse oponent's pieces")
	}
	r.first = !r.first
	return nil
}
func (r *Reversi) isEnd() bool {
	if len(r.playable()) != 0 {
		return false
	}
	r.first = !r.first
	p := len(r.playable())
	r.first = !r.first
	return p == 0
}

type action struct {
	put bool
	y   int
	x   int
}

func parseAction(s string) (action, error) {
	a := strings.Split(s, " ")
	if a[0] == "pass" {
		return action{false, 0, 0}, nil
	}
	if a[0] == "put" {
		a1, err := strconv.Atoi(a[1])
		if err != nil {
			return action{false, 0, 0}, err
		}
		a2, err := strconv.Atoi(a[2])
		if err != nil {
			return action{false, 0, 0}, err
		}
		return action{true, a1, a2}, nil
	}
	return action{false, 0, 0}, fmt.Errorf("Unknown command: %s", s)
}

func (r *Reversi) act(a action) error {
	if a.put {
		return r.put(a.y, a.x)
	} else {
		return r.pass()
	}
}
func (r *Reversi) result() int {
	f := 0
	s := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if r.board[y][x] == 1 {
				f++
			} else if r.board[y][x] == 2 {
				s++
			}
		}
	}
	if f == 0 {
		return -64
	}
	if s == 0 {
		return 64
	}
	return f - s
}
func (r *Reversi) boardStr() string {
	res := ""
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			c := r.board[y][x]
			if c == 0 {
				res += "."
			} else if c == 1 {
				res += "O"
			} else if c == 2 {
				res += "X"
			}
		}
		res += "\n"
	}
	return res
}

const timeout = 5 * time.Second

func (c *CmdRW) Wait() (string, error) {
	c.WriteLn("wait")
	type serr struct {
		res string
		err error
	}
	ch := make(chan serr, 1)
	go func() {
		res, err := c.ReadLn()
		ch <- serr{res, err}
	}()
	select {
	case res := <-ch:
		return res.res, res.err
	case <-time.After(timeout):
		return "", fmt.Errorf("Timeout %v", timeout)
	}
}
func (r *Reversi) Start(players []*CmdRW) (*Result, error) {
	result := &Result{Result: 0, Record: []string{}, Game: "Reversi", Exception: ""}

	p0 := players[0]
	p1 := players[1]

	p0.WriteLn("init 0")
	p1.WriteLn("init 1")

	cp := p0
	op := p1
	cn := 0
	for {
		fmt.Printf("Wait...P%d\n", cn)
		s, err := cp.Wait()
		if err != nil {
			result.Exception = fmt.Sprintf("P%d: Unexpected EOF: %s", cn, err.Error())
			result.Result = -64
			return result, nil
		}
		fmt.Printf("P%d: %v\n", cn, s)
		result.Record = append(result.Record, s)

		a, err := parseAction(s)
		if err != nil {
			result.Exception = fmt.Sprintf("P%d: Unexpected Action: %s", cn, err.Error())
			result.Result = -64
			return result, nil
		}
		err = r.act(a)
		if err != nil {
			result.Exception = fmt.Sprintf("P%d: Wrong Action: %s", cn, err.Error())
			result.Result = -64
			return result, nil
		}
		op.WriteLn(fmt.Sprintf("played %s", s))
		fmt.Printf("%s", r.boardStr())
		if r.isEnd() {
			break
		}

		op, cp = cp, op
		if cn == 0 {
			cn = 1
		} else {
			cn = 0
		}

	}
	result.Result = r.result()
	p0.WriteLn(fmt.Sprintf("result %d", result.Result))
	p1.WriteLn(fmt.Sprintf("result %d", -result.Result))
	return result, nil
}
