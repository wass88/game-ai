package lib

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

const TOKEN_LEN = 64

type PlayoutID struct {
	ID int64
	DB *DB
}

func (db *DB) NewPlayout(gameID int64, aiIDs []int64) (*PlayoutID, error) {
	token, err := GenerateRandomString(TOKEN_LEN)
	if err != nil {
		return nil, err
	}
	res, err := db.DB.Exec(
		`INSERT INTO playout (state, game_id, token) VALUES (?, ?, ?)`,
		"ready", gameID, token)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed creating playout")
	}
	playoutID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	for i, aiID := range aiIDs {
		_, err := db.DB.Exec(
			`INSERT INTO playout_ai (playout_id, ai_id, turn) VALUES (?, ?, ?)`,
			playoutID, aiID, i)
		if err != nil {
			return nil, errors.Wrapf(err, "Faild creating playout_ai")
		}
	}
	return &PlayoutID{playoutID, db}, nil
}

type PlayoutTask struct {
	PlayoutID PlayoutID
	Token     string
	Game      string
	Players   []int64
}

type PlayoutTaskM struct {
	PlayoutID int64  `db:"playout_id"`
	Token     string `db:"token"`
	Game      string `db:"game"`
	Player    int64  `db:"player"`
}

func (db *DB) GetOldestTask() (*PlayoutTask, error) {
	ais := []PlayoutTaskM{}
	err := db.DB.Select(&ais, `
		SELECT p.id AS playout_id, p.token AS token, p.name AS game, playout_ai.id AS player FROM
			(SELECT playout.id, playout.token, game.name
		FROM playout
		INNER JOIN game ON game.id = playout.game_id
		WHERE playout.state = "ready"
		ORDER BY playout.created_at DESC
		LIMIT 1) AS p
		INNER JOIN playout_ai ON playout_ai.playout_id = p.id
		ORDER BY playout_ai.turn ASC
	`)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed select playout")
	}
	if len(ais) == 0 {
		return nil, nil
	}
	players := []int64{}
	for _, ai := range ais {
		players = append(players, ai.Player)
	}

	res := PlayoutTask{
		PlayoutID{ais[0].Player, db},
		ais[0].Token,
		ais[0].Game,
		players,
	}
	return &res, nil
}

func (p *PlayoutID) Run() error {
	_, err := p.DB.DB.Exec(`
		UPDATE playout SET state = "running"
		WHERE playout.id = ?
	`, p.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed update playout to running")
	}
	return nil
}

func (t *PlayoutTask) SpownPlayout(runner RunnerConf) {
	t.PlayoutID.Run()
	runner.Run(*t)
}

func (r *PlayoutID) ValidateToken(token string) (bool, error) {
	type Token struct {
		token string `db:"token"`
	}
	tok := []Token{}
	err := r.DB.DB.Select(&tok, `
		SELECT playout.token AS token
		FROM playout
		WHERE playout.id = ?
	`, r.ID)
	if err != nil {
		return false, errors.Wrapf(err, "Failed Select")
	}
	if len(tok) < 1 {
		return false, errors.New("Missing Playout")
	}
	return tok[0].token == token, nil
}

type ResultA struct {
	Record    string
	Exception string
}

func (r *PlayoutID) Update(result ResultA) error {
	_, err := r.DB.DB.Exec(`
		INSERT INTO playout_result (playout_id, record, exception)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE record=?, exception=?
	`, r.ID, result.Record, result.Exception, result.Record, result.Exception)
	if err != nil {
		return errors.Wrapf(err, "Failed insert playout_result")
	}
	return nil
}

type ResultPlayerA struct {
	Result    int
	Stderr    string
	Exception string
}

func (r *PlayoutID) Complete(results []ResultPlayerA) error {
	for i, result := range results {
		_, err := r.DB.DB.Exec(`
			INSERT INTO playout_result_ai (turn, playout_id, result, stderr, exception)
			VALUES (?, ?, ?, ?, ?)
		`, i, r.ID, result.Result, result.Stderr, result.Exception)
		if err != nil {
			return errors.Wrapf(err, "Failed insert playout_result_ai")
		}
	}
	_, err := r.DB.DB.Exec(`
		UPDATE playout SET state="completed" WHERE playout.id=?
	`, r.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed update playout to completed")
	}
	return nil
}
