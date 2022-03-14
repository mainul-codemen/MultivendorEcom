package storage

import (
	"database/sql"
	"time"
)

const (
	ewpq = "error while prepared query"
)

type (
	Message struct {
		Status  bool
		Message string
	}

	AccountsTransaction struct {
		ID                      string         `db:"id"`
		FromAccountID           string         `db:"from_account_id"`
		FromAccountName         sql.NullString `db:"from_account_name"`
		ToAccountID             string         `db:"to_account_id"`
		ToAccountName           sql.NullString `db:"to_account_name"`
		UserID                  string         `db:"user_id"`
		TransactionAmount       float64        `db:"transaction_amount,omitempty"`
		FromAcntPreviousBalance float64        `db:"from_acnt_previous_balance,omitempty"`
		FromAcntCurrentBalance  float64        `db:"from_acnt_current_balance,omitempty"`
		ToAcntPreviousBalance   float64        `db:"to_acnt_previous_balance,omitempty"`
		ToAcntCurrentBalance    float64        `db:"to_acnt_current_balance,omitempty"`
		TransactionType         string         `db:"transaction_type_id,omitempty"`
		TransactionTypeName     sql.NullString `db:"transaction_type_name,omitempty"`
		TransactionSource       string         `db:"transaction_source_id,omitempty"`
		TransactionSourceName   sql.NullString `db:"transaction_source_name,omitempty"`
		Reference               string         `db:"reference"`
		Note                    string         `db:"note"`
		Status                  int32          `db:"status"`
		AcceptedAt              time.Time      `db:"accepted_at,omitempty"`
		AcceptedBy              string         `db:"accepted_by"`
		CRUDTimeDate
	}
	Accounts struct {
		ID                   string  `db:"id"`
		AccountVisualization string  `db:"account_visualization"`
		AccountName          string  `db:"account_name"`
		AccountNumber        string  `db:"account_number"`
		Amount               float64 `db:"amount"`
		Status               int32   `db:"status"`
		CRUDTimeDate
	}

	IncomeTax struct {
		ID               string    `db:"id"`
		AccountID        string    `db:"account_id"`
		TaxReceiptNumber string    `db:"tax_receipt_number"`
		Status           int32     `db:"status"`
		IncomeTaxDate    time.Time `db:"income_tax_date"`
		TaxAmount        float64   `db:"tax_amount"`
		CRUDTimeDate
	}

	Income struct {
		ID               string  `db:"id"`
		Title            string  `db:"title"`
		TaxAmount        float64 `db:"tax_amount"`
		AccountID        string  `db:"account_id"`
		TaxReceiptNumber float64 `db:"tax_receipt_number"`
		IncomeTaxDate    float64 `db:"income_tax_date"`
		Status           int32   `db:"status"`
		CRUDTimeDate
	}

	Designation struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		Status      int32  `db:"status"`
		Position    int32  `db:"position"`
		CRUDTimeDate
	}

	Department struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		Status      int32  `db:"status"`
		Position    int32  `db:"position"`
		CRUDTimeDate
	}
	UserRole struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		Status      int32  `db:"status"`
		Position    int32  `db:"position"`
		CRUDTimeDate
	}
	Grade struct {
		ID             string  `db:"id"`
		Name           string  `db:"name"`
		Description    string  `db:"description"`
		BasicSalary    float64 `db:"basic_salary"`
		LunchAllowance float64 `db:"lunch_allowance"`
		Transportation float64 `db:"transportation"`
		RentAllowance  float64 `db:"rent_allowance"`
		AbsentPenalty  float64 `db:"absent_penalty"`
		TotalSalary    float64 `db:"total_salary"`
		Status         int32   `db:"status"`
		Position       int32   `db:"position"`
		CRUDTimeDate
	}
	Country struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		Status      int16  `db:"status"`
		Position    int32  `db:"position"`
		CRUDTimeDate
	}
	District struct {
		ID          string         `db:"id"`
		Name        string         `db:"name"`
		CountryID   string         `db:"country_id"`
		CountryName sql.NullString `db:"country_name,omitempty"`
		Status      int16          `db:"status"`
		Position    int32          `db:"position"`
		CRUDTimeDate
	}

	Station struct {
		ID           string         `db:"id"`
		Name         string         `db:"name"`
		DistrictID   string         `db:"district_id"`
		DistrictName sql.NullString `db:"district_name,omitempty"`
		Status       int16          `db:"status"`
		Position     int32          `db:"position"`
		CRUDTimeDate
	}
	Hub struct {
		ID           string         `db:"id"`
		HubName      string         `db:"hub_name"`
		CountryID    string         `db:"country_id"`
		CountryName  sql.NullString `db:"country_name,omitempty"`
		DistrictID   string         `db:"district_id"`
		DistrictName sql.NullString `db:"district_name,omitempty"`
		StationID    string         `db:"station_id"`
		StationName  sql.NullString `db:"station_name,omitempty"`
		HubPhone1    string         `db:"hub_phone_1"`
		HubPhone2    string         `db:"hub_phone_2"`
		HubEmail     string         `db:"hub_email"`
		HubAddress   string         `db:"hub_address"`
		Status       int16          `db:"hub_status"`
		Position     int32          `db:"position"`
		CRUDTimeDate
	}
	DeliveryCharge struct {
		ID             string         `db:"id"`
		CountryID      string         `db:"country_id"`
		CountryName    sql.NullString `db:"country_name,omitempty"`
		DistrictID     string         `db:"district_id"`
		DistrictName   sql.NullString `db:"district_name,omitempty"`
		StationID      string         `db:"station_id"`
		StationName    sql.NullString `db:"station_name,omitempty"`
		WeightMin      float64        `db:"weight_min"`
		WeightMax      float64        `db:"weight_max"`
		DeliveryCharge float64        `db:"delivery_charge"`
		Status         int16          `db:"dc_status"`
		CRUDTimeDate
	}
	Branch struct {
		ID            string         `db:"id"`
		BranchName    string         `db:"branch_name"`
		CountryID     string         `db:"country_id"`
		CountryName   sql.NullString `db:"country_name,omitempty"`
		DistrictID    string         `db:"district_id"`
		DistrictName  sql.NullString `db:"district_name,omitempty"`
		StationID     string         `db:"station_id"`
		StationName   sql.NullString `db:"station_name,omitempty"`
		BranchPhone1  string         `db:"branch_phone_1"`
		BranchPhone2  string         `db:"branch_phone_2"`
		BranchEmail   string         `db:"branch_email"`
		BranchAddress string         `db:"branch_address"`
		BranchStatus  int16          `db:"branch_status"`
		Position      int32          `db:"position"`
		CRUDTimeDate
	}

	DeliveryCompany struct {
		ID             string         `db:"id"`
		CompanyName    string         `db:"company_name"`
		CountryID      string         `db:"country_id"`
		CountryName    sql.NullString `db:"country_name,omitempty"`
		DistrictID     string         `db:"district_id"`
		DistrictName   sql.NullString `db:"district_name,omitempty"`
		StationID      string         `db:"station_id"`
		StationName    sql.NullString `db:"station_name,omitempty"`
		Phone          string         `db:"phone"`
		Email          string         `db:"email"`
		CompanyAddress string         `db:"company_address"`
		CompanyStatus  int16          `db:"company_status"`
		Position       int32          `db:"position"`
		CRUDTimeDate
	}
	PassResRequest struct {
		ID                string    `db:"id"`
		UserID            string    `db:"user_id"`
		Password          string    `db:"password"`
		Token             string    `db:"password_reset_token"`
		PasswordResetTime time.Time `db:"pass_reset_time"`
	}
	TransactionTypes struct {
		ID                   string `db:"id"`
		TransactionTypesName string `db:"transaction_type_name"`
		Status               int32  `db:"status"`
		CRUDTimeDate
	}
	TransactionSource struct {
		ID                    string `db:"id"`
		TransactionSourceName string `db:"transaction_source_name"`
		Status                int32  `db:"status"`
		CRUDTimeDate
	}

	CRUDTimeDate struct {
		CreatedAt time.Time      `db:"created_at,omitempty"`
		CreatedBy string         `db:"created_by"`
		UpdatedAt time.Time      `db:"updated_at,omitempty"`
		UpdatedBy string         `db:"updated_by,omitempty"`
		DeletedAt sql.NullTime   `db:"deleted_at,omitempty"`
		DeletedBy sql.NullString `db:"deleted_by,omitempty"`
	}
)
