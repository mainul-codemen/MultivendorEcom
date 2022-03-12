package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type AccountsTempData struct {
	CSRFField   template.HTML
	Form        AccountsForm
	FormErrors  map[string]string
	Data        []AccountsForm
	FormAction  string
	FormMessage map[string]string
}

func (s AccountsForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.AccountName,
			validation.Required.Error("Account name is required"),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateAccount(srv, s.AccountName, id)),
		),
		validation.Field(&s.AccountNumber,
			validation.Required.Error("The Account Number is required"),
		),
		validation.Field(&s.Amount,
			validation.Required.Error("The Amount is required"),
		),
	)
}
func (s *Server) submitAccountsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit accounts")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form AccountsForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	if err := form.Validate(s, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := AccountsTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	_, err := s.st.CreateAccunt(r.Context(), storage.Accounts{
		AccountName:          trim(form.AccountName),
		AccountNumber:        trim(form.AccountNumber),
		AccountVisualization: trim(form.AccountVisualization),
		Amount:               form.Amount,
		Status:               1,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: uid,
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)

}

func (s *Server) accountsListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("accounts list")
	tmpl := s.lookupTemplate("accounts.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	actdata := s.accountsList(r, w, false)
	data := AccountsTempData{
		Data: actdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) addMoney(w http.ResponseWriter, r *http.Request) {
	logger.Info("add money handler")
	uid, _ := s.GetSetSessionValue(r, w)
	var form AccountsForm

	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	fmt.Println(r)
	fmt.Printf("############### %+v", form)
	validationAddMoney(form, r, w)
	_, err := s.st.AddMoney(r.Context(), storage.Accounts{
		AccountNumber: form.AccountNumber,
		Amount:        form.Amount,
		Status:        1,
		CRUDTimeDate:  storage.CRUDTimeDate{UpdatedBy: uid},
	})
	if err != nil {
		logger.Error("error while update accounts data ." + err.Error())
	}

	json.NewEncoder(w).Encode(msg)

}

func validationAddMoney(form AccountsForm, r *http.Request, w http.ResponseWriter) {
	vErrs := map[string]string{}
	if form.AccountNumber == "" {
		vErrs["AccountNumber"] = "Account Number is Required"
	}
	if form.Amount == 0 {
		vErrs["Amount"] = "Amount is required"
	}
	if len(vErrs) > 0 {
		data := AccountsTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

}
func (s *Server) updateAccountsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update accounts")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form AccountsForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logger.Error(deErr + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}
	if err := form.Validate(s, id); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := AccountsTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err := s.st.UpdateAccounts(r.Context(), storage.Accounts{
		ID:                   id,
		AccountName:          trim(form.AccountName),
		AccountNumber:        trim(form.AccountNumber),
		AccountVisualization: trim(form.AccountVisualization),
		Amount:               form.Amount,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error("error while update accounts data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewAccountsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view accounts form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetAccountsBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get accounts " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := AccountsTempData{
		Form: AccountsForm{
			ID:                   id,
			AccountVisualization: res.AccountVisualization,
			AccountName:          res.AccountName,
			AccountNumber:        res.AccountNumber,
			Amount:               res.Amount,
			Status:               res.Status,
			CreatedAt:            res.CreatedAt,
			CreatedBy:            res.CreatedBy,
			UpdatedAt:            res.UpdatedAt,
			UpdatedBy:            res.UpdatedBy,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateAccountsStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update accounts status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	res, err := s.st.GetAccountsBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get accounts info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateAccountsStatus(r.Context(), storage.Accounts{
			ID:     id,
			Status: 2,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: uid,
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	} else {
		_, err := s.st.UpdateAccountsStatus(r.Context(), storage.Accounts{
			ID:     id,
			Status: 1,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: uid,
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	}
	http.Redirect(w, r, "/admin/"+accountsListPath, http.StatusSeeOther)
}

func (s *Server) deleteAccountsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete accounts")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteAccounts(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete accounts" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	json.NewEncoder(w).Encode(msg)
}
