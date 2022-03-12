package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

const AlrEx = " already exists"
const PosEx = "position %d already exists"
const PhnEx = "phone number %s already exists"

func (s *Server) countryList(r *http.Request, w http.ResponseWriter, sts bool) []CountryForm {
	cntryList, err := s.st.GetAllCountry(r.Context(), sts)
	if err != nil {
		logger.Error("error while get country : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	cntryListForm := make([]CountryForm, 0)
	for _, item := range cntryList {
		cntryData := CountryForm{
			ID:       item.ID,
			Name:     item.Name,
			Status:   item.Status,
			Position: item.Position,
		}
		cntryListForm = append(cntryListForm, cntryData)
	}
	return cntryListForm
}

func (s *Server) strToDisFormByID(w http.ResponseWriter, r *http.Request, id string) DistrictForm {
	res, err := s.st.GetDistrictBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get district " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Data := DistrictForm{
		ID:          res.ID,
		Name:        res.Name,
		CountryID:   res.CountryID,
		CountryName: res.CountryName.String,
		Status:      res.Status,
		Position:    res.Position,
		CreatedBy:   res.CreatedBy,
	}
	return Data
}

func (s *Server) stationList(r *http.Request, w http.ResponseWriter, sts bool) []StationForm {
	stnList, err := s.st.GetStationList(r.Context(), sts)
	if err != nil {
		logger.Error("error while get station : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	stnListForm := make([]StationForm, 0)
	for _, item := range stnList {
		cntryData := StationForm{
			ID:           item.ID,
			Name:         item.Name,
			Status:       item.Status,
			Position:     item.Position,
			DistrictName: item.DistrictName.String,
			CreatedBy:    item.CreatedBy,
		}
		stnListForm = append(stnListForm, cntryData)
	}
	return stnListForm
}

func (s *Server) districtList(r *http.Request, w http.ResponseWriter, sts bool) []DistrictForm {
	dist, err := s.st.GetDistrictList(r.Context(), sts)
	if err != nil {
		logger.Error("error while get district : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	disdata := make([]DistrictForm, 0)
	for _, item := range dist {
		disApnd := DistrictForm{
			ID:       item.ID,
			Name:     item.Name,
			Status:   item.Status,
			Position: item.Position,
		}
		disdata = append(disdata, disApnd)
	}
	return disdata
}

func (s *Server) hubList(r *http.Request, w http.ResponseWriter, sts bool) []HubForm {
	hubList, err := s.st.GetHubList(r.Context(), sts)
	if err != nil {
		logger.Error("error while get hub : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	hbdata := make([]HubForm, 0)
	for _, item := range hubList {
		disApnd := HubForm{
			ID:           item.ID,
			Name:         item.HubName,
			CountryID:    item.CountryID,
			CountryName:  item.CountryName.String,
			DistrictID:   item.DistrictID,
			DistrictName: item.DistrictName.String,
			StationID:    item.StationID,
			StationName:  item.StationName.String,
			HubPhone1:    item.HubPhone1,
			HubPhone2:    item.HubPhone2,
			HubEmail:     item.HubEmail,
			HubAddress:   item.HubAddress,
			Position:     item.Position,
			CreatedAt:    item.CreatedAt,
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt,
			UpdatedBy:    item.UpdatedBy,
		}
		hbdata = append(hbdata, disApnd)
	}
	return hbdata
}

func (s *Server) usrList(r *http.Request, w http.ResponseWriter, sts bool) []UserForm {
	usrList, err := s.st.GetUserList(r.Context(), sts)
	if err != nil {
		logger.Error("error while get user : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	usdata := make([]UserForm, 0)
	for _, item := range usrList {
		usApnd := UserForm{
			ID:                      item.ID,
			DesignationID:           item.DesignationID,
			DesignationName:         item.DesignationName.String,
			UserRole:                item.UserRole,
			EmployeeRole:            item.EmployeeRole,
			VerifiedBy:              item.VerifiedBy,
			JoinBy:                  item.JoinBy,
			CountryID:               item.CountryID,
			CountryName:             item.CountryName.String,
			DistrictID:              item.DistrictID,
			DistrictName:            item.DistrictName.String,
			StationID:               item.StationID,
			StationName:             item.StationName.String,
			Status:                  item.Status,
			UserName:                item.UserName,
			FirstName:               item.FirstName,
			LastName:                item.LastName,
			Email:                   item.Email,
			EmailVerifiedAt:         item.EmailVerifiedAt,
			Password:                item.Password,
			Phone1:                  item.Phone1,
			Phone2:                  item.Phone2,
			PhoneNumberVerifiedAt:   item.PhoneNumberVerifiedAt,
			PhoneNumberVerifiedCode: item.PhoneNumberVerifiedCode,
			EmailVerifiedCode:       item.EmailVerifiedCode,
			ISOTPVerified:           item.ISOTPVerified,
			ISEmailVerified:         item.ISEmailVerified,
			DateOfBirthT:            item.DateOfBirth,
			Gender:                  item.Gender,
			FBID:                    item.FBID,
			Photo:                   item.Photo,
			NIDFrontPhoto:           item.NIDFrontPhoto,
			NIDBackPhoto:            item.NIDBackPhoto,
			NIDNumber:               item.NIDNumber,
			CVPDF:                   item.CVPDF,
			PresentAddress:          item.PresentAddress,
			PermanentAddress:        item.PermanentAddress,
			Reference:               item.Reference,
			RememberToken:           item.RememberToken,
			CreatedAt:               item.CreatedAt,
			CreatedBy:               item.CreatedBy,
			UpdatedAt:               item.UpdatedAt,
			UpdatedBy:               item.UpdatedBy,
		}
		usdata = append(usdata, usApnd)
	}
	return usdata
}

func (s *Server) desList(r *http.Request, w http.ResponseWriter, sts bool) []DesignationForm {
	desList, err := s.st.GetDesignation(r.Context(), sts)
	if err != nil {
		logger.Error("error while get designation : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	desListForm := make([]DesignationForm, 0)
	for _, item := range desList {
		desData := DesignationForm{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			Position:    item.Position,
			CreatedAt:   item.CreatedAt,
			CreatedBy:   item.CreatedBy,
			UpdatedAt:   item.UpdatedAt,
			UpdatedBy:   item.UpdatedBy,
		}
		desListForm = append(desListForm, desData)
	}
	return desListForm
}

func (s *Server) dptList(r *http.Request, w http.ResponseWriter, sts bool) []DepartmentForm {
	desList, err := s.st.GetDepartment(r.Context(), sts)
	if err != nil {
		logger.Error("error while get designation : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	desListForm := make([]DepartmentForm, 0)
	for _, item := range desList {
		desData := DepartmentForm{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			Position:    item.Position,
			CreatedAt:   item.CreatedAt,
			CreatedBy:   item.CreatedBy,
			UpdatedAt:   item.UpdatedAt,
			UpdatedBy:   item.UpdatedBy,
		}
		desListForm = append(desListForm, desData)
	}
	return desListForm
}

func (s *Server) usrRoleList(r *http.Request, w http.ResponseWriter, sts bool) []UserRoleForm {
	desList, err := s.st.GetUserRole(r.Context(), sts)
	if err != nil {
		logger.Error("error while get user role : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	desListForm := make([]UserRoleForm, 0)
	for _, item := range desList {
		desData := UserRoleForm{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			Position:    item.Position,
			CreatedAt:   item.CreatedAt,
			CreatedBy:   item.CreatedBy,
			UpdatedAt:   item.UpdatedAt,
			UpdatedBy:   item.UpdatedBy,
		}
		desListForm = append(desListForm, desData)
	}
	return desListForm
}
func (s *Server) grdList(r *http.Request, w http.ResponseWriter, sts bool) []GradeForm {
	desList, err := s.st.GetGrade(r.Context(), sts)
	if err != nil {
		logger.Error("error while get designation : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	desListForm := make([]GradeForm, 0)
	for _, item := range desList {
		desData := GradeForm{
			ID:             item.ID,
			Name:           item.Name,
			Description:    item.Description,
			BasicSalary:    item.BasicSalary,
			LunchAllowance: item.LunchAllowance,
			RentAllowance:  item.RentAllowance,
			AbsentPenalty:  item.AbsentPenalty,
			TotalSalary:    item.TotalSalary,
			Transportation: item.Transportation,
			Status:         item.Status,
			Position:       item.Position,
			CreatedAt:      item.CreatedAt,
			CreatedBy:      item.CreatedBy,
			UpdatedAt:      item.UpdatedAt,
			UpdatedBy:      item.UpdatedBy,
		}
		desListForm = append(desListForm, desData)
	}
	return desListForm
}

func checkLogin(s *Server, emailorUsername string, pass string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetUserInfoBy(context.Background(), emailorUsername)
		if resp != nil {
			err := bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(pass))
			if err != nil {
				return fmt.Errorf(" Password doesn't match")
			}
		} else {
			return fmt.Errorf(" Please Enter valid credential")
		}
		return nil
	}
}

func formTemplate(s *Server, w http.ResponseWriter, r *http.Request, tmp string) {
	tmpl := s.lookupTemplate(tmp)
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

}

func (s *Server) accountsList(r *http.Request, w http.ResponseWriter, sts bool) []AccountsForm {
	actList, err := s.st.GetAccounts(r.Context(), sts)
	if err != nil {
		logger.Error("error while get accounts : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	actListForm := make([]AccountsForm, 0)
	for _, item := range actList {
		actData := AccountsForm{
			ID:                   item.ID,
			AccountVisualization: item.AccountVisualization,
			AccountNumber:        item.AccountNumber,
			AccountName:          item.AccountName,
			Amount:               item.Amount,
			Status:               item.Status,
		}
		actListForm = append(actListForm, actData)
	}
	return actListForm
}

func (*Server) storageToTransForm(transactionTypesList []storage.TransactionTypes) []TransactionTypesForm {
	transactionTypesListForm := make([]TransactionTypesForm, 0)
	for _, item := range transactionTypesList {
		transactionTypesData := TransactionTypesForm{
			ID:                   item.ID,
			TransactionTypesName: item.TransactionTypesName,
			Status:               item.Status,
			CreatedAt:            item.CreatedAt,
			CreatedBy:            item.CreatedBy,
			UpdatedAt:            item.UpdatedAt,
			UpdatedBy:            item.UpdatedBy,
		}
		transactionTypesListForm = append(transactionTypesListForm, transactionTypesData)
	}
	return transactionTypesListForm
}

func (*Server) storageToTranSrcForm(transactionSourceList []storage.TransactionSource) []TransactionSourceForm {
	transactionSourceListForm := make([]TransactionSourceForm, 0)
	for _, item := range transactionSourceList {
		transactionSourceData := TransactionSourceForm{
			ID:                    item.ID,
			TransactionSourceName: item.TransactionSourceName,
			Status:                item.Status,
			CreatedAt:             item.CreatedAt,
			CreatedBy:             item.CreatedBy,
			UpdatedAt:             item.UpdatedAt,
			UpdatedBy:             item.UpdatedBy,
		}
		transactionSourceListForm = append(transactionSourceListForm, transactionSourceData)
	}
	return transactionSourceListForm
}
func (s *Server) accountsTransactionList(r *http.Request, w http.ResponseWriter, sts bool) []AccountsTransactionForm {
	actList, err := s.st.GetAccountsTransaction(r.Context(), sts)
	if err != nil {
		logger.Error("error while get accounts : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	actListForm := make([]AccountsTransactionForm, 0)
	for _, item := range actList {
		actData := AccountsTransactionForm{
			ID:                      item.ID,
			FromAccountID:           item.FromAccountID,
			FromAccountName:         item.FromAccountName.String,
			ToAccountID:             item.ToAccountID,
			ToAccountName:           item.ToAccountName.String,
			UserID:                  item.UserID,
			TransactionAmount:       item.TransactionAmount,
			TransactionType:         item.TransactionType,
			TransactionTypeName:     item.TransactionTypeName.String,
			TransactionSource:       item.TransactionSource,
			TransactionSourceName:   item.TransactionSourceName.String,
			Reference:               item.Reference,
			Note:                    item.Note,
			Status:                  item.Status,
			FromAcntPreviousBalance: item.FromAcntPreviousBalance,
			FromAcntCurrentBalance:  item.FromAcntCurrentBalance,
			ToAcntPreviousBalance:   item.ToAcntPreviousBalance,
			ToAcntCurrentBalance:    item.ToAcntCurrentBalance,
			CreatedAt:               item.CreatedAt,
			CreatedBy:               item.CreatedBy,
			UpdatedAt:               item.UpdatedAt,
			UpdatedBy:               item.UpdatedBy,
		}
		actListForm = append(actListForm, actData)
	}
	return actListForm
}

func (s *Server) transactionTypesList(r *http.Request, w http.ResponseWriter, sts bool) []TransactionTypesForm {
	actList, err := s.st.GetTransactionTypes(r.Context(), sts)
	if err != nil {
		logger.Error("error while get transaction types : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	ttListForm := make([]TransactionTypesForm, 0)
	for _, item := range actList {
		ttData := TransactionTypesForm{
			ID:                   item.ID,
			TransactionTypesName: item.TransactionTypesName,
			Status:               item.Status,
			CreatedAt:            item.CreatedAt,
			CreatedBy:            item.CreatedBy,
		}
		ttListForm = append(ttListForm, ttData)
	}
	return ttListForm
}

func (s *Server) transactionSourceList(r *http.Request, w http.ResponseWriter, sts bool) []TransactionSourceForm {
	actList, err := s.st.GetTransactionSource(r.Context(), sts)
	if err != nil {
		logger.Error("error while get transaction source : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	ttListForm := make([]TransactionSourceForm, 0)
	for _, item := range actList {
		ttData := TransactionSourceForm{
			ID:                    item.ID,
			TransactionSourceName: item.TransactionSourceName,
			Status:                item.Status,
			CreatedAt:             item.CreatedAt,
			CreatedBy:             item.CreatedBy,
		}
		ttListForm = append(ttListForm, ttData)
	}
	return ttListForm
}
