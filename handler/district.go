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

type DistrictTempData struct {
	CSRFField   template.HTML
	Form        DistrictForm
	FormErrors  map[string]string
	Data        []DistrictForm
	CountryData []CountryForm
	FormAction  string
}

func (s DistrictForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateDistrict(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error("The Position is required"),
			validation.By(checkDistrictPosition(srv, s.Position, id)),
		),
		validation.Field(&s.Status,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("Status is Invalid"),
			validation.Max(2).Error("Status is Invalid"),
		),
		validation.Field(&s.CountryID,
			validation.Required.Error("The Country name is required"),
			validation.By(checkCountryExists(srv, s.CountryID)),
		),
	)
}

func (s *Server) districtListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("district list")
	tmpl := s.lookupTemplate("district.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	disList, err := s.st.GetDistrictList(r.Context(), false) // false = active + inactive
	if err != nil {
		logger.Error("error while get designtion : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	disListForm := storageToDisForm(disList)
	data := DistrictTempData{
		Data: disListForm,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) districtFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("district submit")
	cntrydata := s.countryList(r, w, true)
	data := DistrictTempData{
		CSRFField:   csrf.TemplateField(r),
		CountryData: cntrydata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitDistrictHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("district submit")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form DistrictForm
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
		data := DistrictTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateDistrict(r.Context(), storage.District{
		Name:      trim(form.Name),
		Status:    form.Status,
		Position:  form.Position,
		CountryID: form.CountryID,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: "123",
			UpdatedBy: "123",
		},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateDistrictFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update district form")
	params := mux.Vars(r)
	id := params["id"]
	cntrydata := s.countryList(r, w, true)
	disData := s.strToDisFormByID(w, r, id)
	data := DistrictTempData{
		CSRFField:   csrf.TemplateField(r),
		Form:        disData,
		CountryData: cntrydata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDistrictHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update district")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form DistrictForm
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
		data := DistrictTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.District{
		ID:        id,
		Name:      trim(form.Name),
		Status:    form.Status,
		CountryID: form.CountryID,
		Position:  form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "123",
		},
	}
	_, err := s.st.UpdateDistrict(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update district data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewDistrictHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view district form")
	params := mux.Vars(r)
	id := params["id"]
	disData := s.strToDisFormByID(w, r, id)
	data := DistrictTempData{
		Form: disData,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDistrictStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update district status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetDistrictBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get district info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateDistrictStatus(r.Context(), storage.District{
			ID:     id,
			Status: 2,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: "123",
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	} else {
		_, err := s.st.UpdateDistrictStatus(r.Context(), storage.District{
			ID:     id,
			Status: 1,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: "123",
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	}
	http.Redirect(w, r, "/admin/"+districtListPath, http.StatusSeeOther)
}

func (s *Server) deleteDistrictHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete district")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteDistrict(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete district" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func storageToDisForm(disList []storage.District) []DistrictForm {
	disListForm := make([]DistrictForm, 0)
	for _, item := range disList {
		disdata := DistrictForm{
			ID:          item.ID,
			Name:        item.Name,
			Status:      item.Status,
			Position:    item.Position,
			CountryID:   item.CountryID,
			CreatedAt:   item.CreatedAt,
			CreatedBy:   item.CreatedBy,
			UpdatedAt:   item.UpdatedAt,
			UpdatedBy:   item.UpdatedBy,
			CountryName: item.CountryName.String,
		}
		disListForm = append(disListForm, disdata)
	}
	return disListForm
}
