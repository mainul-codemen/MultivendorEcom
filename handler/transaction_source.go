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

type TransactionSourceTempData struct {
	CSRFField    template.HTML
	Form         TransactionSourceForm
	FormErrors   map[string]string
	Data         []TransactionSourceForm
	DistrictData []DistrictForm
	CountryData  []CountryForm
	StationData  []StationForm
}

func (s TransactionSourceForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.TransactionSourceName,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateTransactionSource(srv, s.TransactionSourceName, id)),
		),
	)
}

func (s *Server) transactionSourceListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("transactionSource list")
	tmpl := s.lookupTemplate("transaction-source.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	transactionSourceList, err := s.st.GetTransactionSource(r.Context(), false)
	if err != nil {
		logger.Error("error while get transactionSource : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	transactionSourceListForm := s.storageToTranSrcForm(transactionSourceList)
	data := TransactionSourceTempData{
		Data: transactionSourceListForm,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) submitTransactionSourceHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("transactionSource submit")
	uid, _ := s.GetSetSessionValue(r, w)
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form TransactionSourceForm
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
		data := TransactionSourceTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateTransactionSource(r.Context(), storage.TransactionSource{
		TransactionSourceName: form.TransactionSourceName,
		Status:                form.Status,
		CRUDTimeDate:          storage.CRUDTimeDate{CreatedBy: uid, UpdatedBy: uid},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateTransactionSourceHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update transactionSource")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form TransactionSourceForm
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
		data := TransactionSourceTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.TransactionSource{
		ID:                    id,
		TransactionSourceName: form.TransactionSourceName,
		Status:                form.Status,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	}
	_, err := s.st.UpdateTransactionSource(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update transactionSource data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewTransactionSourceHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view transactionSource form")
	params := mux.Vars(r)
	id := params["id"]
	brnFrm := s.getTransactionSourceInfo(r, id, w)
	data := TransactionSourceTempData{
		Form: brnFrm,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateTransactionSourceStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update transactionSource status")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetTransactionSourceBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get transactionSource info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateTransactionSourceStatus(r.Context(), storage.TransactionSource{
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
		_, err := s.st.UpdateTransactionSourceStatus(r.Context(), storage.TransactionSource{
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
	http.Redirect(w, r, "/admin/"+transactionSourceListPath, http.StatusSeeOther)
}

func (s *Server) deleteTransactionSourceHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete transactionSource")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteTransactionSource(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete transactionSource" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) getTransactionSourceInfo(r *http.Request, id string, w http.ResponseWriter) TransactionSourceForm {
	res, err := s.st.GetTransactionSourceBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get transactionSource " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	form := TransactionSourceForm{
		ID:                    id,
		TransactionSourceName: res.TransactionSourceName,
		Status:                res.Status,
		CreatedAt:             res.CreatedAt,
		CreatedBy:             res.CreatedBy,
		UpdatedAt:             res.UpdatedAt,
		UpdatedBy:             res.UpdatedBy,
	}
	return form
}
