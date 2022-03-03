package handler

import (
	"context"
	"html/template"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
)

type LoginForm struct {
	CSRFField       template.HTML
	EmailOrUserName string
	Password        string
	FormErrors      map[string]string
}

func (s LoginForm) Validate(srv *Server) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.EmailOrUserName,
			validation.Required.Error("This Field is Required"),
		),
		validation.Field(&s.Password,
			validation.Required.Error("Password is Required"),
			validation.By(checkLogin(srv, s.EmailOrUserName, s.Password)),
		),
	)
}

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	logger.Info("login form handler")
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
			CSRFField:       csrf.TemplateField(r),
			EmailOrUserName: form.EmailOrUserName,
			Password:        form.Password,
			FormErrors:      vErrs,
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
	resp, _ := s.st.GetUserInfoBy(context.Background(), form.EmailOrUserName)
	sesn, _ := s.session.Get(r, "mvec-prod")
	sesn.Values["user_id"] = resp.ID
	sesn.Values["user_role"] = resp.UserRole
	if err := sesn.Save(r, w); err != nil {
		logger.Error("Error while saving Session information" + err.Error())
	}
	http.Redirect(w, r, "/admin/index", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.session.Get(r, "mvec-prod")
	session.Values["user_id"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/public/login", http.StatusSeeOther)
}

func (s *Server) GetSetSessionValue(r *http.Request, w http.ResponseWriter) (string, string) {
	session, _ := s.session.Get(r, "mvec-prod")
	uid := session.Values["user_id"]
	user_role := session.Values["user_role"]
	if uid == nil || user_role == nil {
		http.Redirect(w, r, "/forbiden", http.StatusSeeOther)
	} else {
		usr, _ := s.st.GetUserRoleBy(context.Background(), user_role.(string))
		logger.Info("Logged in user = " + usr.Name)
		return uid.(string), usr.Name
	}
	return "", ""
}
