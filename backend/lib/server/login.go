package server

import (
	"fmt"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

const (
	sessionName   = "game-ai-session"
	sessionGithub = "user_github"
)

func (db *DB) SetSessionHandler(e *echo.Echo) {
	e.POST("/github/logout", logoutHandler(db))
	oauth2Config := &oauth2.Config{
		ClientID:     db.Config.Session.GithubClientID,
		ClientSecret: db.Config.Session.GithubClientSecret,
		RedirectURL:  db.Config.Session.GithubRedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
	}
	// TODO: https only
	stateConfig := gologin.DebugOnlyCookieConfig
	errorHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Error %v", req)
	}
	e.GET("/github/login", echo.WrapHandler(github.StateHandler(
		stateConfig,
		github.LoginHandler(oauth2Config, http.HandlerFunc(errorHandler)),
	)))
	e.GET("/github/callback", echo.WrapHandler(github.StateHandler(
		stateConfig,
		github.CallbackHandler(oauth2Config, db.issueSession(), http.HandlerFunc(errorHandler)),
	)))
	e.GET("/github/logout", logoutHandler(db))
}

// issueSession issues a cookie session after successful Github login
func (db *DB) issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if githubUser == nil {
			http.Error(w, "No githubUser", http.StatusInternalServerError)
		}
		_, err = db.NewUserIfNotExist(*githubUser.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session := db.CookieStore.New(sessionName)
		session.Values[sessionGithub] = *githubUser.Login
		err = session.Save(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

type Session struct {
	ID         UserID
	GithubName string
	Name       string
	Authority  UserAuthority
}

func (s *Session) IsUser() bool {
	return s.Authority == UserAuthorityUser || s.Authority == UserAuthorityAdmin
}

func (db *DB) GetSession(c echo.Context) (*Session, error) {
	s, err := db.CookieStore.Get(c.Request(), sessionName)
	if err != nil {
		return nil, err
	}
	github, ok := s.Values[sessionGithub].(string)
	if !ok {
		return nil, errors.Wrapf(err, "invalid github id %v", s.Values[sessionGithub])
	}
	user, ok, err := db.GetUserByGithub(github)
	if err != nil {
		return nil, errors.Wrapf(err, "Faild Select User")
	}
	if !ok {
		return nil, errors.Errorf("Missing User %v", github)
	}
	return &Session{
		ID:         UserID(user.ID),
		GithubName: github,
		Name:       user.Name,
		Authority:  user.Authority,
	}, nil
}

func logoutHandler(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		db.CookieStore.Destroy(c.Response().Writer, sessionName)
		return c.Redirect(http.StatusFound, "/")
	}
}

type You struct {
	Login      bool    `json:"login"`
	ID         *int64  `json:"id"`
	GithubName *string `json:"github_name"`
	Name       *string `json:"name"`
	Authority  *string `json:"authority"`
	Err        *string `json:"err"`
}

func HandlerYou(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		s, err := db.GetSession(c)
		if err != nil {
			fmt.Printf("Session err %v\n", err)
			e := err.Error()
			return c.JSON(http.StatusOK, You{Err: &e})
		}
		return c.JSON(http.StatusOK, You{true, (*int64)(&s.ID), &s.GithubName, &s.Name, (*string)(&s.Authority), nil})
	}
}

func (db *DB) GetUserByGithub(github string) (*UserM, bool, error) {
	var res []*UserM
	err := db.DB.Select(&res, `
		SELECT * FROM user WHERE github_name = ? LIMIT 1
	`, github)
	if err != nil {
		return nil, false, errors.Wrapf(err, "Failed Select")
	}
	if len(res) == 0 {
		return nil, false, nil
	}
	return res[0], true, nil
}
func (db *DB) NewUserIfNotExist(githubName string) (*UserM, error) {
	user, ok, err := db.GetUserByGithub(githubName)
	if err != nil {
		return nil, errors.Wrapf(err, "Get User")
	}
	if ok {
		return user, nil
	}
	_, err = db.NewUser(githubName)
	if err != nil {
		return nil, err
	}
	user, ok, err = db.GetUserByGithub(githubName)
	if err != nil {
		return nil, errors.Wrapf(err, "Get User")
	}
	if ok {
		return nil, errors.Wrapf(err, "Missing user after new")
	}
	return user, nil
}

func (db *DB) NewUser(githubName string) (UserID, error) {
	res, err := db.DB.Exec(`
		INSERT INTO user (github_name, name, authority)
		VALUES (?, ?, "visitor")
	`, githubName, githubName)
	if err != nil {
		return -1, errors.Wrapf(err, "Failed Insert")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return UserID(id), nil
}
