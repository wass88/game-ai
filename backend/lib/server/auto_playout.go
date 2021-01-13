package server

import (
	"fmt"

	"github.com/pkg/errors"
)

type AutoPlayout struct {
	DB *DB
}


var MaxConcurrentPlayout = 1
func NewAutoPlayout(db *DB) *AutoPlayout {
	return &AutoPlayout {
		DB: db,
	}
}

type PlayoutConfig struct {
	GameID GameID
	AIIDs []AIID
}

func (a *AutoPlayout) RandomPlayoutConfig() (*PlayoutConfig, error) {
	// TODO: Only two players and GameID=2
	// TODO: Fix Order BY ID
	gameID:= 2
	res := []AIID{}
	err := a.DB.DB.Select(&res, `
		SELECT ai.id
		FROM ai_github
		INNER JOIN LATERAL (
			SELECT ai.state, ai.id, ai.ai_github_id
			FROM ai
			WHERE ai.ai_github_id = ai_github.id
			ORDER BY created_at DESC, ai.id DESC
			LIMIT 1
		) AS ai ON ai.ai_github_id = ai_github.id
		WHERE ai.state = "ready" && ai_github.game_id = ?
		ORDER BY RAND()
		LIMIT 2`, gameID)
	if err != nil {
		return nil, errors.Wrapf(err, "Select Random")
	}
	if len(res) < 2 {
		return nil, errors.Errorf("Missing enough ai config")
	}
	return &PlayoutConfig{GameID(gameID), res}, nil
}

func (p *AutoPlayout) CreatePlayoutFromConfig(c *PlayoutConfig) (error) {
	_, err := p.DB.CreatePlayout(c.GameID, c.AIIDs, true)
	if err != nil {
		return errors.Wrapf(err, "Create Playout")
	}
	return nil
}

func (p *AutoPlayout) Kick() (error) {
	fmt.Printf("[Auto Playout] Start\n")
	n, err := p.DB.GetOldestPlayoutTask()
	if err != nil {
		return errors.Wrapf(err, "Count Wating")
	}
	if n != nil {
		fmt.Printf("[Auto Playout] Already enqueued = %v\n", n)
		return nil
	}
	c, err := p.RandomPlayoutConfig()
	if err != nil {
		return errors.Wrapf(err, "Random Playout")
	}
	err = p.CreatePlayoutFromConfig(c)
	if err != nil {
		return errors.Wrapf(err, "Create Playout")
	}
	fmt.Printf("[Auto Playout] Completed \n")
	return nil

}