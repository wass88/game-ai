package server

import (
	"fmt"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

const (
	sessionName     = "game-ai-session"
	sessionUserKey  = "githubID"
	sessionUsername = "githubUsername"
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
		session := db.CookieStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Values[sessionUsername] = *githubUser.Login
		session.Save(w)
		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func (db *DB) GetSession(c echo.Context) (*sessions.Session, error) {
	return db.CookieStore.Get(c.Request(), sessionName)
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		db.CookieStore.Destroy(c.Response().Writer, sessionName)
		return c.HTML(http.StatusFound, "/")
	}
}

type You struct {
	Login bool    `json:"login"`
	Name  *string `json:"name"`
	Err   *string `json:"err"`
}

func HandlerYou(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		s, err := db.GetSession(c)
		if err != nil {
			e := err.Error()
			return c.JSON(http.StatusOK, You{false, nil, &e})
		}
		user := s.Values[sessionUsername].(string)
		return c.JSON(http.StatusOK, You{true, &user, nil})
	}
}
