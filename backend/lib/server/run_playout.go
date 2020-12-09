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

// KickPlayout kicks the oldest unplayed playout
func (db *DB) KickPlayout() error {
	fmt.Printf("[Playout] Start\n")
	t, err := db.GetOldestPlayoutTask()
	if err != nil {
		return errors.Wrapf(err, "Kick Playout")
	}
	if t == nil {
		fmt.Printf("[Playout] No Task\n")
		return nil
	}
	err = t.SpownPlayout(db.Config)
	if err != nil {
		err2 := t.PlayoutID.MakeFailed()
		if err != nil {
			return errors.Wrapf(err, "Make failed playout with %v", err2)
		}
		return errors.Wrapf(err, "Spown Playout")
	}
	fmt.Printf("[Playout] Complated")
	return nil
}

//GetPlayoutIDAndCheckToken returns Playout by ID with check TOKEN
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

//HandlerResultsUpdate handles the playout updating
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

//HandlerResultComplete makes the playout completed
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

//HandlerAddMatch appends the new playout
func HandlerAddMatch(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		s, err := db.GetSession(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Session Error %s", err))
		}
		if !s.IsUser() {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Session Error. your not user"))
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
		id, err := db.CreatePlayout(GameID(req.GameID), aiid, false)
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

//TokenLen is length of tokens to update playout
const TokenLen = 32

//PlayoutID has DB Config
type PlayoutID struct {
	ID int64
	DB *DB
}

//PlayerCmd is AI config
type PlayerCmd struct {
	ID     AIID
	Github string
	Branch string
	Commit string
}

//PlayoutTask is Playout Config
type PlayoutTask struct {
	PlayoutID PlayoutID
	Token     string
	Game      string
	Players   []PlayerCmd
}

//GetOldestPlayoutTask fetchs oldest playout
func (db *DB) GetOldestPlayoutTask() (*PlayoutTask, error) {
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
		ORDER BY playout_ai.turn
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

//Run runs the playout
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

//SpownPlayout spowns the playout
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

//ValidateToken validates the token of playout updating
func (p *PlayoutID) ValidateToken(token string) (bool, error) {
	tok := []sql.NullString{}
	err := p.DB.DB.Select(&tok, `
		SELECT playout.token AS token
		FROM playout
		WHERE playout.id = ?
	`, p.ID)
	if err != nil {
		return false, errors.Wrapf(err, "Failed Select")
	}
	if len(tok) < 1 {
		return false, errors.New("Missing Playout")
	}
	return tok[0].String == token, nil
}

//Update append the result to the playout.
func (p *PlayoutID) Update(result protocol.ResultA) error {
	_, err := p.DB.DB.Exec(`
		INSERT INTO playout_result (playout_id, record, exception)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE record=?, exception=?
	`, p.ID, result.Record, result.Exception, result.Record, result.Exception)
	if err != nil {
		return errors.Wrapf(err, "Failed insert playout_result")
	}
	return nil
}

//MakeFailed makes the playout failed
func (r *PlayoutID) MakeFailed() error {
	_, err := r.DB.DB.Exec(`
		UPDATE playout SET state = "failed"
		WHERE id = ?
	`, r.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed insert playout_result")
	}
	return nil
}

//Complete appends final result to the playout  and makes the playout finished
func (p *PlayoutID) Complete(results []protocol.ResultPlayerA) error {
	//TODO: Using transactions
	rate, err := p.CalclateRate(results)
	if err != nil {
		return errors.Wrapf(err, "Calclate Rate")
	}
	for i, result := range results {
		_, err := p.DB.DB.Exec(`
			INSERT INTO playout_result_ai (turn, playout_id, result, stderr, exception, rate)
			VALUES (?, ?, ?, ?, ?, ?)
		`, i, p.ID, result.Result, result.Stderr, result.Exception, rate[i])
		if err != nil {
			return errors.Wrapf(err, "insert playout_result_ai")
		}
	}
	_, err = p.DB.DB.Exec(`
		UPDATE playout SET state="completed" WHERE playout.id=?
	`, p.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed update playout to completed")
	}
	return nil
}

func (p *PlayoutID) CalclateRate(results []protocol.ResultPlayerA) ([]float64, error) {
	rated, err := p.FetchRated()
	if err != nil {
		return nil, errors.Wrapf(err, "Fetch Rated")
	}
	rate, selfMatch, err := p.FetchLatestRate()
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}
	if selfMatch {
		rated = false
	}
	if !rated {
		return rate, nil
	}
	score := []int{}
	for _, result := range results {
		score = append(score, result.Result)
	}
	return eroRating.Rating(rate, score), nil
}

//CreatePlayout creates new playout
func (db *DB) CreatePlayout(gameID GameID, ais []AIID, rated bool) (int64, error) {
	// TODO: not rated playout match
	tx, err := db.DB.Beginx()
	if err != nil {
		return -1, errors.Wrapf(err, "Beginx")
	}
	token, err := GenerateRandomString(TokenLen)
	if err != nil {
		return -1, errors.Wrapf(err, "Generate Random")
	}
	res, err := tx.Exec(`
		INSERT INTO playout (game_id, state, token, rated)
		VALUES (?, "ready", ?, ?);
	`, gameID, token, rated)
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
