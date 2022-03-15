package handler

import (
	"database/sql"
	"time"
)

type UserForm struct {
	ID                      string
	DesignationID           string
	DesignationName         string
	UserRole                string
	EmployeeRole            string
	DepartmentID            int16
	HubID                   string
	HubName                 string
	VerifiedBy              string
	JoinBy                  string
	CountryID               string
	CountryName             string
	DistrictID              string
	DistrictName            string
	StationID               string
	StationName             string
	Status                  int16
	UserName                string
	FirstName               string
	LastName                string
	Email                   string
	EmailVerifiedAt         time.Time
	Password                string
	Phone1                  string
	Phone2                  string
	PhoneNumberVerifiedAt   time.Time
	PhoneNumberVerifiedCode string
	EmailVerifiedCode       string
	ISOTPVerified           bool
	ISEmailVerified         bool
	DateOfBirth             string
	JoinDate                string
	DateOfBirthT            time.Time
	JoinDateT               time.Time
	Gender                  int16
	GradeID                 string
	FBID                    string
	Photo                   string
	NIDFrontPhoto           string
	NIDBackPhoto            string
	NIDNumber               string
	CVPDF                   string
	PresentAddress          string
	PermanentAddress        string
	Reference               string
	RememberToken           string
	CreatedAt               time.Time
	CreatedBy               string
	UpdatedAt               time.Time
	UpdatedBy               string
	DeletedAt               sql.NullTime
}

type DesignationForm struct {
	ID          string
	Name        string
	Description string
	Status      int32
	Position    int32
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   sql.NullTime
}
type DepartmentForm struct {
	ID          string
	Name        string
	Description string
	Status      int32
	Position    int32
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   sql.NullTime
}
type GradeForm struct {
	ID             string
	Name           string
	Description    string
	BasicSalary    float64
	LunchAllowance float64
	Transportation float64
	RentAllowance  float64
	AbsentPenalty  float64
	TotalSalary    float64
	Status         int32
	Position       int32
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	DeletedAt      sql.NullTime
}
type UserRoleForm struct {
	ID          string
	Name        string
	Description string
	Status      int32
	Position    int32
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   sql.NullTime
}
type DistrictForm struct {
	ID          string
	Name        string
	Description string
	CountryID   string
	CountryName string
	Status      int16
	Position    int32
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   sql.NullTime
}

type StationForm struct {
	ID           string
	Name         string
	DistrictID   string
	DistrictName string
	Status       int16
	Position     int32
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    sql.NullTime
}

type HubForm struct {
	ID           string
	Name         string
	CountryID    string
	CountryName  string
	DistrictID   string
	DistrictName string
	StationID    string
	StationName  string
	HubPhone1    string
	HubPhone2    string
	HubEmail     string
	HubAddress   string
	Status       int16
	Position     int32
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    sql.NullTime
}

type BranchForm struct {
	ID            string
	Name          string
	CountryID     string
	CountryName   string
	DistrictID    string
	DistrictName  string
	StationID     string
	StationName   string
	BranchName    string
	BranchPhone1  string
	BranchPhone2  string
	BranchEmail   string
	BranchAddress string
	BranchStatus  int16
	Position      int32
	CreatedAt     time.Time
	CreatedBy     string
	UpdatedAt     time.Time
	UpdatedBy     string
	DeletedAt     sql.NullTime
}

type DeliveryCompanyForm struct {
	ID             string
	CompanyName    string
	CountryID      string
	CountryName    string
	DistrictID     string
	DistrictName   string
	StationID      string
	StationName    string
	Phone          string
	Email          string
	CompanyAddress string
	CompanyStatus  int16
	Position       int32
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	DeletedAt      sql.NullTime
}

type DeliveryChargeForm struct {
	ID                   string
	CompanyName          string
	CountryID            string
	CountryName          string
	DistrictID           string
	DistrictName         string
	StationID            string
	StationName          string
	DeliveryChargeStatus int16
	DeliveryCharge       float64
	WeightMin            float64
	WeightMax            float64
	CreatedAt            time.Time
	CreatedBy            string
	UpdatedAt            time.Time
	UpdatedBy            string
	DeletedAt            sql.NullTime
}

type CountryForm struct {
	ID          string
	Name        string
	Description string
	Status      int16
	Position    int32
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   sql.NullTime
}
type AccountsForm struct {
	ID                   string
	AccountVisualization string
	AccountName          string
	AccountNumber        string
	Amount               float64
	Status               int32
	CreatedAt            time.Time
	CreatedBy            string
	UpdatedAt            time.Time
	UpdatedBy            string
	DeletedAt            sql.NullTime
}

type AccountsTransactionForm struct {
	ID                      string
	FromAccountID           string
	FromAccountName         string
	ToAccountID             string
	ToAccountName           string
	UserID                  string
	TransactionAmount       float64
	FromAcntPreviousBalance float64
	FromAcntCurrentBalance  float64
	ToAcntPreviousBalance   float64
	ToAcntCurrentBalance    float64
	TransactionType         string
	TransactionTypeName     string
	TransactionSource       string
	TransactionSourceName   string
	Reference               string
	Note                    string
	Status                  int32
	AcceptedAt              time.Time
	AcceptedBy              string
	CreatedAt               time.Time
	CreatedBy               string
	UpdatedAt               time.Time
	UpdatedBy               string
	DeletedAt               sql.NullTime
}

type TransactionTypesForm struct {
	ID                   string
	TransactionTypesName string
	Status               int32
	CreatedAt            time.Time
	CreatedBy            string
	UpdatedAt            time.Time
	UpdatedBy            string
	DeletedAt            sql.NullTime
}
type TransactionSourceForm struct {
	ID                    string
	TransactionSourceName string
	Status                int32
	CreatedAt             time.Time
	CreatedBy             string
	UpdatedAt             time.Time
	UpdatedBy             string
	DeletedAt             sql.NullTime
}

type IncomeTaxForm struct {
	ID               string
	AccountID        string
	AccountNumber    sql.NullString
	AccountName      sql.NullString
	TaxReceiptNumber string
	Status           int32
	IncomeTaxDate    time.Time
	TaxAmount        float64
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
	DeletedAt        sql.NullTime
}

type IncomeForm struct {
	ID            string
	Title         string
	AccountID     string
	Note          string
	AccountNumber sql.NullString
	AccountName   sql.NullString
	IncomeAmount  float64
	IncomeDate    string
	Status        int32
	CreatedAt     time.Time
	CreatedBy     string
	UpdatedAt     time.Time
	UpdatedBy     string
	DeletedAt     sql.NullTime
}
