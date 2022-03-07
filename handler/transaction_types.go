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

type TransactionTypesTempData struct {
	CSRFField    template.HTML
	Form         TransactionTypesForm
	FormErrors   map[string]string
	Data         []TransactionTypesForm
	DistrictData []DistrictForm
	CountryData  []CountryForm
	StationData  []StationForm
}

func (s TransactionTypesForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.TransactionTypesName,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateTransactionTypes(srv, s.TransactionTypesName, id)),
		),
	)
}

func (s *Server) transactionTypesListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("transactionTypes list")
	tmpl := s.lookupTemplate("transaction-types.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	transactionTypesList, err := s.st.GetTransactionTypes(r.Context(), false)
	if err != nil {
		logger.Error("error while get transactionTypes : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	transactionTypesListForm := s.storageToTransForm(transactionTypesList)
	data := TransactionTypesTempData{
		Data: transactionTypesListForm,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) submitTransactionTypesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("transactionTypes submit")
	uid, _ := s.GetSetSessionValue(r, w)
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form TransactionTypesForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
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
		data := TransactionTypesTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateTransactionTypes(r.Context(), storage.TransactionTypes{
		TransactionTypesName: form.TransactionTypesName,
		Status:               form.Status,
		CRUDTimeDate:         storage.CRUDTimeDate{CreatedBy: uid, UpdatedBy: uid},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateTransactionTypesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update transactionTypes")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form TransactionTypesForm
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
		data := TransactionTypesTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.TransactionTypes{
		ID:                   id,
		TransactionTypesName: form.TransactionTypesName,
		Status:               form.Status,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	}
	_, err := s.st.UpdateTransactionTypes(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update transactionTypes data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewTransactionTypesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view transactionTypes form")
	params := mux.Vars(r)
	id := params["id"]
	brnFrm := s.getTransactionTypesInfo(r, id, w)
	data := TransactionTypesTempData{
		Form: brnFrm,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateTransactionTypesStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update transactionTypes status")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetTransactionTypesBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get transactionTypes info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateTransactionTypesStatus(r.Context(), storage.TransactionTypes{
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
		_, err := s.st.UpdateTransactionTypesStatus(r.Context(), storage.TransactionTypes{
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
	http.Redirect(w, r, "/admin/"+transactionTypesListPath, http.StatusSeeOther)
}

func (s *Server) deleteTransactionTypesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete transactionTypes")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteTransactionTypes(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete transactionTypes" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) getTransactionTypesInfo(r *http.Request, id string, w http.ResponseWriter) TransactionTypesForm {
	res, err := s.st.GetTransactionTypesBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get transactionTypes " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	form := TransactionTypesForm{
		ID:                   id,
		TransactionTypesName: res.TransactionTypesName,
		Status:               res.Status,
		CreatedAt:            res.CreatedAt,
		CreatedBy:            res.CreatedBy,
		UpdatedAt:            res.UpdatedAt,
		UpdatedBy:            res.UpdatedBy,
	}
	return form
}
