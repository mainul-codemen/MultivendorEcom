package handler

import (
	"database/sql"
	"time"
)

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
	HubName      string
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
