package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type IncomeTaxTempData struct {
	CSRFField   template.HTML
	Form        IncomeTaxForm
	FormErrors  map[string]string
	Data        []IncomeTaxForm
	Accounts    []AccountsForm
	FormAction  string
	FormMessage map[string]string
}

func (s IncomeTaxForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.AccountID,
			validation.Required.Error("The Account Number is required"),
		),
		validation.Field(&s.IncomeTaxDate,
			validation.Required.Error("The date is required"),
		),
		validation.Field(&s.TaxAmount,
			validation.Required.Error("The amount is required"),
		),
	)
}

func (s *Server) incomeTaxFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("income submit")
	s.GetSetSessionValue(r, w)
	disdata := s.accountsList(r, w, false)
	data := IncomeTaxTempData{
		CSRFField: csrf.TemplateField(r),
		Accounts:  disdata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitIncomeTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit income")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form IncomeTaxForm
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
		data := IncomeTaxTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	_, err := s.st.CreateIncomeTax(r.Context(), storage.IncomeTax{
		AccountID:        form.AccountID,
		TaxReceiptNumber: form.TaxReceiptNumber,
		Status:           1,
		IncomeTaxDate:    s.stringToDate(form.IncomeTaxDate),
		TaxAmount:        0,
		CRUDTimeDate:     storage.CRUDTimeDate{CreatedBy: uid, UpdatedBy: uid},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	ttdb, tdb := trnsTypesAndSource(s, r, w, "Cash Out", "Income Tax")
	_, err = s.st.CreateAccountsTransaction(r.Context(), storage.AccountsTransaction{
		FromAccountID:     form.AccountID,
		TransactionAmount: form.TaxAmount,
		TransactionType:   ttdb.ID,
		TransactionSource: tdb.ID,
		Status:            1,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: uid,
		},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)

}

func (s *Server) incomeTaxListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("income list")
	tmpl := s.lookupTemplate("income-tax.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	data := IncomeTaxTempData{
		Data:     s.incomeTaxList(r, w, false),
		Accounts: s.accountsList(r, w, false),
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}

func (s *Server) updateIncomeTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update income tax")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form IncomeTaxForm
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
		data := IncomeTaxTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err := s.st.UpdateIncomeTax(r.Context(), storage.IncomeTax{
		ID:               id,
		AccountID:        form.AccountID,
		TaxReceiptNumber: trim(form.TaxReceiptNumber),
		Status:           1,
		IncomeTaxDate:    s.stringToDate(form.IncomeTaxDate),
		TaxAmount:        form.TaxAmount,
		CRUDTimeDate:     storage.CRUDTimeDate{UpdatedBy: uid},
	})
	if err != nil {
		logger.Error("error while update income data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewIncomeTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view income form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetIncomeTaxBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get income " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := IncomeTaxTempData{
		Form: IncomeTaxForm{
			ID:            id,
			AccountID:     res.AccountID,
			AccountNumber: res.AccountNumber,
			AccountName:   res.AccountName,
			IncomeTaxDate: id,
			TaxAmount:     res.TaxAmount,
			Status:        res.Status,
			CreatedAt:     res.CreatedAt,
			CreatedBy:     res.CreatedBy,
			UpdatedAt:     res.UpdatedAt,
			UpdatedBy:     res.UpdatedBy,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateIncomeTaxStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update income status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	res, err := s.st.GetIncomeTaxBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get income info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateIncomeTaxStatus(r.Context(), storage.IncomeTax{
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
		_, err := s.st.UpdateIncomeTaxStatus(r.Context(), storage.IncomeTax{
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

func (s *Server) deleteIncomeTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete income")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteIncomeTax(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete income" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	json.NewEncoder(w).Encode(msg)
}
