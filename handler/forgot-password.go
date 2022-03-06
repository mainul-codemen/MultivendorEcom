package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/MultivendorEcom/serviceutil/logger"
	lutil "github.com/MultivendorEcom/serviceutil/loginUtils"
	mail "github.com/MultivendorEcom/serviceutil/mail"
	"github.com/MultivendorEcom/serviceutil/otp"
	"github.com/MultivendorEcom/storage"
	"github.com/gorilla/csrf"
)

const fpt = "forgot-password.html"
const fpt2 = "forgot-password-2.html"
const pref = "send-confirmation-forgot-password-.html"
const mailmsg = "This is your Reset Password Link.Plese Reset password With bellow link"
const sub = "This is Password Reset Email"
const rpl = "http://localhost:8090/public/recovery-password-2"

type PassRecForm struct {
	CSRFField  template.HTML
	Email      string
	Password   string
	FormErrors map[string]string
}

func (s *Server) passwordRecoverForm(w http.ResponseWriter, r *http.Request) {
	logger.Info("password recovery form")
	formTemplate(s, w, r, fpt)
}

func (s *Server) submitpasswordRecover(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit password recovery")
	tmpl := s.lookupTemplate(fpt)
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	err := r.ParseForm()
	if err != nil {
		logger.Info(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusBadRequest)
	}
	var form PassRecForm
	err = s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error("unable to decode information" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	formErrs := make(map[string]string)
	if form.Email == "" {
		logger.Info("could not find user information")
		formErrs["Email"] = "Email is can't be empty"
		data := PassRecForm{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: formErrs,
		}
		forgotFormTemplate(s, w, r, fpt, data)
		return
	}
	res, err := s.st.GetUserInfoBy(context.Background(), form.Email)
	if err != nil || res == nil {
		logger.Info("could not find user information")
		formErrs["Email"] = "Email is not registered"
		data := PassRecForm{
			CSRFField:  csrf.TemplateField(r),
			Email:      form.Email,
			FormErrors: formErrs,
		}
		forgotFormTemplate(s, w, r, fpt, data)
		return
	} else {
		uid := res.ID
		tkn := otp.GenerateRandomToken()
		s.st.PassResRequest(context.Background(), storage.PassResRequest{UserID: res.ID, Token: tkn})
		mail.SendingMail(form.Email, mail.MailStruct{
			Token:              tkn,
			ResetPasswordLinks: fmt.Sprintf("%s?token=%s&uvxy=%s", rpl, tkn, uid),
			Message:            mailmsg,
			Subject:            sub,
		})
	}	
	forgotFormTemplate(s, w, r, pref, PassRecForm{})
}

func (s *Server) savePasswordResetForm(w http.ResponseWriter, r *http.Request) {
	logger.Info("password recovery save form")
	forgotFormTemplate(s, w, r, fpt2, PassRecForm{})
}

func (s *Server) savePasswordReset(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit password recovery")
	tmpl := s.lookupTemplate(fpt2)
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	err := r.ParseForm()
	if err != nil {
		logger.Info(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusBadRequest)
	}

	var form PassRecForm
	err = s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error("unable to decode information" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}

	u, _ := url.Parse(r.Referer())
	values, _ := url.ParseQuery(u.RawQuery)
	token := values.Get("token")
	uid := values.Get("uvxy")
	formErrs := make(map[string]string)
	if form.Password == "" {
		logger.Info("password can not be empty")
		formErrs["Password"] = "Password can't be empty"
		data := PassRecForm{
			CSRFField:  csrf.TemplateField(r),
			Password:   form.Password,
			FormErrors: formErrs,
		}
		forgotFormTemplate(s, w, r, fpt, data)
		return
	} else {
		_, err := s.st.GetPassResRequestInfo(r.Context(), uid, token)
		if err != nil {
			logger.Error("Error while get password reset information" + err.Error())
			http.Redirect(w, r, pref, http.StatusBadRequest)
			return
		}
		hp, err := lutil.HashPassword(trim(form.Password))
		if err != nil {
			log.Fatalln("Unable to hashed password")
		}
		_, err = s.st.SavePasswordReset(r.Context(), storage.PassResRequest{
			UserID:   uid,
			Password: hp,
		})
		if err != nil {
			logger.Error("Error while get password reset information" + err.Error())
			http.Redirect(w, r, pref, http.StatusBadRequest)
			return
		}
		tmpl := s.lookupTemplate("login.html")
		if tmpl == nil {
			logger.Error(ult)
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
		if err := tmpl.Execute(w, nil); err != nil {
			logger.Error(ewte + err.Error())
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
	}
}

func forgotFormTemplate(s *Server, w http.ResponseWriter, r *http.Request, tmp string, data PassRecForm) {
	tmpl := s.lookupTemplate(tmp)
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}
