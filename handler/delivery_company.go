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

type DeliveryCompanyTempData struct {
	CSRFField    template.HTML
	Form         DeliveryCompanyForm
	FormErrors   map[string]string
	Data         []DeliveryCompanyForm
	DistrictData []DistrictForm
	CountryData  []CountryForm
	StationData  []StationForm
	FormAction   string
}

func (s DeliveryCompanyForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.CompanyName,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateDeliveryCompany(srv, s.CompanyName, id)),
		),
		validation.Field(&s.Phone,
			validation.Required.Error("The deliveryCompany phone 1 is required"),
			validation.Length(3, 11).Error("Please insert name between 3 to 11"),
			validation.Match(regexp.MustCompile("^[0-9_ ]*$")).Error("Must be digit. No alphabet is allowed."),
			validation.By(checkDuplicateDeliveryCompanyPhone(srv, s.Phone, id)),
		),
		validation.Field(&s.Email,
			validation.Required.Error("The name email is required"),
			validation.Length(3, 40).Error("Please insert name between 3 to 40"),
			validation.By(checkDuplicateDeliveryCompanyEmail(srv, s.Email, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error(posReq),
			validation.By(checkDeliveryCompanyPosition(srv, s.Position, id)),
		),
		validation.Field(&s.CompanyStatus,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("DeliveryCompanyStatus is Invalid"),
			validation.Max(2).Error("DeliveryCompanyStatus is Invalid"),
		),
		validation.Field(&s.DistrictID,
			validation.Required.Error("The District name is required"),
			validation.By(checkDistrictExists(srv, s.DistrictID)),
		),
		validation.Field(&s.CountryID,
			validation.Required.Error("The country name is required"),
			validation.By(checkCountryExists(srv, s.DistrictID)),
		),
		validation.Field(&s.StationID,
			validation.Required.Error("The station name is required"),
			validation.By(checkDuplicateStation(srv, s.StationID, id)),
		),
	)
}

func (s *Server) deliveryCompanyListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("deliveryCompany list")
	tmpl := s.lookupTemplate("delivery-company.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	deliveryCompanyList, err := s.st.GetDeliveryCompanyList(r.Context())
	if err != nil {
		logger.Error("error while get deliveryCompany : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	deliveryCompanyListForm := s.storageToDCForm(deliveryCompanyList)
	data := DeliveryCompanyTempData{
		Data: deliveryCompanyListForm,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) deliveryCompanyFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("deliveryCompany submit")
	disdata := s.districtList(r, w, true)
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w,true)
	data := DeliveryCompanyTempData{
		CSRFField:    csrf.TemplateField(r),
		FormErrors:   map[string]string{},
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitDeliveryCompanyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("deliveryCompany submit")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form DeliveryCompanyForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	fmt.Print("##############", form)
	if err := form.Validate(s, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := DeliveryCompanyTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateDeliveryCompany(r.Context(), storage.DeliveryCompany{
		CompanyName:    trim(form.CompanyName),
		CountryID:      form.CountryID,
		DistrictID:     form.DistrictID,
		StationID:      form.StationID,
		Phone:          trim(form.Phone),
		Email:          trim(form.Email),
		CompanyAddress: trim(form.CompanyAddress),
		CompanyStatus:  form.CompanyStatus,
		Position:       form.Position,
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

func (s *Server) updateDeliveryCompanyFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update deliveryCompany form")
	params := mux.Vars(r)
	id := params["id"]
	disdata := s.districtList(r, w, true)
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w,true)
	brnFrm := s.getDeliveryCompanyInfo(r, id, w)
	data := DeliveryCompanyTempData{
		Form:         brnFrm,
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDeliveryCompanyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update deliveryCompany")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form DeliveryCompanyForm
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
		data := DeliveryCompanyTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.DeliveryCompany{
		ID:             id,
		CompanyName:    trim(form.CompanyName),
		CountryID:      form.CountryID,
		DistrictID:     form.DistrictID,
		StationID:      form.StationID,
		Phone:          trim(form.Phone),
		Email:          trim(form.Email),
		CompanyAddress: trim(form.CompanyAddress),
		CompanyStatus:  form.CompanyStatus,
		Position:       form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "123",
		},
	}
	_, err := s.st.UpdateDeliveryCompany(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update deliveryCompany data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewDeliveryCompanyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view deliveryCompany form")
	params := mux.Vars(r)
	id := params["id"]
	brnFrm := s.getDeliveryCompanyInfo(r, id, w)
	data := DeliveryCompanyTempData{
		Form: brnFrm,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDeliveryCompanyStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update deliveryCompany status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetDeliveryCompanyBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get deliveryCompany info " + err.Error())
	}
	if res.CompanyStatus == 1 {
		_, err := s.st.UpdateDeliveryCompanyStatus(r.Context(), storage.DeliveryCompany{
			ID:            id,
			CompanyStatus: 2,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: "123",
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	} else {
		_, err := s.st.UpdateDeliveryCompanyStatus(r.Context(), storage.DeliveryCompany{
			ID:            id,
			CompanyStatus: 1,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: "123",
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	}
	http.Redirect(w, r, "/admin/"+deliveryCompanyListPath, http.StatusSeeOther)
}

func (s *Server) deleteDeliveryCompanyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete deliveryCompany")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteDeliveryCompany(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete deliveryCompany" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) getDeliveryCompanyInfo(r *http.Request, id string, w http.ResponseWriter) DeliveryCompanyForm {
	res, err := s.st.GetDeliveryCompanyBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get deliveryCompany " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Form := DeliveryCompanyForm{
		ID:             id,
		CompanyName:    res.CompanyName,
		CountryID:      res.CountryID,
		CountryName:    res.CountryName.String,
		DistrictID:     res.DistrictID,
		DistrictName:   res.DistrictName.String,
		StationID:      res.StationID,
		StationName:    res.StationName.String,
		Phone:          res.Phone,
		Email:          res.Email,
		CompanyAddress: res.CompanyAddress,
		CompanyStatus:  res.CompanyStatus,
		Position:       res.Position,
		CreatedAt:      res.CreatedAt,
		CreatedBy:      res.CreatedBy,
		UpdatedAt:      res.UpdatedAt,
		UpdatedBy:      res.UpdatedBy,
	}
	return Form
}

func (*Server) storageToDCForm(deliveryCompanyList []storage.DeliveryCompany) []DeliveryCompanyForm {
	deliveryCompanyListForm := make([]DeliveryCompanyForm, 0)
	for _, item := range deliveryCompanyList {
		deliveryCompanyData := DeliveryCompanyForm{
			ID:             item.ID,
			CompanyName:    item.CompanyName,
			CountryID:      item.CountryID,
			CountryName:    item.CountryName.String,
			DistrictID:     item.DistrictID,
			DistrictName:   item.DistrictName.String,
			StationID:      item.StationID,
			StationName:    item.StationName.String,
			Phone:          item.Phone,
			Email:          item.Email,
			CompanyAddress: item.CompanyAddress,
			CompanyStatus:  item.CompanyStatus,
			Position:       item.Position,
			CreatedAt:      item.CreatedAt,
			CreatedBy:      item.CreatedBy,
			UpdatedAt:      item.UpdatedAt,
			UpdatedBy:      item.UpdatedBy,
		}
		deliveryCompanyListForm = append(deliveryCompanyListForm, deliveryCompanyData)
	}
	return deliveryCompanyListForm
}
