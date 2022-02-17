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

type UserRoleData struct {
	CSRFField   template.HTML
	Form        UserRoleForm
	FormErrors  map[string]string
	Data        []UserRoleForm
	FormAction  string
	FormMessage map[string]string
}

func (s UserRoleForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name is required"),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Must be alphabet. No digit or special character is allowed"),
			validation.By(checkDuplicateUserRole(srv, s.Name, id)),
		),
		validation.Field(&s.Position,
			validation.Required.Error("The Position is required"),
			validation.By(checkUserRolePosition(srv, s.Position, id)),
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

func (s *Server) submitUserRole(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit userRole")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form UserRoleForm
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
		data := UserRoleData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	_, err = s.st.CreateUserRole(r.Context(), storage.UserRole{
		Name:        trim(form.Name),
		Description: trim(form.Description),
		Status:      form.Status,
		Position:    form.Position,
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

func (s *Server) userRoleList(w http.ResponseWriter, r *http.Request) {
	logger.Info("userRole list")
	tmpl := s.lookupTemplate("user-role.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	dptlist := s.usrRoleList(r, w, false)
	data := UserRoleData{
		Data: dptlist,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) deleteUserRole(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete userRole")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteUserRole(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete userRole" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateUserRole(w http.ResponseWriter, r *http.Request) {
	logger.Info("update userRole")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form UserRoleForm
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
		data := UserRoleData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.UserRole{
		ID:          id,
		Name:        trim(form.Name),
		Description: trim(form.Description),
		Status:      form.Status,
		Position:    form.Position,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "123",
		},
	}
	_, err := s.st.UpdateUserRole(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update designatio data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewUserRole(w http.ResponseWriter, r *http.Request) {
	logger.Info("view userRole form")
	params := mux.Vars(r)
	id := params["id"]
	res, err := s.st.GetUserRoleBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get userRole " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	data := UserRoleData{
		Form: UserRoleForm{
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

func (s *Server) updateUserRoleStatus(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update userRole status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetUserRoleBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get userRole info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateUserRoleStatus(r.Context(), storage.UserRole{
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
		_, err := s.st.UpdateUserRoleStatus(r.Context(), storage.UserRole{
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
	http.Redirect(w, r, "/admin/"+userRoleListPath, http.StatusSeeOther)
}
