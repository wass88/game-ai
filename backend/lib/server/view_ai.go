package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type AIRG struct {
	ID         AIID       `json:"id"`
	AIGIthubID AIGithubID `json:"ai_github_id"`
	Commit     string     `json:"commit"`
	State      AIState    `json:"state"`
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
	}
	var sel []Result
	err := db.DB.Select(&sel, `
		SELECT g.id AS id, g.game_id AS game_id, g.user_id AS user_id,
			gm.name AS game_name, u.name AS user_name,
			g.updating AS updating,
			g.github AS github, g.branch AS branch,
			ai.id AS ai_id, ai.state AS ai_status, ai.commit AS ai_commit
		FROM ai_github AS g
		INNER JOIN game AS gm ON gm.id = g.game_id
		INNER JOIN user AS u ON u.id = g.user_id
		INNER JOIN LATERAL (SELECT * FROM ai
			WHERE ai.ai_github_id = g.id
			ORDER BY ai.created_at DESC LIMIT 1) AS ai ON ai.ai_github_id = g.id
		WHERE g.game_id = ?
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
			}
		}
		res = append(res, aig)
	}
	return res, nil
}
