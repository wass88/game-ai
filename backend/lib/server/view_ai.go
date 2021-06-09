package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type AIRG struct {
	ID         AIID       `json:"id"`
	AIGIthubID AIGithubID `json:"ai_github_id"`
	Commit     string     `json:"commit"`
	State      AIState    `json:"state"`
	UpdatedAt  time.Time `json:"updated_at"`
	Rate       *float64 `json:"rate"`
}

type AIGithubRAI struct {
	ID       AIGithubID       `json:"id"`
	Github   string           `json:"github"`
	Branch   string           `json:"branch"`
	Updating AIGithubUpdating `json:"updating"`
	User     UserR            `json:"user"`
	Game     GameR            `json:"game"`
	LatestAI *AIRG            `json:"latest_ai"`
}

func HandlerViewAIGithubByGame(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		idp := c.Param("id")
		id, err := strconv.Atoi(idp)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}
		res, err := db.GetAIGithubsByGame((GameID)(id))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed db select")
		}
		return c.JSON(http.StatusOK, res)
	}
}

func HandlerViewLatestByGame(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		idp := c.Param("id")
		id, err := strconv.Atoi(idp)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}
		res, err := db.GetLatestAIByGame((GameID)(id))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed db select")
		}
		return c.JSON(http.StatusOK, res)
	}
}

func (db *DB) GetAIGithubsByGame(id GameID) ([]AIGithubRAI, error) {
	type Result struct {
		ID       AIGithubID       `db:"id"`
		GameID   GameID           `db:"game_id"`
		GameName string           `db:"game_name"`
		UserID   UserID           `db:"user_id"`
		UserName string           `db:"user_name"`
		Updating AIGithubUpdating `db:"updating"`
		Github   string           `db:"github"`
		Branch   string           `db:"branch"`
		AIID     *AIID            `db:"ai_id"`
		AIStatus *AIState         `db:"ai_status"`
		AICommit *string          `db:"ai_commit"`
		AIUpdatedAt *time.Time    `db:"ai_updated_at"`
		AIRate   *float64         `db:"ai_rate"`
	}
	var sel []Result
	err := db.DB.Select(&sel, `
	WITH playout_rate AS (
		SELECT o_playout_ai.ai_id, o_result_ai.rate, o_playout.game_id, o_result_ai.created_at, o_result_ai.id
		FROM playout AS o_playout
		INNER JOIN playout_ai AS o_playout_ai ON o_playout_ai.playout_id = o_playout.id
		INNER JOIN playout_result_ai AS o_result_ai 
			ON o_result_ai.playout_id = o_playout_ai.playout_id
			AND o_result_ai.turn = o_playout_ai.turn
	),
	
	playout_rate_latest AS (
		SELECT o.ai_id, MAX(o.id) AS max_id
		FROM playout_rate AS o
		GROUP BY o.ai_id
	),
	
	rate AS (
		SELECT o.ai_id, o.rate, o.game_id
		FROM playout_rate AS o
		LEFT JOIN playout_rate_latest AS t ON t.ai_id = o.ai_id
		WHERE o.id = t.max_id
	)
	
	SELECT g.id AS id, g.game_id AS game_id, g.user_id AS user_id,
		gm.name AS game_name, u.name AS user_name,
		g.updating AS updating,
		g.github AS github, g.branch AS branch,
		ai.id AS ai_id, ai.state AS ai_status, ai.commit AS ai_commit,
		ai.updated_at AS ai_updated_at, rate.rate AS ai_rate
	FROM ai_github AS g
	INNER JOIN game AS gm ON gm.id = g.game_id
	INNER JOIN user AS u ON u.id = g.user_id
	LEFT JOIN LATERAL (SELECT * FROM ai
		WHERE ai.ai_github_id = g.id
		ORDER BY ai.created_at DESC LIMIT 1) AS ai ON ai.ai_github_id = g.id
	LEFT JOIN rate ON rate.game_id = g.game_id AND rate.ai_id = ai.id
	WHERE g.game_id = ?
	ORDER BY rate.rate DESC;
	`, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed select game %d", id)
	}
	res := []AIGithubRAI{}
	for _, s := range sel {
		aig := AIGithubRAI{
			ID: s.ID,
			Game: GameR{
				ID:   s.GameID,
				Name: s.GameName,
			},
			User: UserR{
				ID:   s.UserID,
				Name: s.UserName,
			},
			Updating: s.Updating,
			Github:   s.Github,
			Branch:   s.Branch,
		}
		if s.AIID != nil {
			aig.LatestAI = &AIRG{
				ID:         *s.AIID,
				AIGIthubID: s.ID,
				State:      *s.AIStatus,
				Commit:     *s.AICommit,
				UpdatedAt:  *s.AIUpdatedAt,
				Rate:       s.AIRate,
			}
		}
		res = append(res, aig)
	}
	return res, nil
}

type AIShortInfo struct {
	ID     AIID   `db:"id" json:"id"`
	Github string `db:"github" json:"github"`
	Branch string `db:"branch" json:"branch"`
	Commit string `db:"commit" json:"commti"`
	Rate *float64 `db:"rate" join:"rate"`
}

func (db *DB) GetLatestAIByGame(id GameID) ([]AIShortInfo, error) {
	var res []AIShortInfo
	err := db.DB.Select(&res, `
		SELECT ai.id AS id, ai.commit AS commit, g.github AS github, g.branch AS branch, rate.rate
		FROM ai
		LEFT JOIN ai AS b ON (ai.created_at < b.created_at AND ai.ai_github_id = b.ai_github_id)
		INNER JOIN ai_github AS g ON g.id = ai.ai_github_id
		LEFT JOIN (
		SELECT o_playout_ai.ai_id, o_result_ai.rate, o_playout.game_id
		FROM playout AS o_playout
		INNER JOIN playout_ai AS o_playout_ai ON o_playout_ai.playout_id = o_playout.id
		INNER JOIN playout_result_ai AS o_result_ai ON o_result_ai.turn = o_playout_ai.turn AND o_result_ai.playout_id = o_playout_ai.playout_id
		WHERE NOT EXISTS(
		SELECT 1 FROM playout_ai AS t_playout_ai
		INNER JOIN playout_result_ai AS t_result_ai ON t_result_ai.turn = t_playout_ai.turn AND t_result_ai.playout_id = t_playout_ai.playout_id
		INNER JOIN playout AS t_playout ON t_playout.id = t_playout_ai.playout_id
		WHERE o_playout.game_id = t_playout.game_id
			AND o_playout_ai.ai_id = t_playout_ai.ai_id
			AND o_result_ai.created_at <= t_result_ai.created_at
			AND o_result_ai.id < t_result_ai.id
		)
		) AS rate ON rate.game_id = g.game_id AND rate.ai_id = ai.id
		WHERE ai.state = "ready"
		AND g.game_id = ?
		AND b.created_at IS NULL
	`, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed select game %d", id)
	}
	return res, nil
}
