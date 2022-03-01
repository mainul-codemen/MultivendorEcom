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

type CountryTempData struct {
	CSRFField   template.HTML
	Form        CountryForm
	FormErrors  map[string]string
	Data        []CountryForm
	FormAction  string
	FormMessage map[string]string
}

func (s CountryForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateCountry(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error("The Position is required"),
			validation.By(checkCountryPosition(srv, s.Position, id)),
		),
		validation.Field(&s.Status,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("Status is Invalid"),
			validation.Max(2).Error("Status is Invalid"),
		),
	)
}
func (s *Server) submitCountry(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit country")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form CountryForm
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
		data := CountryTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err := s.st.CreateCountry(r.Context(), storage.Country{
		Name:     trim(form.Name),
		Status:   form.Status,
		Position: form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: s.GetSetSessionValue(r),
			UpdatedBy: s.GetSetSessionValue(r),
		},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)

}

func (s *Server) countryListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("country list")
	tmpl := s.lookupTemplate("country.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	cntrydata := s.countryList(r, w, false)
	data := CountryTempData{
		Data: cntrydata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) updateCountryHadler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update country")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form CountryForm
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
		data := CountryTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err := s.st.UpdateCountry(r.Context(), storage.Country{
		ID:       id,
		Name:     trim(form.Name),
		Status:   form.Status,
		Position: form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: s.GetSetSessionValue(r),
		},
	})
	if err != nil {
		logger.Error("error while update country data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewCountryHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view country form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetCountryBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get country " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := CountryTempData{
		Form: CountryForm{
			ID:        id,
			Name:      res.Name,
			Status:    res.Status,
			Position:  res.Position,
			CreatedAt: res.CreatedAt,
			CreatedBy: res.CreatedBy,
			UpdatedAt: res.UpdatedAt,
			UpdatedBy: res.UpdatedBy,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateCountryStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update country status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetCountryBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get country info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateCountryStatus(r.Context(), storage.Country{
			ID:     id,
			Status: 2,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: s.GetSetSessionValue(r),
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	} else {
		_, err := s.st.UpdateCountryStatus(r.Context(), storage.Country{
			ID:     id,
			Status: 1,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: s.GetSetSessionValue(r),
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	}
	http.Redirect(w, r, "/admin/"+countryListPath, http.StatusSeeOther)
}
func (s *Server) deleteCountryHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete country")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteCountry(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete country" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	json.NewEncoder(w).Encode(msg)
}
