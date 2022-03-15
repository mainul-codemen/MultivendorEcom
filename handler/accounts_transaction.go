package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type AccountsTransactionTempData struct {
	CSRFField             template.HTML
	Form                  AccountsTransactionForm
	FormErrors            map[string]string
	Data                  []AccountsTransactionForm
	FormAction            string
	FormMessage           map[string]string
	AccountsData          []AccountsForm
	TransactionTypeData   []TransactionTypesForm
	TransactionSourceData []TransactionSourceForm
}

func (s AccountsTransactionForm) Validate(srv *Server, r *http.Request, w http.ResponseWriter, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.FromAccountID,
			validation.Required.Error("From Account name is required"),
			validation.By(validateBalance(srv, w, r, s.FromAccountID, s.ToAccountID, s.TransactionAmount)),
		),
		validation.Field(&s.ToAccountID,
			validation.Required.Error("To Account Number is required"),
		),
		validation.Field(&s.TransactionAmount,
			validation.Required.Error("Transaciton Amount is required"),
		),
	)
}

func (s *Server) accountsTransactionFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("account transaction form")
	actdata := s.accountsList(r, w, true)
	ttdata := s.transactionTypesList(r, w, true)
	tsdata := s.transactionSourceList(r, w, true)
	data := AccountsTransactionTempData{
		CSRFField:             csrf.TemplateField(r),
		AccountsData:          actdata,
		TransactionTypeData:   ttdata,
		TransactionSourceData: tsdata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitAccountsTransactionHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("submit accounts transaction information")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	var form AccountsTransactionForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}

	if err := form.Validate(s, r, w, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := AccountsTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	resp, err := s.st.GetAccountsBy(context.Background(), trim(form.FromAccountID))
	if err != nil {
		logger.Error("error while get account data" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}
	fromAcntpreviousAmount := resp.Amount
	// deduct money from account : Current Amount
	fromAcntCurrentAmount := fromAcntpreviousAmount - form.TransactionAmount
	_, err = s.st.UpdateBalance(context.Background(), storage.Accounts{
		ID:     resp.ID,
		Amount: fromAcntCurrentAmount,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error("error while updating data " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}
	// add money to the to account
	resp2, err := s.st.GetAccountsBy(context.Background(), trim(form.ToAccountID))
	if err != nil {
		logger.Error("error while get account data" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}
	toAcntpreviousAmount := resp2.Amount
	toAcntCurrentBalance := toAcntpreviousAmount + form.TransactionAmount
	_, err = s.st.UpdateBalance(context.Background(), storage.Accounts{
		ID:     resp2.ID,
		Amount: toAcntCurrentBalance,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error("error while get account data" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}

	_, err = s.st.CreateAccountsTransaction(r.Context(), storage.AccountsTransaction{
		FromAccountID:     trim(form.FromAccountID),
		ToAccountID:       trim(form.ToAccountID),
		UserID:            uid,
		TransactionAmount: form.TransactionAmount,
		TransactionType:   form.TransactionType,
		TransactionSource: form.TransactionSource,
		Reference:         form.Reference,
		Note:              form.Note,
		Status:            1,
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

func (s *Server) accountsTransactionListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("account transaction list")
	tmpl := s.lookupTemplate("accounts-transaction.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	actdata := s.accountsTransactionList(r, w, false)
	data := AccountsTransactionTempData{
		Data: actdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) updateAccountTransactionFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update account transaction form")
	params := mux.Vars(r)
	id := params["id"]
	form := getTransactionByID(s, r, id, w)
	actdata := s.accountsList(r, w, true)
	ttdata := s.transactionTypesList(r, w, true)
	tsdata := s.transactionSourceList(r, w, true)
	data := AccountsTransactionTempData{
		CSRFField:             csrf.TemplateField(r),
		AccountsData:          actdata,
		TransactionTypeData:   ttdata,
		TransactionSourceData: tsdata,
		Form:                  form,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateAccountsTransactionHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update accounts transaction")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form AccountsTransactionForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logger.Error(deErr + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}
	if err := form.Validate(s, r, w, id); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := AccountsTransactionTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err := s.st.UpdateAccountsTransaction(r.Context(), storage.AccountsTransaction{
		ID:                id,
		FromAccountID:     trim(form.FromAccountID),
		ToAccountID:       trim(form.ToAccountID),
		UserID:            uid,
		TransactionAmount: form.TransactionAmount,
		TransactionType:   form.TransactionType,
		TransactionSource: form.TransactionSource,
		Reference:         form.Reference,
		Note:              form.Note,
		Status:            form.Status,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error("error while update accounts data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewAccountsTransactionHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view accounts transaction")
	params := mux.Vars(r)
	id := params["id"]
	data := getTransactionByID(s, r, id, w)
	json.NewEncoder(w).Encode(data)
}

func getTransactionByID(s *Server, r *http.Request, id string, w http.ResponseWriter) AccountsTransactionForm {
	res, err := s.st.GetAccountsTransactionBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get accounts " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	form := AccountsTransactionForm{
		ID:                    id,
		FromAccountID:         res.FromAccountID,
		FromAccountName:       res.FromAccountName.String,
		ToAccountID:           res.ToAccountID,
		ToAccountName:         res.ToAccountName.String,
		TransactionAmount:     res.TransactionAmount,
		TransactionSource:     res.TransactionSource,
		TransactionSourceName: res.TransactionSourceName.String,
		TransactionType:       res.TransactionType,
		TransactionTypeName:   res.TransactionTypeName.String,
		Reference:             res.Reference,
		Note:                  res.Note,
		UserID:                res.UserID,
		Status:                res.Status,
		CreatedAt:             res.CreatedAt,
		CreatedBy:             res.CreatedBy,
		UpdatedAt:             res.UpdatedAt,
		UpdatedBy:             res.UpdatedBy,
	}
	return form
}

func (s *Server) updateAccountsTransactionStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update accounts status")
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	uid, _ := s.GetSetSessionValue(r, w)
	res, err := s.st.GetAccountsTransactionBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get accounts info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateAccountsTransactionStatus(r.Context(), storage.AccountsTransaction{
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
		_, err := s.st.UpdateAccountsTransactionStatus(r.Context(), storage.AccountsTransaction{
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
	http.Redirect(w, r, "/admin/"+accountsTransactionListPath, http.StatusSeeOther)
}

func (s *Server) deleteAccountsTransactionHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete accounts")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteAccountsTransaction(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete accounts" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	json.NewEncoder(w).Encode(msg)
}

func validateBalance(s *Server, w http.ResponseWriter, r *http.Request, acnt, toantid string, amount float64) validation.RuleFunc {
	return func(value interface{}) error {
		// check equal bank: not possible to transfer
		if acnt == toantid {
			return fmt.Errorf(" Balance Transfer of Same account is not Possible")
		}
		resp, err := s.st.GetAccountsBy(context.Background(), trim(acnt))
		if err != nil {
			return fmt.Errorf(" Account is not valid")
		} else {
			// check enough money in the account
			stredAmt := resp.Amount
			if stredAmt < amount {
				return fmt.Errorf(" In This Accont You Don't have enough amount to transfer")
			} else {
				return nil
			}
		}
	}
}
