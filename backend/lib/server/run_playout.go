package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/wass88/gameai/lib/protocol"
)

func (db *DB) KickPlayout() error {
	fmt.Printf("[Playout] Start\n")
	t, err := db.GetOldestTask()
	if err != nil {
		return errors.Wrapf(err, "Kick Playout")
	}
	if t == nil {
		fmt.Printf("[Playout] No Task\n")
		return nil
	}
	err = t.SpownPlayout(db.Config)
	if err != nil {
		return errors.Wrapf(err, "Failed Playout %v", err)
	}
	fmt.Printf("[Playout] Complated")
	return nil
}

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
func HandlerAddMatch(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, err := db.GetSession(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Session Error %s", err))
		}
		var req struct {
			GameID int64   `json:"game_id" validate:"required"`
			AIID   []int64 `json:"ai_id" validate:"required"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Json Parse"))
		}
		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Json Validate"))
		}
		if len(req.AIID) != 2 {
			return echo.NewHTTPError(http.StatusBadRequest, "AI is two")
		}
		aiid := []AIID{}
		for _, ai := range req.AIID {
			aiid = append(aiid, AIID(ai))
		}
		id, err := db.CreatePlayout(GameID(req.GameID), aiid)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Create %v", err)
		}
		var res struct {
			ID int64 `json:"playout_id"`
		}
		res.ID = id
		return c.JSON(http.StatusOK, res)
	}
}

const TOKEN_LEN = 32

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

type PlayerCmd struct {
	ID     AIID
	Github string
	Branch string
	Commit string
}
type PlayoutTask struct {
	PlayoutID PlayoutID
	Token     string
	Game      string
	Players   []PlayerCmd
}

func (db *DB) GetOldestTask() (*PlayoutTask, error) {
	ais := []struct {
		PlayoutID int64  `db:"playout_id"`
		Token     string `db:"token"`
		Game      string `db:"game"`
		AIID      AIID   `db:"ai_id"`
		Github    string `db:"github"`
		Branch    string `db:"branch"`
		Commit    string `db:"commit"`
	}{}
	err := db.DB.Select(&ais, `
		SELECT p.id AS playout_id, p.token AS token, p.name AS game,
		playout_ai.ai_id AS ai_id, ai.commit AS commit,
		ai_github.branch AS branch, ai_github.github AS github
		FROM
		(SELECT playout.id, playout.token, game.name
			FROM playout
			INNER JOIN game ON game.id = playout.game_id
			INNER JOIN playout_ai ON playout_ai.playout_id = playout.id
			INNER JOIN ai ON ai.id = playout_ai.ai_id
			WHERE playout.state = "ready"
			GROUP BY playout.id
			HAVING min(ai.state) = max(ai.state) AND max(ai.state) = "ready"
			ORDER BY playout.created_at DESC
			LIMIT 1) AS p
		INNER JOIN playout_ai ON playout_ai.playout_id = p.id
		INNER JOIN ai ON ai.id = playout_ai.ai_id
		INNER JOIN ai_github ON ai_github.id = ai.ai_github_id
		ORDER BY playout_ai.turn ASC
	`)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed select playout")
	}
	if len(ais) == 0 {
		return nil, nil
	}
	players := []PlayerCmd{}
	for _, ai := range ais {
		players = append(players, PlayerCmd{ai.AIID, ai.Github, ai.Branch, ai.Commit})
	}
	res := PlayoutTask{
		PlayoutID{ais[0].PlayoutID, db},
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

func (t *PlayoutTask) SpownPlayout(c *Config) error {
	err := t.PlayoutID.Run()
	if err != nil {
		return errors.Wrapf(err, "Update DB Run playout")
	}
	cmd := t.Cmd(c)
	fmt.Printf("playout cmd = %v\n", strings.Join(cmd.Args, `   `))
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "Faild Run %v", cmd.Args)
	}
	return nil
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

func (db *DB) CreatePlayout(gameID GameID, ais []AIID) (int64, error) {
	// TODO: not rated playout match
	tx, err := db.DB.Beginx()
	if err != nil {
		return -1, errors.Wrapf(err, "Beginx")
	}
	token, err := GenerateRandomString(TOKEN_LEN)
	if err != nil {
		return -1, errors.Wrapf(err, "Generate Random")
	}
	res, err := tx.Exec(`
		INSERT INTO playout (game_id, state, token)
		VALUES (?, "ready", ?);
	`, gameID, token)
	playoutID, err := res.LastInsertId()
	if err != nil {
		return -1, errors.Wrapf(err, "Insert playout")
	}
	for turn, ai := range ais {
		_, err := tx.Exec(`
			INSERT INTO playout_ai (playout_id, ai_id, turn)
			VALUES (?, ?, ?)
		`, playoutID, ai, turn)
		if err != nil {
			return -1, errors.Wrapf(err, "Insert %d %d", ai, turn)
		}
	}
	err = tx.Commit()
	if err != nil {
		return -1, errors.Wrapf(err, "Commit Playout")
	}
	return playoutID, nil
}
