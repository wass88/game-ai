package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/pkg/errors"
)

func (db *DB) KickSetupAI() error {
	ain, err := db.FindAIGithubNeedUpdate()
	if err != nil {
		return errors.Wrapf(err, "Faild FindAIGithubNeedUpdate")
	}
	for _, ai := range ain {
		db.CreateAI(&ai)
	}

	ais, err := db.GetNeedSetupAI()
	if err != nil {
		return errors.Wrapf(err, "Failed GetNeedSetupAI")
	}
	for _, ai := range ais {
		err := ai.KickSetup(db)
		if err != nil {
			return errors.Wrapf(err, "Failed KickSetup")
		}
		break
	}
	return nil
}

type AIGithubA struct {
	GameID int64
	UserID int64
	Github string
	Branch string
}
type AIGithubID int64

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
				WHERE id = ai_github.id
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
			return nil, errors.Wrapf(err, "Failed Get LastCommit")
		}

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
	if resp.StatusCode != 200 {
		return "", errors.Errorf("Bad Status %d", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Fail ReadAll")
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

type AIID int64

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

func (ai *AINeedSetup) KickSetup(db *DB) error {
	cmd := ai.SetupCmd(db.Config.AIRunner)
	err := ai.ID.UpdateState(db, "setup")
	if err != nil {
		return errors.Wrapf(err, "Failed Update State")
	}
	err = cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "Failed cmd Start: %v", cmd)
	}
	return nil
}
