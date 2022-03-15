package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type IncomeTempData struct {
	CSRFField   template.HTML
	Form        IncomeForm
	FormErrors  map[string]string
	Data        []IncomeForm
	Accounts    []AccountsForm
	FormAction  string
	FormMessage map[string]string
}

func (s IncomeForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Title,
			validation.Required.Error("Title is required"),
			validation.Length(3, 50).Error("Please insert title between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateAccount(srv, s.Title, id)),
		),
		validation.Field(&s.AccountID,
			validation.Required.Error("The Account Number is required"),
		),
		validation.Field(&s.IncomeAmount,
			validation.Required.Error("The Amount is required"),
		),
	)
}

func (s *Server) incomeFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("income submit")
	s.GetSetSessionValue(r, w)
	disdata := s.accountsList(r, w, false)
	data := IncomeTempData{
		CSRFField: csrf.TemplateField(r),
		Accounts:  disdata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitIncomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit income")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form IncomeForm
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
		data := IncomeTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	_, err := s.st.CreateIncome(r.Context(), storage.Income{
		Title:        trim(form.Title),
		AccountID:    form.AccountID,
		IncomeAmount: form.IncomeAmount,
		Note:         trim(form.Note),
		IncomeDate:   s.stringToDate(form.IncomeDate),
		Status:       1,
		CRUDTimeDate: storage.CRUDTimeDate{CreatedBy: uid, UpdatedBy: uid},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)

}

func (s *Server) incomeListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("income list")
	tmpl := s.lookupTemplate("income.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	actdata := s.incomeList(r, w, false)
	data := IncomeTempData{
		Data: actdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) updateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update income")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form IncomeForm
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
		data := IncomeTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err := s.st.UpdateIncome(r.Context(), storage.Income{
		ID:           id,
		Title:        trim(form.Title),
		AccountID:    form.AccountID,
		IncomeAmount: form.IncomeAmount,
		Note:         trim(form.Note),
		IncomeDate:   s.stringToDate(form.IncomeDate),
		Status:       1,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error("error while update income data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewIncomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view income form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetIncomeBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get income " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := IncomeTempData{
		Form: IncomeForm{
			ID:            id,
			Title:         res.Title,
			AccountID:     res.AccountID,
			Note:          res.Note,
			AccountNumber: res.AccountNumber,
			AccountName:   res.AccountName,
			IncomeAmount:  res.IncomeAmount,
			IncomeDate:    id,
			Status:        res.Status,
			CreatedAt:     res.CreatedAt,
			CreatedBy:     res.CreatedBy,
			UpdatedAt:     res.UpdatedAt,
			UpdatedBy:     res.UpdatedBy,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateIncomeStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update income status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	res, err := s.st.GetIncomeBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get income info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateIncomeStatus(r.Context(), storage.Income{
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
		_, err := s.st.UpdateIncomeStatus(r.Context(), storage.Income{
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
	http.Redirect(w, r, "/admin/"+incomeListPath, http.StatusSeeOther)
}

func (s *Server) deleteIncomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete income")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteIncome(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete income" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	json.NewEncoder(w).Encode(msg)
}
