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

type BranchTempData struct {
	CSRFField    template.HTML
	Form         BranchForm
	FormErrors   map[string]string
	Data         []BranchForm
	DistrictData []DistrictForm
	CountryData  []CountryForm
	StationData  []StationForm
}

func (s BranchForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateBranch(srv, s.Name, id)),
		),
		validation.Field(&s.BranchPhone1,
			validation.Required.Error("The branch phone 1 is required"),
			validation.Length(11, 11).Error("Please insert 11 Number between"),
			validation.Match(regexp.MustCompile("^[0-9_ ]*$")).Error("Must be digit. No alphabet is allowed."),
			validation.By(checkDuplicateBranchPhone(srv, s.BranchPhone1, id)),
		),
		validation.Field(&s.BranchEmail,
			validation.Required.Error("The name email is required"),
			validation.Length(3, 40).Error("Please insert name between 3 to 40"),
			validation.By(checkDuplicateBranchEmail(srv, s.BranchEmail, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error(posReq),
			validation.By(checkBranchPosition(srv, s.Position, id)),
		),
		validation.Field(&s.BranchStatus,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("BranchStatus is Invalid"),
			validation.Max(2).Error("BranchStatus is Invalid"),
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

func (s *Server) branchListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("branch list")
	tmpl := s.lookupTemplate("branch.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	branchList, err := s.st.GetBranchList(r.Context())
	if err != nil {
		logger.Error("error while get branch : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	branchListForm := storageToBLForm(branchList)
	data := BranchTempData{
		Data: branchListForm,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) branchFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("branch submit")
	disdata := s.districtList(r, w, true) // status = active
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w, true)
	data := BranchTempData{
		CSRFField:    csrf.TemplateField(r),
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitBranchHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("branch submit")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form BranchForm
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
		data := BranchTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateBranch(r.Context(), storage.Branch{
		BranchName:    trim(form.Name),
		CountryID:     form.CountryID,
		DistrictID:    form.DistrictID,
		StationID:     form.StationID,
		BranchPhone1:  trim(form.BranchPhone1),
		BranchPhone2:  trim(form.BranchPhone2),
		BranchEmail:   trim(form.BranchEmail),
		BranchAddress: trim(form.BranchAddress),
		BranchStatus:  form.BranchStatus,
		Position:      form.Position,
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

func (s *Server) updateBranchFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update branch form")
	params := mux.Vars(r)
	id := params["id"]
	disdata := s.districtList(r, w, true) // status = active
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w, true)
	brnFrm := s.getBranchInfo(r, id, w)
	data := BranchTempData{
		Form:         brnFrm,
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateBranchHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update branch")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form BranchForm
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
		data := BranchTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.Branch{
		ID:            id,
		BranchName:    trim(form.Name),
		CountryID:     form.CountryID,
		DistrictID:    form.DistrictID,
		StationID:     form.StationID,
		BranchPhone1:  trim(form.BranchPhone1),
		BranchPhone2:  trim(form.BranchPhone2),
		BranchEmail:   trim(form.BranchEmail),
		BranchAddress: trim(form.BranchAddress),
		BranchStatus:  form.BranchStatus,
		Position:      form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "123",
		},
	}
	_, err := s.st.UpdateBranch(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update branch data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewBranchHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view branch form")
	params := mux.Vars(r)
	id := params["id"]
	brnFrm := s.getBranchInfo(r, id, w)
	data := BranchTempData{
		Form: brnFrm,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateBranchStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update branch status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetBranchBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get branch info " + err.Error())
	}
	if res.BranchStatus == 1 {
		_, err := s.st.UpdateBranchStatus(r.Context(), storage.Branch{
			ID:           id,
			BranchStatus: 2,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: "123",
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	} else {
		_, err := s.st.UpdateBranchStatus(r.Context(), storage.Branch{
			ID:           id,
			BranchStatus: 1,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: "123",
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	}
	http.Redirect(w, r, "/admin/"+branchListPath, http.StatusSeeOther)
}

func (s *Server) deleteBranchHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete branch")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteBranch(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete branch" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) getBranchInfo(r *http.Request, id string, w http.ResponseWriter) BranchForm {
	res, err := s.st.GetBranchBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get branch " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Form := BranchForm{
		ID:            id,
		Name:          res.BranchName,
		CountryID:     res.CountryID,
		CountryName:   res.CountryName.String,
		DistrictID:    res.DistrictID,
		DistrictName:  res.DistrictName.String,
		StationID:     res.StationID,
		StationName:   res.StationName.String,
		BranchName:    res.BranchName,
		BranchPhone1:  res.BranchPhone1,
		BranchPhone2:  res.BranchPhone2,
		BranchEmail:   res.BranchEmail,
		BranchAddress: res.BranchAddress,
		BranchStatus:  res.BranchStatus,
		Position:      res.Position,
		CreatedAt:     res.CreatedAt,
		CreatedBy:     res.CreatedBy,
		UpdatedAt:     res.UpdatedAt,
		UpdatedBy:     res.UpdatedBy,
	}
	return Form
}

func storageToBLForm(branchList []storage.Branch) []BranchForm {
	branchListForm := make([]BranchForm, 0)
	for _, item := range branchList {
		branchData := BranchForm{
			ID:            item.ID,
			Name:          item.BranchName,
			CountryID:     item.CountryID,
			CountryName:   item.CountryName.String,
			DistrictID:    item.DistrictID,
			DistrictName:  item.DistrictName.String,
			StationID:     item.StationID,
			StationName:   item.StationName.String,
			BranchPhone1:  item.BranchPhone1,
			BranchPhone2:  item.BranchPhone2,
			BranchEmail:   item.BranchEmail,
			BranchAddress: item.BranchAddress,
			BranchStatus:  item.BranchStatus,
			Position:      item.Position,
			CreatedAt:     item.CreatedAt,
			CreatedBy:     item.CreatedBy,
			UpdatedAt:     item.UpdatedAt,
			UpdatedBy:     item.UpdatedBy,
		}
		branchListForm = append(branchListForm, branchData)
	}
	return branchListForm
}
