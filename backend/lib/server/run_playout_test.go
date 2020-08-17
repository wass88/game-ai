package server

import (
	"testing"

	"github.com/wass88/gameai/lib/protocol"
)

func getDB() *DB {
	dbname := `root:goodpassword@tcp(127.0.0.1:13306)/dev`
	return NewDB(dbname)
}

func mockDB() *DB {
	db := getDB()
	dropMock(db)
	setupGame(db)
	setupAI(db)
	setupPlayout(db)
	return db
}

func dropMock(db *DB) {
	db.DB.MustExec(`DELETE FROM playout_result
		WHERE playout_result.playout_id = 1`)
	db.DB.MustExec(`DELETE FROM playout_result_ai
		WHERE playout_result_ai.playout_id = 1`)
	db.DB.MustExec(`DELETE FROM playout_ai
	    WHERE playout_ai.playout_id = 1 OR playout_ai.ai_id = 1 OR playout_ai.playout_id = 2`)
	db.DB.MustExec(`DELETE FROM playout_ai
	    WHERE playout_id IN
			(SELECT id FROM playout WHERE game_id = 1)
		`)
	db.DB.MustExec(`DELETE FROM playout WHERE id = 1 OR game_id = 1`)
	db.DB.MustExec(`DELETE FROM ai WHERE id = 1`)
	db.DB.MustExec(`DELETE FROM ai WHERE id = 2`)
	db.DB.MustExec(`DELETE FROM ai_github WHERE id = 1`)
	db.DB.MustExec(`DELETE FROM ai_github WHERE id = 2`)
	db.DB.MustExec(`DELETE FROM user WHERE id = 1`)
	db.DB.MustExec(`DELETE FROM game WHERE id = 1`)
}
func setupGame(db *DB) {
	db.DB.MustExec(`
		INSERT INTO game (id, name)
		VALUE (1, "reversi")`)
}

func setupAI(db *DB) {
	db.DB.MustExec(`
		INSERT INTO user (id, name, twitter_token)
		VALUE (1, "test_user", "test_token")`)
	db.DB.MustExec(`
		INSERT INTO ai_github (id, user_id, github, branch)
		VALUE (1, 1, "git_addr", "master")`)
	db.DB.MustExec(`
		INSERT INTO ai_github (id, user_id, github, branch)
		VALUE (2, 1, "git_addr", "super")`)
	db.DB.MustExec(`
		INSERT INTO ai (id, ai_github_id, commit)
		VALUE (1, 1, "000001")`)
	db.DB.MustExec(`
		INSERT INTO ai (id, ai_github_id, commit)
		VALUE (2, 2, "000002")`)
}

func setupPlayout(db *DB) {
	db.DB.MustExec(`
		INSERT INTO playout (id, state, game_id, token)
		VALUES (1, "ready", 1, "TOKEN")
	`)
	db.DB.MustExec(`
		INSERT INTO playout_ai (id, ai_id, playout_id, turn)
		VALUES (1, 1, 1, 0)
	`)
	db.DB.MustExec(`
		INSERT INTO playout_ai (id, ai_id, playout_id, turn)
		VALUES (2, 2, 1, 1)
	`)
}

func TestNewPlayout(t *testing.T) {
	db := mockDB()
	_, err := db.NewPlayout(1, []int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePlayout(t *testing.T) {
	db := mockDB()
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
	db := mockDB()
	playoutID := PlayoutID{1, db}
	res := []protocol.ResultPlayerA{{-12, "stderr", ""}, {12, "stderr", ""}}
	err := playoutID.Complete(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOldestPlayout(t *testing.T) {
	db := mockDB()
	task, err := db.GetOldestTask()
	if err != nil {
		t.Fatal(err)
	}
	if task.Game != "reversi" {
		t.Fatalf("%s is not reversi", task.Game)
	}
	if task.PlayoutID.ID != 1 {
		t.Fatalf("%d is not 1", task.PlayoutID.ID)
	}
	if task.Token != "TOKEN" {
		t.Fatalf("%s is not TOKEN", task.Token)
	}
	if task.Players[0] != 1 {
		t.Fatalf("%d is not 1", task.Players[0])
	}
	if task.Players[1] != 2 {
		t.Fatalf("%d is not 2", task.Players[2])
	}
}

func TestRunPlayout(t *testing.T) {
	db := mockDB()
	task, err := db.GetOldestTask()
	if err != nil {
		t.Fatal(err)
	}
	err = task.PlayoutID.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckToken(t *testing.T) {
	db := mockDB()
	playoutID := PlayoutID{1, db}
	ok, err := playoutID.ValidateToken("TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("Token is not expected")
	}
}
