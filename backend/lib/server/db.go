package server

import (
	"time"

	"github.com/dghubble/sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB          *sqlx.DB
	CookieStore *sessions.CookieStore
	Config      *Config
}

func (c *Config) NewDB() *DB {
	db, err := sqlx.Open("mysql", c.DBName)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	cookieStore := sessions.NewCookieStore([]byte(c.Session.SessionSecret), nil)
	return &DB{db, cookieStore, c}
}

type UserID int64
type UserAuthority string

const (
	UserAuthorityVistor UserAuthority = "visitor"
	UserAuthorityUser   UserAuthority = "user"
	UserAuthorityAdmin  UserAuthority = "admin"
)

type UserM struct {
	ID         int64         `db:"id"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
	Name       string        `db:"name"`
	GithubName string        `db:"github_name"`
	Authority  UserAuthority `db:"authority"`
}

type AIGithubUpdating string

const (
	AIGithubActive AIGithubUpdating = "active"
	AIGithubIgnore AIGithubUpdating = "ignore"
)

type AIGithubID int64
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
	AIFailed AIState = "failed"
	AIPurged AIState = "purged"
)

type AIID int64
type AIM struct {
	ID         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	AIGithubID int64     `db:"ai_github_id"`
	Commit     string    `db:"commit"`
	State      AIState   `db:"state"`
}

type GameID int64
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
	PlayoutFailed    PlayoutState = "failed"
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
