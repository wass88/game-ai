package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (db *DB) KickSetupAI() error {
	fmt.Printf("[Setup AI]Start\n")
	ain, err := db.FindAIGithubNeedUpdate()
	if err != nil {
		return errors.Wrapf(err, "[Setup AI] Faild FindAIGithubNeedUpdate")
	}
	for _, ai := range ain {
		fmt.Printf("[Setup AI] Regiter AI %v\n", ai)
		db.CreateAI(&ai)
	}

	ais, err := db.GetNeedSetupAI()
	if err != nil {
		return errors.Wrapf(err, "[Setup AI] Failed GetNeedSetupAI")
	}
	if len(ais) == 0 {
		fmt.Printf("[Setup AI] No AI to need to update\n")
	}
	for _, ai := range ais {
		fmt.Printf("[Setup AI] Kick Create AI %v\n", ai)
		err := ai.KickSetup(db)
		if err != nil {
			return errors.Wrapf(err, "Failed KickSetup")
		}
		fmt.Printf("[Setup AI] Done Create AI %v\n", ai)
		break
	}
	return nil
}

func HandlerReadyContainer(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		var req struct {
			Github string `json:"github"`
			Branch string `json:"branch"`
			Commit string `json:"commit"`
		}
		err := c.Bind(&req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid json %v", err))
		}
		c.Logger().Debugf("%v", req)
		err = c.Validate(&req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed validation json %v", err))
		}

		err = db.ReadyContianersByCommit(req.Github, req.Branch, req.Commit)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed db %v", err))
		}
		return nil
	}
}

func HandlerAddAIGithub(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		s, err := db.GetSession(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Session Error %s", err))
		}
		if !s.IsUser() {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Session Error. your not user"))
		}
		id := s.ID
		var req struct {
			GameID int64  `json:"game_id" validate:"required""`
			Github string `json:"github" validate:"required"`
			Branch string `json:"branch" validate:"required"`
		}
		err = c.Bind(&req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid json %v", err))
		}
		err = c.Validate(&req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed validation json %v", err))
		}
		github := AIGithubA{
			UserID: id,
			GameID: req.GameID,
			Github: req.Github,
			Branch: req.Branch,
		}
		res, err := db.CreateAIGithub(&github)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed db create %v", err))
		}

		type Result struct {
			AIGithubID int64 `json:"ai_github_id"`
		}

		return c.JSON(http.StatusOK, Result{(int64)(res)})
	}
}

type AIGithubA struct {
	GameID int64
	UserID UserID
	Github string
	Branch string
}

func (db *DB) CreateAIGithub(ai *AIGithubA) (AIGithubID, error) {
	res, err := db.DB.Exec(`
		INSERT INTO ai_github (game_id, user_id, github, branch, updating)
		VALUES (?, ?, ?, ?, ?)
	`, ai.GameID, ai.UserID, ai.Github, ai.Branch, "active")
	if err != nil {
		return -1, errors.Wrapf(err, "Create AIGithub")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, errors.Wrapf(err, "LastInsertId")
	}
	return AIGithubID(id), nil
}

type AIGithubAct struct {
	Github string `db:"github"`
	Branch string `db:"branch"`
	// IF NULL, ""
	LastAICommit string     `db:"last_commit"`
	ID           AIGithubID `db:"id"`
}

func (db *DB) GetActiveAI() ([]AIGithubAct, error) {
	res := []AIGithubAct{}
	err := db.DB.Select(&res, `
		SELECT github, branch, id,
			ifnull((SELECT commit FROM ai
				WHERE ai_github_id = ai_github.id
				ORDER BY created_at DESC
				LIMIT 1),"") AS last_commit
		FROM ai_github
		WHERE updating = "active"
	`)
	if err != nil {
		return nil, errors.Wrapf(err, "Select github_ai")
	}
	return res, nil
}

type AIGithubNeedUpdate struct {
	Github string
	Branch string
	Commit string
	ID     AIGithubID
}

func (db *DB) FindAIGithubNeedUpdate() ([]AIGithubNeedUpdate, error) {
	ai, err := db.GetActiveAI()
	if err != nil {
		return nil, err
	}
	res := []AIGithubNeedUpdate{}
	for _, a := range ai {
		commit, err := a.GetLastCommit()
		if err != nil {
			fmt.Printf("Failed Get LastCommit %v\n", err)
			continue
		}
		fmt.Printf("Current: %#v\n%s\n%s\n", a, a.LastAICommit, commit)
		if a.LastAICommit != commit {
			n := AIGithubNeedUpdate{
				Github: a.Github,
				Branch: a.Branch,
				Commit: commit,
				ID:     a.ID,
			}
			res = append(res, n)
		}
	}
	return res, nil
}

func FetchCommitFromGithub(userRepo, branch string) (string, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/branches/%s", userRepo, branch)

	resp, err := http.Get(u)
	if err != nil {
		return "", errors.Wrapf(err, "Failed Get %s", u)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Fail ReadAll")
	}
	if resp.StatusCode != 200 {
		return "", errors.Errorf("Bad Status %d from %s\nBody: %s", resp.StatusCode, u, bytes)
	}
	var res struct {
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return "", errors.Wrapf(err, "Faild Unmarshal")
	}
	return res.Commit.SHA, nil
}

func (a AIGithubAct) GetLastCommit() (string, error) {
	commit, err := FetchCommitFromGithub(a.Github, a.Branch)
	if err != nil {
		return "", errors.Wrapf(err, "Fetch Commit")
	}
	return commit, nil
}

func (db *DB) CreateAI(a *AIGithubNeedUpdate) (AIID, error) {
	res, err := db.DB.Exec(`
		INSERT INTO ai (state, ai_github_id, commit)
		VALUE (?, ?, ?)
	`, "found", a.ID, a.Commit)
	if err != nil {
		return -1, errors.Wrapf(err, "Insert %v", a)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, errors.Wrapf(err, "Fail LastInsertId")
	}
	return AIID(id), nil
}

type AINeedSetup struct {
	ID     AIID   `db:"id"`
	Github string `db:"github"`
	Branch string `db:"branch"`
	Commit string `db:"commit"`
}

func (db *DB) GetNeedSetupAI() ([]AINeedSetup, error) {
	res := []AINeedSetup{}
	err := db.DB.Select(&res, `
		SELECT ai.id AS id, g.github AS github, g.branch AS branch, ai.commit AS commit
		FROM ai INNER JOIN ai_github AS g ON ai.ai_github_id = g.id
		WHERE ai.state = "found"
	`)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed Select")
	}
	return res, nil
}

func (ai *AINeedSetup) SetupCmd(c AIRunnerConf) *exec.Cmd {
	args := []string{"-api", c.API, "-dir", c.Dir,
		"-github", ai.Github, "-branch", ai.Branch,
		"-commit", ai.Commit, "setup",
	}
	return exec.Command(c.Cmd, args...)
}

func (a *AIID) UpdateState(db *DB, state AIState) error {
	_, err := db.DB.Exec(`
		UPDATE ai
		SET state = ?
		WHERE id = ?
	`, state, a)
	if err != nil {
		return errors.Wrapf(err, "Update State")
	}
	return nil
}

func (db *DB) ReadyContianersByCommit(github, branch, commit string) error {
	_, err := db.DB.Exec(`
		UPDATE ai
		INNER JOIN ai_github AS g ON g.id = ai.ai_github_id
		SET ai.state = "ready"
		WHERE commit = ? AND g.branch = ? AND g.github = ?
	`, commit, branch, github)
	if err != nil {
		return errors.Wrapf(err, "Update States")
	}
	return nil
}

func (ai *AINeedSetup) KickSetup(db *DB) error {
	cmd := ai.SetupCmd(db.Config.AIRunner)
	err := ai.ID.UpdateState(db, "setup")
	if err != nil {
		return errors.Wrapf(err, "Failed Update State")
	}
	bytes, err := cmd.CombinedOutput()
	log.Printf("Output: %s\n", bytes)
	if err != nil {
		ai.ID.UpdateState(db, "failed")
		return errors.Wrapf(err, "Failed to Run: err")
	}
	return nil
}
