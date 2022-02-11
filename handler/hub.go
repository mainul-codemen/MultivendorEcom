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

type HubTempData struct {
	CSRFField    template.HTML
	Form         HubForm
	FormErrors   map[string]string
	Data         []HubForm
	DistrictData []DistrictForm
	CountryData  []CountryForm
	StationData  []StationForm
	FormAction   string
}

func (s HubForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateHub(srv, s.Name, id)),
		),
		validation.Field(&s.HubPhone1,
			validation.Required.Error("The hub phone 1 is required"),
			validation.Length(3, 11).Error("Please insert name between 3 to 11"),
			validation.Match(regexp.MustCompile("^[0-9_ ]*$")).Error("Must be digit. No alphabet is allowed."),
			validation.By(checkDuplicateHubPhone(srv, s.HubPhone1, id)),
		),
		validation.Field(&s.HubEmail,
			validation.Required.Error("The name email is required"),
			validation.Length(3, 40).Error("Please insert name between 3 to 40"),
			validation.By(checkDuplicateHubEmail(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error(posReq),
			validation.By(checkHubPosition(srv, s.Position, id)),
		),
		validation.Field(&s.HubAddress,
			validation.Required.Error("The hub address is required"),
		),
		validation.Field(&s.Status,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("Status is Invalid"),
			validation.Max(2).Error("Status is Invalid"),
		),
		validation.Field(&s.DistrictID,
			validation.Required.Error("The District name is required"),
			validation.By(checkDistrictExists(srv, s.DistrictID)),
		),
		validation.Field(&s.CountryID,
			validation.Required.Error("The country name is required"),
			validation.By(checkCountryExists(srv, s.CountryID)),
		),
		validation.Field(&s.StationID,
			validation.Required.Error("The station name is required"),
			validation.By(checkStationExists(srv, s.StationID)),
		),
	)
}

func (s *Server) hubListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("hub list")
	tmpl := s.lookupTemplate("hub.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	hubdata := s.hubList(r, w, false)
	data := HubTempData{
		Data: hubdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}

func (s *Server) hubFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("hub form handler")
	cntrydata := s.countryList(r, w, true)
	disdata := s.districtList(r, w, true)
	stndata := s.stationList(r, w, true)
	data := HubTempData{
		CSRFField:    csrf.TemplateField(r),
		CountryData:  cntrydata,
		DistrictData: disdata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitHubHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("hub submit handler")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form HubForm
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
		data := HubTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateHub(r.Context(), storage.Hub{
		HubName:    trim(form.Name),
		CountryID:  form.CountryID,
		DistrictID: form.DistrictID,
		StationID:  form.StationID,
		HubPhone1:  trim(form.HubPhone1),
		HubPhone2:  trim(form.HubPhone2),
		HubEmail:   trim(form.HubEmail),
		HubAddress: trim(form.HubAddress),
		Status:     form.Status,
		Position:   form.Position,
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

func (s *Server) updateHubFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update hub form")
	params := mux.Vars(r)
	id := params["id"]
	disdata := s.districtList(r, w, true)
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w, true)
	hubInfo := s.strgToHubByID(r, id, w)
	data := HubTempData{
		CSRFField:    csrf.TemplateField(r),
		Form:         hubInfo,
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateHubHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update hub")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form HubForm
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
		data := HubTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.Hub{
		ID:         id,
		HubName:    trim(form.Name),
		CountryID:  form.CountryID,
		DistrictID: form.DistrictID,
		StationID:  form.StationID,
		HubPhone1:  trim(form.HubPhone1),
		HubPhone2:  trim(form.HubPhone2),
		HubEmail:   trim(form.HubEmail),
		HubAddress: trim(form.HubAddress),
		Status:     form.Status,
		Position:   form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "123",
		},
	}
	_, err := s.st.UpdateHub(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update hub data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewHubHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view hub form")
	params := mux.Vars(r)
	id := params["id"]
	hubinfo := s.strgToHubByID(r, id, w)
	data := HubTempData{
		Form: hubinfo,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateHubStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update hub status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetHubBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get hub info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateHubStatus(r.Context(), storage.Hub{
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
		_, err := s.st.UpdateHubStatus(r.Context(), storage.Hub{
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
	http.Redirect(w, r, "/admin/"+hubListPath, http.StatusSeeOther)
}

func (s *Server) deleteHubHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete hub")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteHub(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete hub" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) strgToHubByID(r *http.Request, id string, w http.ResponseWriter) HubForm {
	res, err := s.st.GetHubBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get hub " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Form := HubForm{
		ID:           id,
		Name:         res.HubName,
		CountryID:    res.CountryID,
		CountryName:  res.CountryName.String,
		DistrictID:   res.DistrictID,
		DistrictName: res.DistrictName.String,
		StationID:    res.StationID,
		StationName:  res.StationName.String,
		HubPhone1:    res.HubPhone1,
		HubPhone2:    res.HubPhone2,
		HubEmail:     res.HubEmail,
		HubAddress:   res.HubAddress,
		Status:       res.Status,
		Position:     res.Position,
		CreatedAt:    res.CreatedAt,
		CreatedBy:    res.CreatedBy,
		UpdatedAt:    res.UpdatedAt,
		UpdatedBy:    res.UpdatedBy,
	}
	return Form
}
