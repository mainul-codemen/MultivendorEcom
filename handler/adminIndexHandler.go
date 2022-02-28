package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	UserName   string
	Email      string
	Password   string
	FormErrors map[string]string
}

type RegistrationForm struct {
	ID         string
	UserName   string
	FirstName  string
	LastName   string
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
		data := RegistrationForm{
			UserName: form.UserName,
			Email:    form.Email,
			Password: form.Password,
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
					fmt.Println(value)
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
