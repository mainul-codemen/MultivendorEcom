package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/MultivendorEcom/serviceutil/logger"
	lutil "github.com/MultivendorEcom/serviceutil/loginUtils"
	mail "github.com/MultivendorEcom/serviceutil/mail"
	"github.com/MultivendorEcom/serviceutil/otp"
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
	DepartmentData  []DepartmentForm
	UserRoleData    []UserRoleForm
	HubData         []HubForm
	GradeData       []GradeForm
	FormAction      string
}

func (s UserForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.UserName,
			validation.Required.Error(nameReq),
			validation.Length(3, 50).Error("Please insert name between 3 to 50"),
			validation.By(checkDuplicateUserName(srv, s.UserName)),
		),
		validation.Field(&s.Phone1,
			validation.Required.Error("The User phone is required"),
			validation.Length(3, 11).Error("Please insert phone between 3 to 11"),
			validation.Match(regexp.MustCompile("^[0-9_ ]*$")).Error("Must be digit. No alphabet is allowed."),
			validation.By(checkDuplicateUserPhone(srv, s.Phone1, id)),
		),
		validation.Field(&s.Password,
			validation.Required.Error("Password is Required"),
			validation.By(validatePassword(srv, s.Password)),
		),
		validation.Field(&s.CountryID,
			validation.Required.Error("Country Name is Required"),
		),
		validation.Field(&s.DesignationID,
			validation.Required.Error("Designation Name is Required"),
		),
		validation.Field(&s.DistrictID,
			validation.Required.Error("District Name is Required"),
		),
		validation.Field(&s.FirstName,
			validation.Required.Error("First Name is Required"),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("Last Name is Required"),
		),
		validation.Field(&s.DateOfBirth,
			validation.Required.Error("Date Of Birth is Required"),
		),
		validation.Field(&s.Email,
			validation.Required.Error("Email is Required"),
			validation.By(checkDuplicateUserEmail(srv, s.Email)),
		),
		validation.Field(&s.Gender,
			validation.Required.Error("Gender is Required"),
		),
		validation.Field(&s.Phone1,
			validation.Required.Error("Phone 1 is Required"),
		),
		validation.Field(&s.HubID,
			validation.Required.Error("Hub Name is Required"),
		),
		validation.Field(&s.JoinDate,
			validation.Required.Error("Join Date is Required"),
		),
		validation.Field(&s.StationID,
			validation.Required.Error("Station Name is Required"),
		),
		validation.Field(&s.NIDNumber,
			validation.Required.Error("NID Number is Required"),
		),
		validation.Field(&s.UserRole,
			validation.Required.Error("User Role is Required"),
		),
		validation.Field(&s.UserRole,
			validation.Required.Error("Employee Role is Required"),
		),
		validation.Field(&s.GradeID,
			validation.Required.Error("Grade is Required"),
		),
		validation.Field(&s.GradeID,
			validation.Required.Error("Grade is Required"),
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
	data := UsrTempData{
		CSRFField:       csrf.TemplateField(r),
		CountryData:     s.countryList(r, w, true),  // true = only active status
		DistrictData:    s.districtList(r, w, true), // true = only active status
		StationData:     s.stationList(r, w, true),  // true = only active status
		DesignationData: s.desList(r, w, true),      // true = only active status
		UserRoleData:    s.usrRoleList(r, w, true),  // true = only active status
		GradeData:       s.grdList(r, w, true),      // true = only active status
		HubData:         s.hubList(r, w, true),      // true = only active status
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
			CountryData:     s.countryList(r, w, true),  // true = only active status
			DistrictData:    s.districtList(r, w, true), // true = only active status
			StationData:     s.stationList(r, w, true),  // true = only active status
			DesignationData: s.desList(r, w, true),      // true = only active status
			UserRoleData:    s.usrRoleList(r, w, true),  // true = only active status
			GradeData:       s.grdList(r, w, true),      // true = only active status
			HubData:         s.hubList(r, w, true),      // true = only active status
			CSRFField:       csrf.TemplateField(r),
			FormErrors:      vErrs,
			Form:            form,
		}
		tmpl := s.lookupTemplate("user-create.html")
		if tmpl == nil {
			logger.Error(ult)
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
		if err := tmpl.Execute(w, data); err != nil {
			logger.Error(ewte + err.Error())
			http.Redirect(w, r, "/admin"+createUserPath, http.StatusSeeOther)
		}
		return
	}
	hp, err := lutil.HashPassword(trim(form.Password))
	if err != nil {
		log.Fatalln("Unable to hashed password")
	}
	// send mail
	emailCode := otp.SendEmailVerificationCode(form.Email)
	err = mail.SendingMail(trim(form.Email), emailCode)
	if err != nil {
		fmt.Println("User Data saved. But Email is not send for some reason.")
		emailCode = ""
	}
	_, err = s.st.RegisterUser(r.Context(), storage.Users{
		DesignationID:           form.DesignationID,
		CountryID:               form.CountryID,
		HubID:                   form.HubID,
		DistrictID:              form.DistrictID,
		StationID:               form.StationID,
		JoinBy:                  "124",
		EmployeeRole:            form.EmployeeRole,
		UserRole:                form.UserRole,
		VerifiedBy:              s.GetSetSessionValue(r),
		Status:                  1,
		GradeID:                 form.GradeID,
		UserName:                trim(form.UserName),
		FirstName:               trim(form.FirstName),
		LastName:                trim(form.LastName),
		Email:                   trim(form.Email),
		Password:                hp,
		Phone1:                  trim(form.Phone1),
		Phone2:                  trim(form.Phone2),
		PhoneNumberVerifiedCode: otp.SendVerificationCode(form.Phone1),
		PhoneNumberVerifiedAt:   time.Now().Add(time.Minute * 5),
		EmailVerifiedCode:       emailCode,
		EmailVerifiedAt:         time.Now().Add(time.Minute * 30),
		ISOTPVerified:           false,
		JoinDate:                form.JoinDateT,
		DateOfBirth:             form.DateOfBirthT,
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
		CRUDTimeDate:            storage.CRUDTimeDate{CreatedBy: s.GetSetSessionValue(r), UpdatedBy: "123"},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admin"+userListPath, http.StatusSeeOther)
}

func (s *Server) viewVerificationForm(w http.ResponseWriter, r *http.Request) {
	logger.Info("view verification form")
	params := mux.Vars(r)
	id := params["id"]
	vf := s.strgToFormUserVC(r, id, w)
	data := UsrTempData{
		CSRFField: csrf.TemplateField(r),
		Form:      vf,
	}
	json.NewEncoder(w).Encode(data)
}
func (s *Server) submitVerificationCode(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit verification code")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	vc := s.strgToFormUserVC(r, id, w)
	var form UserForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	// check valid otp
	vErrs := map[string]string{}
	if vc.PhoneNumberVerifiedCode != form.PhoneNumberVerifiedCode || form.PhoneNumberVerifiedCode == "" {
		vErrs["OTPErr"] = "OTP is Not valid"
		data := UsrTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       vc,
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	data := storage.Users{
		ID:                      id,
		PhoneNumberVerifiedCode: "",
		ISOTPVerified:           true,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "1234",
		},
	}
	_, err = s.st.VerifyPhoneNumber(r.Context(), data)
	if err != nil {
		logger.Error("error while verify phone ." + err.Error())
	}

	json.NewEncoder(w).Encode(msg)
}

func (s *Server) submitEmailVerificationCode(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit email verification code")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	vc := s.strgToFormUserVC(r, id, w)
	var form UserForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	// check valid code
	vErrs := map[string]string{}
	if vc.EmailVerifiedCode != form.EmailVerifiedCode || form.EmailVerifiedCode == "" {
		vErrs["EmailErr"] = "Verification Code is Not valid"
		data := UsrTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       vc,
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	data := storage.Users{
		ID:                id,
		EmailVerifiedCode: "",
		ISEmailVerified:   true,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "1234",
		},
	}
	_, err = s.st.VerifyEmail(r.Context(), data)
	if err != nil {
		logger.Error("error while verify phone ." + err.Error())
	}

	json.NewEncoder(w).Encode(msg)
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
				UpdatedBy: s.GetSetSessionValue(r),
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
				UpdatedBy: s.GetSetSessionValue(r),
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

func (s *Server) strgToFormUserVC(r *http.Request, id string, w http.ResponseWriter) UserForm {
	res, err := s.st.GetUserInfoBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get verify phone " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Form := UserForm{
		ID:                      res.ID,
		Phone1:                  res.Phone1,
		PhoneNumberVerifiedCode: res.PhoneNumberVerifiedCode,
		Email:                   res.Email,
		EmailVerifiedCode:       res.EmailVerifiedCode,
		PhoneNumberVerifiedAt:   res.PhoneNumberVerifiedAt,
		EmailVerifiedAt:         res.EmailVerifiedAt,
	}
	return Form
}
