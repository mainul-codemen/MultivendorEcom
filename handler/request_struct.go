package handler

import (
	"database/sql"
	"time"
)

type UserForm struct {
	ID                      string
	DesignationID           string
	UserRole                int16
	EmployeeRole            int16
	DepartmentID            int16
	HubID                   int32
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
	DateOfBirth             time.Time
	JoinDate                time.Time
	Gender                  int16
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
