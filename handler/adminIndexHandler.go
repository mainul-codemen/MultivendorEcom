package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	CSRFField  template.HTML
	UserName   string
	Email      string
	Password   string
	FormErrors map[string]string
}

type RegistrationForm struct {
	CSRFField  template.HTML
	ID         string
	UserName   string
	FirstName  string
	LastName   string
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

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s LoginForm) Validate(srv *Server) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email,
			validation.Required.Error("Email is Required"),
		),
		validation.Field(&s.Password,
			validation.Required.Error("Password is Required"),
			validation.By(checkLogin(srv, s.Email, s.Password)),
		),
	)
}

func (s RegistrationForm) Validate(srv *Server) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.UserName,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.By(checkDuplicateUserName(srv, s.UserName)),
		),
		validation.Field(&s.Email,
			validation.Required.Error("Email is Required"),
		),
		validation.Field(&s.Password,
			validation.Required.Error("Password is Required"),
			validation.By(validatePassword(srv, s.Password)),
		),
		validation.Field(&s.FirstName,
			validation.Required.Error("First Name is Required"),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("Last Name is Required"),
		),
	)
}

func (s *Server) adminindex(w http.ResponseWriter, r *http.Request) {
	logger.Info("index handler")
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		data := LoginForm{
			Email: userSession.email,
		}
		tmpl := s.lookupTemplate("index.html")
		if tmpl == nil {
			logger.Error(ult)
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
		if err := tmpl.Execute(w, data); err != nil {
			logger.Error(ewte + err.Error())
			http.Redirect(w, r, "/admin/index", http.StatusSeeOther)
		}
	}
}

func (s Server) refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code from this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route
	// If the previous session is valid, create a new session token for the current user
	newSessionToken := uuid.New().String()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the user whom it represents
	sessions[newSessionToken] = session{
		email:  userSession.email,
		expiry: expiresAt,
	}

	// Delete the older session token
	delete(sessions, sessionToken)

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	logger.Info("login form handler")
	formTemplate(s, w, r, "login.html")
}

// each session contains the email of the user and the time at which it expires

func (s *Server) submitLogin(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit login handler")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form LoginForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	if err := form.Validate(s); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := RegistrationForm{
			CSRFField:  csrf.TemplateField(r),
			UserName:   form.UserName,
			Email:      form.Email,
			Password:   form.Password,
			FormErrors: vErrs,
		}
		tmpl := s.lookupTemplate("login.html")
		if tmpl == nil {
			logger.Error(ult)
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
		if err := tmpl.Execute(w, data); err != nil {
			logger.Error(ewte + err.Error())
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		}
		return
	}
	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
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
	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds

	http.Redirect(w, r, "/admin/index", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

func (s *Server) registrationForm(w http.ResponseWriter, r *http.Request) {
	logger.Info("registration form")
	formTemplate(s, w, r, "registration-form.html")
}

func (s *Server) submitRegistration(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit user registration handler")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form RegistrationForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	if err := form.Validate(s); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := RegistrationForm{
			UserName:   form.UserName,
			Email:      form.Email,
			Password:   form.Password,
			FirstName:  form.FirstName,
			LastName:   form.LastName,
			FormErrors: vErrs,
		}
		tmpl := s.lookupTemplate("registration-form.html")
		if tmpl == nil {
			logger.Error(ult)
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
		if err := tmpl.Execute(w, data); err != nil {
			logger.Error(ewte + err.Error())
			http.Redirect(w, r, "/admin/registration", http.StatusSeeOther)
		}
		return
	}

	_, err = s.st.RegisterUser(r.Context(), storage.Users{
		UserName:     form.UserName,
		Email:        form.Email,
		Password:     form.Password,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		CRUDTimeDate: storage.CRUDTimeDate{CreatedBy: "123", UpdatedBy: "123"},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/admin"+userListPath, http.StatusSeeOther)

}
func (s *Server) passwordRecoverForm(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "forgot-password.html")
}

func (s *Server) screenLockForm(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "screen-lock.html")
}

func formTemplate(s *Server, w http.ResponseWriter, r *http.Request, tmp string) {
	tmpl := s.lookupTemplate(tmp)
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}

func checkLogin(s *Server, email string, pass string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetUserInfoBy(context.Background(), email)
		if resp != nil {
			err := bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(pass))
			if err != nil {
				return fmt.Errorf(" Password doesn't match")
			}
		} else {
			return fmt.Errorf(" Please Enter valid credential")
		}
		return nil
	}
}
