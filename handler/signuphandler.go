package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type UsrTempData struct {
	CSRFField       template.HTML
	Form            UserForm
	FormErrors      map[string]string
	Data            []UserForm
	DistrictData    []DistrictForm
	CountryData     []CountryForm
	StationData     []StationForm
	DesignationData []DesignationForm
	HubData         []HubForm
	FormAction      string
}

func (s UserForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.UserName,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.By(checkDuplicateHub(srv, s.UpdatedBy, id)),
		),
		validation.Field(&s.Phone1,
			validation.Required.Error("The User phone is required"),
			validation.Length(3, 11).Error("Please insert phone between 3 to 11"),
			validation.Match(regexp.MustCompile("^[0-9_ ]*$")).Error("Must be digit. No alphabet is allowed."),
			validation.By(checkDuplicateHubPhone(srv, s.Phone1, id)),
		),
	)
}

func (s *Server) userListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("user list")
	tmpl := s.lookupTemplate("user.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	usrdata := s.usrList(r, w, false)
	data := UsrTempData{
		Data: usrdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}

func (s *Server) usrFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("usr form handler")
	tmpl := s.lookupTemplate("user-create.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	cntrydata := s.countryList(r, w, true) // true = only active status
	disdata := s.districtList(r, w, true)  // true = only active status
	stndata := s.stationList(r, w, true)   // true = only active status
	desdata := s.desList(r, w, true)       // true = only active status
	data := UsrTempData{
		CSRFField:       csrf.TemplateField(r),
		CountryData:     cntrydata,
		DistrictData:    disdata,
		StationData:     stndata,
		DesignationData: desdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}

func (s *Server) submitUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("user submit handler")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form UserForm
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
		data := UsrTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.RegisterUser(r.Context(), storage.Users{
		DesignationID:           form.DesignationID,
		UserRole:                form.UserRole,
		EmployeeRole:            form.EmployeeRole,
		VerifiedBy:              "123",
		JoinBy:                  "124",
		CountryID:               form.CountryID,
		DistrictID:              form.DistrictID,
		StationID:               form.StationID,
		Status:                  form.Status,
		UserName:                trim(form.UserName),
		FirstName:               trim(form.FirstName),
		LastName:                trim(form.LastName),
		Email:                   trim(form.Email),
		EmailVerifiedAt:         form.EmailVerifiedAt,
		Password:                trim(form.Password),
		Phone1:                  trim(form.Phone1),
		Phone2:                  trim(form.Phone2),
		PhoneNumberVerifiedAt:   time.Now(),
		PhoneNumberVerifiedCode: trim(form.PhoneNumberVerifiedCode),
		DateOfBirth:             form.DateOfBirth,
		Gender:                  form.Gender,
		FBID:                    trim(form.FBID),
		Photo:                   trim(form.Photo),
		NIDFrontPhoto:           trim(form.NIDFrontPhoto),
		NIDBackPhoto:            trim(form.NIDBackPhoto),
		NIDNumber:               trim(form.NIDNumber),
		CVPDF:                   trim(form.CVPDF),
		PresentAddress:          trim(form.PresentAddress),
		PermanentAddress:        trim(form.PermanentAddress),
		Reference:               trim(form.Reference),
		RememberToken:           trim(form.RememberToken),
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

// func (s *Server) updateUserFormHandler(w http.ResponseWriter, r *http.Request) {
// 	logger.Info("view update user form")
// 	params := mux.Vars(r)
// 	id := params["id"]
// 	disdata := s.districtList(r, w, true)
// 	cntrydata := s.countryList(r, w, true)
// 	stndata := s.stationList(r, w, true)
// 	hubInfo := s.hubList(r, w, true)
// 	uinfo, err := s.st.GetUserInfoBy(r.Context(), id)
// 	data := UsrTempData{
// 		CSRFField:    csrf.TemplateField(r),
// 		HubData:      hubInfo,
// 		Data:         []UserForm{},
// 		DistrictData: disdata,
// 		CountryData:  cntrydata,
// 		StationData:  stndata,
// 	}
// 	json.NewEncoder(w).Encode(data)
// }

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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
		data := UsrTempData{
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

func (s *Server) viewUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view hub form")
	// params := mux.Vars(r)
	// id := params["id"]
	// hubinfo := s.strgToHubByID(r, id, w)
	data := UsrTempData{
		// Form: hubinfo,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateUserStatusHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

// func (s *Server) strgToHubByID(r *http.Request, id string, w http.ResponseWriter) HubForm {
// 	res, err := s.st.GetHubBy(r.Context(), id)
// 	if err != nil {
// 		logger.Error("error while get hub " + err.Error())
// 		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
// 	}
// 	Form := HubForm{
// 		ID:           id,
// 		Name:         res.HubName,
// 		CountryID:    res.CountryID,
// 		CountryName:  res.CountryName.String,
// 		DistrictID:   res.DistrictID,
// 		DistrictName: res.DistrictName.String,
// 		StationID:    res.StationID,
// 		StationName:  res.StationName.String,
// 		HubPhone1:    res.HubPhone1,
// 		HubPhone2:    res.HubPhone2,
// 		HubEmail:     res.HubEmail,
// 		HubAddress:   res.HubAddress,
// 		Status:       res.Status,
// 		Position:     res.Position,
// 		CreatedAt:    res.CreatedAt,
// 		CreatedBy:    res.CreatedBy,
// 		UpdatedAt:    res.UpdatedAt,
// 		UpdatedBy:    res.UpdatedBy,
// 	}
// 	return Form
// }
