package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

const (
	sessionName     = "game-ai-session"
	sessionUserKey  = "githubID"
	sessionUsername = "githubUsername"
	sessionUserID   = "githubUserID"
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
		fmt.Printf("Issue Session\n")
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("HEY %v\n", githubUser)
		if githubUser == nil {
			http.Error(w, "No githubUser", http.StatusInternalServerError)
		}
		id, err := db.NewUserIfNotExist(*githubUser.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session := db.CookieStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Values[sessionUsername] = *githubUser.Login
		session.Values[sessionUserID] = strconv.Itoa((int)(id))
		err = session.Save(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("set session %v\n", session)
		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

type Session struct {
	ID         UserID
	GithubName string
}

func (db *DB) GetSession(c echo.Context) (*Session, error) {
	s, err := db.CookieStore.Get(c.Request(), sessionName)
	if err != nil {
		return nil, err
	}
	ids := s.Values[sessionUserID]
	id, err := strconv.Atoi(ids.(string))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid ids %s", ids)
	}
	return &Session{
		ID:         UserID(id),
		GithubName: s.Values[sessionUsername].(string),
	}, nil
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		db.CookieStore.Destroy(c.Response().Writer, sessionName)
		return c.Redirect(http.StatusFound, "/")
	}
}

type You struct {
	Login bool    `json:"login"`
	Name  *string `json:"name"`
	ID    *int64  `json:"id"`
	Err   *string `json:"err"`
}

func HandlerYou(db *DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		s, err := db.GetSession(c)
		if err != nil {
			fmt.Printf("Session err %v\n", err)
			e := err.Error()
			return c.JSON(http.StatusOK, You{false, nil, nil, &e})
		}
		fmt.Printf("got session = %v\n", s)
		if s == nil {
			e := "Session Error"
			return c.JSON(http.StatusOK, You{false, nil, nil, &e})
		}

		user := s.GithubName
		id := (int64)(s.ID)
		return c.JSON(http.StatusOK, You{true, &user, &id, nil})
	}
}

func (db *DB) NewUserIfNotExist(githubName string) (UserID, error) {
	var res []struct {
		UserID *UserID `db:"id"`
	}
	err := db.DB.Select(&res, `
		SELECT id FROM user WHERE github_name = ? LIMIT 1
	`, githubName)
	if err != nil {
		return -1, errors.Wrapf(err, "Failed Select")
	}
	if len(res) == 0 {
		id, err := db.NewUser(githubName)
		if err != nil {
			return -1, err
		}
		return id, nil
	}
	return *res[0].UserID, nil
}

func (db *DB) NewUser(githubName string) (UserID, error) {
	res, err := db.DB.Exec(`
		INSERT INTO user (github_name, name)
		VALUES (?, ?)
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
