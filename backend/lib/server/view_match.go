package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type GameR struct {
	ID   GameID `json:"id"`
	Name string `json:"name"`
}

type UserR struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}

type AIR struct {
	ID       AIID      `json:"id"`
	Commit   string    `json:"commit"`
	AIGithub AIGithubR `json:"ai_github"`
}

type AIGithubR struct {
	ID       AIGithubID       `json:"id"`
	Github   string           `json:"github"`
	Branch   string           `json:"branch"`
	Updating AIGithubUpdating `json:"updating"`
	User     UserR            `json:"user"`
	Game     GameR            `json:"game"`
}

type MatchID int64
type MatchShortR struct {
	ID        MatchID        `json:"id"`
	Game      GameR          `json:"game"`
	State     PlayoutState   `json:"state"`
	Exception string         `json:"exception"`
	Results   []ResultShortR `json:"results"`
}

type ResultShortR struct {
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

		pageP := c.QueryParam("page")
		var page *int = nil
		if pageP != "" {
			val, err := strconv.Atoi(pageP)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid pages")
			}
			page = &val
		}

		//TODO: aiP := c.QueryParam("ais")

		res, err := db.GetMatches((GameID)(id), page)
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

type PagedMatchShortR struct {
	Pages int `json:"pages"`
	Matches []*MatchShortR `json:"matches"`
}
func (db *DB) GetMatches(id GameID, page *int) (*PagedMatchShortR, error) {
	type Result struct {
		Count int `json:"count"`

		ID        MatchID      `db:"id"`
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

	pageSize := 100
	offset := 0
	if page != nil {
		offset = pageSize * *page
	}
	var sel []Result
	err := db.DB.Select(&sel, `
	SELECT
		COUNT(1) OVER() AS count,
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
	LIMIT ? OFFSET ?
	`, id, pageSize, offset)

	if err != nil {
		return nil, errors.Wrapf(err, "Faild select: Game %d", id)
	}

	pages := 1
	if len(sel) > 0 {
		pages = (sel[0].Count + pageSize - 1) / pageSize
	}

	res := []*MatchShortR{}
	var c *MatchShortR
	for _, s := range sel {
		if c != nil && c.ID != s.ID {
			res = append(res, c)
			c = nil
		}
		if c == nil {
			c = &MatchShortR{
				ID:    s.ID,
				State: s.State,
				Game: GameR{
					ID:   s.GameID,
					Name: s.GameName,
				},
				Exception: s.Exception,
				Results:   []ResultShortR{},
			}
		}
		r := ResultShortR{
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
	return &PagedMatchShortR{Pages: pages, Matches: res}, nil
}

type MatchR struct {
	ID        MatchID      `json:"id"`
	Game      GameR        `json:"game"`
	State     PlayoutState `json:"state"`
	Exception string       `json:"exception"`
	Results   []ResultR    `json:"results"`
	Record    string       `json:"record"`
}

type ResultR struct {
	AI        AIR     `json:"ai"`
	Result    *int    `json:"result"`
	Exception *string `json:"exeption"`
	Stderr    *string `json:"stderr"`
}

func HandlerViewMatch(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		idP := c.Param("id")
		id, err := strconv.Atoi(idP)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}
		res, err := db.GetMatch((MatchID)(id))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Failed db select")
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (db *DB) GetMatch(id MatchID) (*MatchR, error) {
	type Result struct {
		ID        MatchID      `db:"id"`
		State     PlayoutState `db:"state"`
		Exception string       `db:"exception"`
		Record    string       `db:"record"`

		GameID   GameID `db:"game_id"`
		GameName string `db:"game_name"`

		Result     **int      `db:"result"`
		RException *string    `db:"rexception"`
		RStderr    *string    `db:"rstderr"`
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
			p.id AS id, p.state AS state,
			IFNULL(r.exception, "") AS exception,
			IFNULL(r.record, "") AS record,
			p.game_id AS game_id, g.name AS game_name,
			rai.result AS result, rai.exception AS rexception,
			rai.stderr AS rstderr,
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
		WHERE p.id = ?
		ORDER BY p.id DESC, pai.turn ASC
	`, id)

	if err != nil {
		return nil, errors.Wrapf(err, "Faild select: Game %d", id)
	}

	var res *MatchR
	for _, s := range sel {
		if res == nil {
			res = &MatchR{
				ID:    s.ID,
				State: s.State,
				Game: GameR{
					ID:   s.GameID,
					Name: s.GameName,
				},
				Exception: s.Exception,
				Record:    s.Record,
				Results:   []ResultR{},
			}
		}
		r := ResultR{
			Result:    dInt(s.Result),
			Exception: s.RException,
			Stderr:    s.RStderr,
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
		res.Results = append(res.Results, r)
	}
	if res == nil {
		return res, errors.Wrapf(err, "Missing Match %d", id)
	}
	return res, nil
}
