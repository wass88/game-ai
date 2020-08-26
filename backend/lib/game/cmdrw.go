package game

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type CmdRW struct {
	in     io.Writer
	out    *bufio.Reader
	stderr []byte
}

func (r *CmdRW) WriteLn(s string) error {
	fmt.Printf("--> %s\n", s)
	_, err := io.WriteString(r.in, s+"\n")
	if err != nil {
		return err
	}
	return nil
}
func (r *CmdRW) ReadLn() (string, error) {
	l := []byte{}
	c := []byte{}
	p := true
	var err error
	for p {
		c, p, err = r.out.ReadLine()
		if err != nil {
			return "", err
		}
		l = append(l, c...)
	}
	fmt.Printf("<-- %s\n", string(l))
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
	res := &CmdRW{out: bufio.NewReader(out), in: in, stderr: []byte{}}
	go func() {
		b := make([]byte, 1024)
		for {
			k, err := serr.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
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
