package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	UserName   string
	Email      string
	Password   string
	FormErrors map[string]string
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
func (s *Server) adminindex(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "index.html")
}

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "login.html")
}

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
		data := LoginForm{
			UserName:   form.UserName,
			Password:   form.Password,
			Email:      form.Email,
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

	http.Redirect(w, r, "/admin/index", http.StatusSeeOther)

}

func (s *Server) screenLockForm(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "screen-lock.html")
}

func (s *Server) registrationForm(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "registration-form.html")
}

func (s *Server) passwordRecoverForm(w http.ResponseWriter, r *http.Request) {
	formTemplate(s, w, r, "forgot-password.html")
}

// func loadLoginTemplate(s *Server, w http.ResponseWriter, r *http.Request, data LoginForm) {
// 	tmpl := s.lookupTemplate("login.html")
// 	if tmpl == nil {
// 		logger.Error(ult)
// 		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
// 	}
// 	if err := tmpl.Execute(w, data); err != nil {
// 		logger.Error(ewte + err.Error())
// 		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
// 	}
// }

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
