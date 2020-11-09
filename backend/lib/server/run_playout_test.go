package server

import (
	"testing"

	"github.com/wass88/gameai/lib/protocol"
)

func TestNewPlayout(t *testing.T) {
	db := mockPlayoutDB()
	_, err := db.NewPlayout(1, []int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePlayout(t *testing.T) {
	db := mockPlayoutDB()
	playoutID := PlayoutID{1, db}
	err := playoutID.Update(protocol.ResultA{"put 0 0", ""})
	if err != nil {
		t.Fatal(err)
	}
	err = playoutID.Update(protocol.ResultA{"put 1 1", ""})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletePlayout(t *testing.T) {
	db := mockPlayoutDB()
	playoutID := PlayoutID{1, db}
	res := []protocol.ResultPlayerA{{-12, "stderr", ""}, {12, "stderr", ""}}
	err := playoutID.Complete(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOldestPlayout(t *testing.T) {
	db := mockPlayoutDB()
	task, err := db.GetOldestPlayoutTask()
	if err != nil {
		t.Fatal(err)
	}
	if task.Game != "reversi" {
		t.Fatalf("%s is not reversi", task.Game)
	}
	if task.PlayoutID.ID != 2 {
		t.Fatalf("%d is not 2", task.PlayoutID.ID)
	}
	if task.Token != "TOKEN" {
		t.Fatalf("%s is not TOKEN", task.Token)
	}
	if task.Players[0].ID != 2 {
		t.Fatalf("%v is not 2", task.Players[0])
	}
	if task.Players[1].ID != 2 {
		t.Fatalf("%v is not 2", task.Players[1])
	}
}

func TestRunPlayout(t *testing.T) {
	db := mockPlayoutDB()
	task, err := db.GetOldestPlayoutTask()
	if err != nil {
		t.Fatal(err)
	}
	err = task.PlayoutID.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckToken(t *testing.T) {
	db := mockPlayoutDB()
	playoutID := PlayoutID{1, db}
	ok, err := playoutID.ValidateToken("TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("Token is not expected")
	}
}

func TestCreatePlayout(t *testing.T) {
	db := mockPlayoutDB()
	id, err := db.CreatePlayout(1, []AIID{1, 1})
	if err != nil {
		t.Fatal(err)
	}
	if id < 0 {
		t.Fatalf("id %d < 0 ", id)
	}
}
