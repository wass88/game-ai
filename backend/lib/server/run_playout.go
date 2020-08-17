package server

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/wass88/gameai/lib/protocol"
)

func GetPlayoutIDAndCheckToken(c echo.Context, db *DB) (*PlayoutID, error) {
	idn := c.Param("id")
	id, err := strconv.Atoi(idn)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Unexpected id"))
	}
	token := c.QueryParam("token")
	playoutID := PlayoutID{int64(id), db}
	ok, err := playoutID.ValidateToken(token)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "DB error on validate token"))
	}
	if !ok {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Bad Token")
	}
	return &playoutID, nil

}
func HandlerResultsUpdate(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := GetPlayoutIDAndCheckToken(c, db)
		if err != nil {
			return err
		}
		result := new(protocol.ResultA)
		if err := c.Bind(result); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Json Parse Error")
		}
		err = id.Update(*result)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "DB error on Update"))
		}
		return c.String(http.StatusOK, "")
	}
}
func HandlerResultsComplete(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := GetPlayoutIDAndCheckToken(c, db)
		if err != nil {
			return err
		}
		result := new([]protocol.ResultPlayerA)
		if err := c.Bind(result); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Json Parse Error")
		}
		err = id.Complete(*result)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "DB error", err)
		}
		return c.String(http.StatusOK, "")
	}
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
	tok := []sql.NullString{}
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
	return tok[0].String == token, nil
}

func (r *PlayoutID) Update(result protocol.ResultA) error {
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

func (r *PlayoutID) Complete(results []protocol.ResultPlayerA) error {
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
