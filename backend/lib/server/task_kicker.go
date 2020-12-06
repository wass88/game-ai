package server

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Kick() error
}

type TaskKick struct {
	Task   Task
	Period time.Duration
	Last   time.Time
}

type TaskKicker struct {
	tasks []*TaskKick
}

type taskSetupAI DB

func (db *taskSetupAI) Kick() error {
	return (*DB)(db).KickSetupAI()
}

type taskPlayout DB

func (db *taskPlayout) Kick() error {
	return (*DB)(db).KickPlayout()
}

func NewTaskKicker(db *DB) *TaskKicker {
	tasks := []*TaskKick{
		NewTaskKick((*taskSetupAI)(db), time.Second*600), // Github Rate...
		NewTaskKick((*taskPlayout)(db), time.Second*10),
		NewTaskKick(NewAutoPlayout(db), time.Second*10),
	}
	return &TaskKicker{tasks}
}

func (t *TaskKicker) Start() {
	for _, task := range t.tasks {
		go task.Start()
	}
}

func NewTaskKick(t Task, p time.Duration) *TaskKick {
	return &TaskKick{t, p, time.Now()}
}

func (t *TaskKick) Start() {
	for {
		t.Last = time.Now()
		w := sync.WaitGroup{}
		w.Add(1)
		go func() {
			defer w.Done()
			err := t.Task.Kick()
			if err != nil {
				fmt.Printf("Task Error %v\n", err)
			}
		}()
		w.Wait()
		delta := time.Now().Sub(t.Last)
		if delta < t.Period {
			time.Sleep(t.Period - delta)
		}
	}
}

