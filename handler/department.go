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

type DepartmentData struct {
	CSRFField   template.HTML
	Form        DepartmentForm
	FormErrors  map[string]string
	Data        []DepartmentForm
	FormAction  string
	FormMessage map[string]string
}

func (s DepartmentForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name is required"),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateDepartment(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error("The Position is required"),
			validation.By(checkDepartmentPosition(srv, s.Position, id)),
		),
		validation.Field(&s.Description,
			validation.Required.Error("The Description is required"),
			validation.Length(3, 100).Error("Please insert description between 3 to 100 character"),
		),
		validation.Field(&s.Status,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("Status is Invalid"),
			validation.Max(2).Error("Status is Invalid"),
		),
	)
}

func (s *Server) submitDepartment(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit department")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form DepartmentForm
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
		data := DepartmentData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	_, err = s.st.CreateDepartment(r.Context(), storage.Department{
		Name:        trim(form.Name),
		Description: trim(form.Description),
		Status:      form.Status,
		Position:    form.Position,
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

func (s *Server) departmentList(w http.ResponseWriter, r *http.Request) {
	logger.Info("department list")
	tmpl := s.lookupTemplate("department.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	dptlist := s.dptList(r, w, false)
	data := DepartmentData{
		Data: dptlist,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) deleteDepartment(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete department")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteDepartment(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete department" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateDepartment(w http.ResponseWriter, r *http.Request) {
	logger.Info("update department")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form DepartmentForm
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
		data := DepartmentData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.Department{
		ID:          id,
		Name:        trim(form.Name),
		Description: trim(form.Description),
		Status:      form.Status,
		Position:    form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: s.GetSetSessionValue(r),
		},
	}
	_, err := s.st.UpdateDepartment(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update designatio data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewDepartment(w http.ResponseWriter, r *http.Request) {
	logger.Info("view department form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetDepartmentBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get department " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := DepartmentData{
		Form: DepartmentForm{
			ID:          id,
			Name:        res.Name,
			Description: res.Description,
			Status:      res.Status,
			Position:    res.Position,
			CreatedAt:   res.CreatedAt,
			CreatedBy:   res.CreatedBy,
			UpdatedAt:   res.UpdatedAt,
			UpdatedBy:   res.UpdatedBy,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDepartmentStatus(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update department status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetDepartmentBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get department info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateDepartmentStatus(r.Context(), storage.Department{
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
		_, err := s.st.UpdateDepartmentStatus(r.Context(), storage.Department{
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
	http.Redirect(w, r, "/admin/"+departmentListPath, http.StatusSeeOther)
}
