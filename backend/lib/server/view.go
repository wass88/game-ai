package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type GameID int64

type GameR struct {
	ID   GameID `json:"id"`
	Name string `json:"name"`
}

type UserID int64
type UserR struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}

type AIID int64
type AIR struct {
	ID       AIID      `json:"id"`
	Commit   string    `json:"commit"`
	AIGithub AIGithubR `json:"ai_github"`
}

type AIGithubID int64
type AIGithubR struct {
	ID     AIGithubID `json:"id"`
	Github string     `json:"github"`
	Branch string     `json:"branch"`
	User   UserR      `json:"user"`
	Game   GameR      `json:"game"`
}

type ViewMatcheID int64
type ViewMatcheShortR struct {
	ID        ViewMatcheID       `json:"id"`
	Game      GameR              `json:"game"`
	State     PlayoutState       `json:"state"`
	Exception string             `json:"exception"`
	Results   []ViewResultShortR `json:"results"`
}

type ViewResultShortR struct {
	AI        AIR     `json:"ai"`
	Result    *int    `json:"result"`
	Exception *string `json:"exeption"`
}

func HandlerViewMatches(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		idP := c.Param("id")
		id, err := strconv.Atoi(idP)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}
		res, err := db.GetMatches((GameID)(id))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Failed db select")
		}

		return c.JSON(http.StatusOK, res)
	}
}
func dInt(i **int) *int {
	if i == nil {
		return nil
	}
	return *i
}
func (db *DB) GetMatches(id GameID) ([]*ViewMatcheShortR, error) {
	type Result struct {
		ID        ViewMatcheID `db:"id"`
		State     PlayoutState `db:"state"`
		Exception string       `db:"exception"`

		GameID   GameID `db:"game_id"`
		GameName string `db:"game_name"`

		Result     **int      `db:"result"`
		RException *string    `db:"rexception"`
		AITurn     int        `db:"ai_turn"`
		AIID       AIID       `db:"ai_id"`
		AICommit   string     `db:"ai_commit"`
		AIGithubID AIGithubID `db:"aigithub_id"`
		AIBranch   string     `db:"ai_branch"`
		AIGithub   string     `db:"ai_github"`
		UserID     UserID     `db:"user_id"`
		UserName   string     `db:"user_name"`
	}
	var sel []Result
	err := db.DB.Select(&sel, `
		SELECT
			p.id AS id, p.state AS state, IFNULL(r.exception, "") AS exception,
			p.game_id AS game_id, g.name AS game_name,
			rai.result AS result, rai.exception AS rexception,
			ai.id AS ai_id, ai.commit AS ai_commit,
			ai.ai_github_id AS aigithub_id, ag.github AS ai_github, ag.branch AS ai_branch,
			u.id AS user_id, u.name AS user_name
		FROM playout AS p
		INNER JOIN game AS g ON g.id = p.game_id 
		INNER JOIN playout_ai AS pai ON pai.playout_id = p.id
		INNER JOIN ai AS ai ON ai.id = pai.ai_id
		INNER JOIN ai_github AS ag ON ag.id = ai.ai_github_id
		INNER JOIN user AS u ON u.id = ag.user_id
		LEFT JOIN playout_result AS r ON r.playout_id = p.id
		LEFT JOIN playout_result_ai AS rai ON rai.playout_id = p.id AND rai.turn = pai.turn
		WHERE g.id = ?
		ORDER BY p.id DESC, pai.turn ASC
	`, id)

	if err != nil {
		return nil, errors.Wrapf(err, "Faild select: Game %d", id)
	}

	res := []*ViewMatcheShortR{}
	var c *ViewMatcheShortR
	for _, s := range sel {
		fmt.Printf("%#+v\n", s)
		if c != nil && c.ID != s.ID {
			res = append(res, c)
			c = nil
		}
		if c == nil {
			c = &ViewMatcheShortR{
				ID:    s.ID,
				State: s.State,
				Game: GameR{
					ID:   s.GameID,
					Name: s.GameName,
				},
				Exception: s.Exception,
				Results:   []ViewResultShortR{},
			}
		}
		r := ViewResultShortR{
			Result:    dInt(s.Result),
			Exception: s.RException,
			AI: AIR{
				ID:     s.AIID,
				Commit: s.AICommit,
				AIGithub: AIGithubR{
					ID:     s.AIGithubID,
					Github: s.AIGithub,
					Branch: s.AIBranch,
					User: UserR{
						ID:   s.UserID,
						Name: s.UserName,
					},
				},
			},
		}
		c.Results = append(c.Results, r)
	}
	if c != nil {
		res = append(res, c)
	}
	return res, nil
}
