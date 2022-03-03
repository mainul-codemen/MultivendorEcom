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

type GradeData struct {
	CSRFField   template.HTML
	Form        GradeForm
	FormErrors  map[string]string
	Data        []GradeForm
	FormAction  string
	FormMessage map[string]string
}

func (s GradeForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name is required"),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateGrade(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error("The Position is required"),
			validation.By(checkGradePosition(srv, s.Position, id)),
		),
		validation.Field(&s.Status,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("Status is Invalid"),
			validation.Max(2).Error("Status is Invalid"),
		),
	)
}

func (s *Server) submitGrade(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit grade")
	uid, _ := s.GetSetSessionValue(r, w)
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form GradeForm
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
					fmt.Println(value)
					vErrs[key] = value.Error()
				}
			}
		}
		data := GradeData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	_, err = s.st.CreateGrade(r.Context(), storage.Grade{
		Name:           trim(form.Name),
		Description:    trim(form.Description),
		BasicSalary:    form.BasicSalary,
		LunchAllowance: form.LunchAllowance,
		Transportation: form.Transportation,
		RentAllowance:  form.RentAllowance,
		AbsentPenalty:  form.AbsentPenalty,
		TotalSalary:    form.TotalSalary,
		Status:         form.Status,
		Position:       form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: uid,
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) gradeList(w http.ResponseWriter, r *http.Request) {
	logger.Info("grade list")
	tmpl := s.lookupTemplate("grade.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	dptlist := s.grdList(r, w, false)
	data := GradeData{
		Data: dptlist,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) deleteGrade(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete grade")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteGrade(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete grade" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateGrade(w http.ResponseWriter, r *http.Request) {
	logger.Info("update grade")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form GradeForm
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
		data := GradeData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.Grade{
		ID:             id,
		Name:           trim(form.Name),
		Description:    trim(form.Description),
		BasicSalary:    form.BasicSalary,
		LunchAllowance: form.LunchAllowance,
		Transportation: form.Transportation,
		RentAllowance:  form.RentAllowance,
		AbsentPenalty:  form.AbsentPenalty,
		TotalSalary:    form.TotalSalary,
		Status:         form.Status,
		Position:       form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	}
	_, err := s.st.UpdateGrade(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update designatio data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewGrade(w http.ResponseWriter, r *http.Request) {
	logger.Info("view grade form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetGradeBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get grade " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := GradeData{
		Form: GradeForm{
			ID:             id,
			Name:           res.Name,
			Description:    res.Description,
			BasicSalary:    res.BasicSalary,
			LunchAllowance: res.LunchAllowance,
			RentAllowance:  res.RentAllowance,
			Transportation: res.Transportation,
			AbsentPenalty:  res.AbsentPenalty,
			TotalSalary:    res.TotalSalary,
			Status:         res.Status,
			Position:       res.Position,
			CreatedAt:      res.CreatedAt,
			CreatedBy:      res.CreatedBy,
			UpdatedAt:      res.UpdatedAt,
			UpdatedBy:      res.UpdatedBy,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateGradeStatus(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update grade status")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetGradeBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get grade info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateGradeStatus(r.Context(), storage.Grade{
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
		_, err := s.st.UpdateGradeStatus(r.Context(), storage.Grade{
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
	http.Redirect(w, r, "/admin/"+gradeListPath, http.StatusSeeOther)
}
