package server

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB *sqlx.DB
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
	return &DB{db}
}

type UserM struct {
	ID           int64     `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Name         string    `db:"name"`
	TwitterToken string    `db:"twitter_token"`
}

type AIGithubM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserID    int64     `db:"user_id"`
	Github    string    `db:"github"`
	Branch    string    `db:"branch"`
}

type AIM struct {
	ID         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	AIGithubID int64     `db:"ai_github_id"`
	Commit     string    `db:"commit"`
	State      string    `db:"state"`
}

type GameM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
}

type PlayoutM struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	State     string    `db:"state"`
	GameID    string    `db:"game_id"`
	Token     string    `db:"token"`
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
