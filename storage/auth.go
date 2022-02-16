package storage

import (
	"database/sql"
	"time"
)

type Users struct {
	ID                      string         `db:"id"`
	DesignationID           string         `db:"designation_id"`
	DesignationName         sql.NullString `db:"designation_name,omitempty"`
	UserRole                int16          `db:"user_role"`
	EmployeeRole            int16          `db:"employee_role"`
	VerifiedBy              string         `db:"verified_by"`
	JoinBy                  string         `db:"join_by"`
	CountryID               string         `db:"country_id"`
	CountryName             sql.NullString `db:"country_name,omitempty"`
	DistrictID              string         `db:"district_id"`
	DistrictName            sql.NullString `db:"district_name,omitempty"`
	StationID               string         `db:"station_id"`
	StationName             sql.NullString `db:"station_name,omitempty"`
	Status                  int16          `db:"status"`
	UserName                string         `db:"user_name"`
	FirstName               string         `db:"first_name"`
	LastName                string         `db:"last_name"`
	Email                   string         `db:"email"`
	EmailVerifiedAt         time.Time      `db:"email_verified_at"`
	Password                string         `db:"password"`
	Phone1                  string         `db:"phone_1"`
	Phone2                  string         `db:"phone_2"`
	PhoneNumberVerifiedAt   time.Time      `db:"phone_number_verified_at"`
	PhoneNumberVerifiedCode string         `db:"phone_number_verified_code"`
	DateOfBirth             time.Time      `db:"date_of_birth"`
	Gender                  int16          `db:"gender"`
	FBID                    string         `db:"fb_id"`
	Photo                   string         `db:"photo"`
	NIDFrontPhoto           string         `db:"nid_front_photo"`
	NIDBackPhoto            string         `db:"nid_back_photo"`
	NIDNumber               string         `db:"nid_number"`
	CVPDF                   string         `db:"cv_pdf"`
	PresentAddress          string         `db:"present_address"`
	PermanentAddress        string         `db:"permanent_address"`
	Reference               string         `db:"reference"`
	RememberToken           string         `db:"remember_token"`
	CRUDTimeDate
}
