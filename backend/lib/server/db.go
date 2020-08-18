package server

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB     *sqlx.DB
	Config *Config
}

type Config struct {
	AIRunner AIRunnerConf
}

type AIRunnerConf struct {
	Cmd string
	API string
	Dir string
}

func NewDB(dbname string) *DB {
	db, err := sqlx.Open("mysql", dbname)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &DB{db, nil}
}

type UserM struct {
	ID           int64     `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Name         string    `db:"name"`
	TwitterToken string    `db:"twitter_token"`
}

type AIGithubUpdating string

const (
	AIGithubActive AIGithubUpdating = "active"
	AIGithubIgnore AIGithubUpdating = "ignore"
)

type AIGithubM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserID    int64     `db:"user_id"`
	GameID    int64     `db:"game_id"`
	Github    string    `db:"github"`
	Branch    string    `db:"branch"`
	Updating  string    `db:"updating"`
}

type AIState string

const (
	AIFound  AIState = "found"
	AISetup  AIState = "setup"
	AIReady  AIState = "ready"
	AIPurged AIState = "purged"
)

type AIM struct {
	ID         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	AIGithubID int64     `db:"ai_github_id"`
	Commit     string    `db:"commit"`
	State      AIState   `db:"state"`
}

type GameM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
}

type PlayoutState string

var (
	PlayoutReady     PlayoutState = "ready"
	PlayoutRunning   PlayoutState = "running"
	PlayoutCompleted PlayoutState = "completed"
)

type PlayoutM struct {
	ID        int64        `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	State     PlayoutState `db:"state"`
	GameID    string       `db:"game_id"`
	Token     string       `db:"token"`
}

type PlayoutAIM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	AIID      int64     `db:"ai_id"`
	Turn      int64     `db:"turn"`
}

type PlayoutResultM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	PlayoutID int64     `db:"playout_id"`
	Records   string    `db:"record"`
	Exception string    `db:"exception"`
}

type PlayoutResultAIM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Result    string    `db:"result"`
	Stderr    string    `db:"stderr"`
	Exception string    `db:"exception"`
}
