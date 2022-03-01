package sessionmgm

import (
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginForm struct {
	CSRFField  template.HTML
	UserName   string
	Email      string
	Password   string
	FormErrors map[string]string
}

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var sessions = map[string]session{}

// each session contains the email of the user and the time at which it expires
type session struct {
	email    string
	username string
	isAdmin  string
	expiry   time.Time
}
type sesMethod interface {
	isExpired()
	saveSessionCookie()
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
func (s session) SaveSessionCookie(w http.ResponseWriter, form LoginForm) {
	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(120 * time.Second)
	sessions[sessionToken] = session{
		email:    form.Email,
		username: form.UserName,
		isAdmin:  "",
		expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
}
