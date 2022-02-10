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

type StationTempData struct {
	CSRFField    template.HTML
	Form         StationForm
	FormErrors   map[string]string
	Data         []StationForm
	DistrictData []DistrictForm
	FormAction   string
}

func (s StationForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name is required"),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateStation(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error("The Position is required"),
			validation.By(checkStationPosition(srv, s.Position, id)),
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
	)
}

func (s *Server) stationListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("station list")
	tmpl := s.lookupTemplate("station.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	stnListForm := s.stationList(r, w,false)
	data := StationTempData{
		Data: stnListForm,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) stationFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("station form")
	disdata := s.districtList(r, w, true)
	data := StationTempData{
		CSRFField:    csrf.TemplateField(r),
		DistrictData: disdata,
	}
	fmt.Println(disdata)
	json.NewEncoder(w).Encode(data)

}

func (s *Server) submitStationHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("station submit")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form StationForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	fmt.Println("###### FORM DATA########")
	fmt.Printf("\n%+v\n\n", form)
	fmt.Println("###### FORM DATA########")
	if err := form.Validate(s, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := StationTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateStation(r.Context(), storage.Station{
		Name:       trim(form.Name),
		Status:     form.Status,
		Position:   form.Position,
		DistrictID: form.DistrictID,
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

func (s *Server) updateStationFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update station form")
	params := mux.Vars(r)
	id := params["id"]
	disdata := s.districtList(r, w, true)
	stnForm := s.strgToFormgetStationByID(r, id, w)
	data := StationTempData{
		Form:         stnForm,
		DistrictData: disdata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateStationHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update station")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form StationForm
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
		data := StationTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.Station{
		ID:         id,
		Name:       trim(form.Name),
		Status:     form.Status,
		DistrictID: form.DistrictID,
		Position:   form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "123",
		},
	}
	_, err := s.st.UpdateStation(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update station data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewStationHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view station form")
	params := mux.Vars(r)
	id := params["id"]
	stnForm := s.strgToFormgetStationByID(r, id, w)
	data := StationTempData{
		Form: stnForm,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateStationStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update station status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetStationBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get station info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateStationStatus(r.Context(), storage.Station{
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
		_, err := s.st.UpdateStationStatus(r.Context(), storage.Station{
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
	http.Redirect(w, r, "/admin/"+stationListPath, http.StatusSeeOther)
}

func (s *Server) deleteStationHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete station")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteStation(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete station" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) strgToFormgetStationByID(r *http.Request, id string, w http.ResponseWriter) StationForm {
	res, err := s.st.GetStationBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get station " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Form := StationForm{
		ID:           res.ID,
		Name:         res.Name,
		DistrictID:   res.DistrictID,
		DistrictName: res.DistrictName.String,
		Status:       res.Status,
		Position:     res.Position,
		CreatedBy:    res.CreatedBy,
	}
	return Form
}
