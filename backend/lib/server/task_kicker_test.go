package server

import (
	"testing"
	"time"
)

type mockTask struct {
	ch    chan int
	wait  time.Duration
	start time.Time
	t     *testing.T
}

func (e *mockTask) Kick() error {
	e.t.Logf("Kick   %v\n", time.Now().Sub(e.start))
	time.Sleep(e.wait)
	e.t.Logf("Finish %v\n", time.Now().Sub(e.start))
	e.ch <- 1
	return nil
}

func TestTaskKick(t *testing.T) {
	ch := make(chan int, 1)

	r := func(k *TaskKick) {
		go k.Start()
		i := 0
		for {
			select {
			case <-ch:
				i++
				if i > 2 {
					return
				}
			}
		}
	}
	t.Log("wait 10")
	r(NewTaskKick(&mockTask{ch, time.Millisecond * 10, time.Now(), t}, time.Millisecond*20))
	t.Log("wait 30")
	r(NewTaskKick(&mockTask{ch, time.Millisecond * 30, time.Now(), t}, time.Millisecond*20))
}
