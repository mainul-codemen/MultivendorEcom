package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertusr = `
INSERT INTO users(
	designation_id,
	user_role,
	employee_role,
	verified_by,
	join_by,
	country_id,
	district_id,
	station_id,
	status,
	user_name,
	first_name,
	last_name,
	email,
	email_verified_at,
	password,
	phone_1,
	phone_2,
	phone_number_verified_at,
	phone_number_verified_code,
	date_of_birth,
	gender,
	fb_id,
	photo,
	nid_front_photo,
	nid_back_photo,
	nid_number,
	cv_pdf,
	present_address,
	permanent_address,
	reference,
	remember_token,
	created_by,
	updated_by
) VALUES (
	:designation_id,
	:user_role,
	:employee_role,
	:verified_by,
	:join_by,
	:country_id,
	:district_id,
	:station_id,
	:status,
	:user_name,
	:first_name,
	:last_name,
	:email,
	:email_verified_at,
	:password,
	:phone_1,
	:phone_2,
	:phone_number_verified_at,
	:phone_number_verified_code,
	:date_of_birth,
	:gender,
	:fb_id,
	:photo,
	:nid_front_photo,
	:nid_back_photo,
	:nid_number,
	:cv_pdf,
	:present_address,
	:permanent_address,
	:reference,
	:remember_token,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) RegisterUser(con context.Context, des storage.Users) (string, error) {
	logger.Info("create users db")
	stmt, err := s.db.PrepareNamed(insertusr)
	if err != nil {
		logger.Error(ewpq + err.Error())
		return "", err
	}

	var id string
	if err := stmt.Get(&id, des); err != nil {
		logger.Error(ewmq + err.Error())
		return "", err
	}
	return id, nil
}

const usrLisQry = `
SELECT 
	designation.name AS designation_name,
	u.user_role,
	u.employee_role,
	u.verified_by,
	u.join_by,
	country.name AS country_name,
	district.name AS district_name,
	station.name AS station_name,
	u.status,
	u.user_name,
	u.first_name,
	u.last_name,
	u.email,
	u.email_verified_at,
	u.password,
	u.phone_1,
	u.phone_2,
	u.phone_number_verified_at,
	u.phone_number_verified_code,
	u.date_of_birth,
	u.gender,
	u.fb_id,
	u.photo,
	u.nid_front_photo,
	u.nid_back_photo,
	u.nid_number,
	u.cv_pdf,
	u.present_address,
	u.permanent_address,
	u.reference,
	u.remember_token,
	u.created_by,
	u.updated_by
FROM users AS u
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE u.deleted_at IS NULL
`

func (s *Storage) GetUserList(ctx context.Context, sts bool) ([]storage.Users, error) {
	logger.Info("get all users")
	ul := make([]storage.Users, 0)
	if err := s.db.Select(&ul, usrLisQry); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return ul, nil
}

const gusrq = `
SELECT 
	designation.name AS designation_name,
	u.user_role,
	u.employee_role,
	u.verified_by,
	u.join_by,
	country.name AS country_name,
	district.name AS district_name,
	station.name AS station_name,
	u.status,
	u.user_name,
	u.first_name,
	u.last_name,
	u.email,
	u.email_verified_at,
	u.password,
	u.phone_1,
	u.phone_2,
	u.phone_number_verified_at,
	u.phone_number_verified_code,
	u.date_of_birth,
	u.gender,
	u.fb_id,
	u.photo,
	u.nid_front_photo,
	u.nid_back_photo,
	u.nid_number,
	u.cv_pdf,
	u.present_address,
	u.permanent_address,
	u.reference,
	u.remember_token,
	u.created_by,
	u.updated_by
FROM users AS u
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE (u.id = $1 OR u.user_name=$1 OR u.email=$1 OR u.phone_1=$1 OR u.phone_2=$1) AND  dc.deleted_at IS NULL
`

func (s *Storage) GetUserInfoBy(ctx context.Context, idname string) (*storage.Users, error) {
	var usr storage.Users
	if err := s.db.Get(&usr, gusrq, idname); err != nil {
		logger.Error("error while get users data. " + err.Error())
		return nil, fmt.Errorf("executing users details: %w", err)
	}
	return &usr, nil
}

const updateUsr = `
UPDATE users SET
	designation_id = :designation_id,
	user_role = :user_role,
	employee_role = :employee_role,
	verified_by = :verified_by,
	join_by = :join_by,
	country_id = :country_id,
	district_id = :district_id,
	station_id = :station_id,
	status = :status,
	user_name = :user_name,
	first_name = :first_name,
	last_name = :last_name,
	email = :email,
	email_verified_at = :email_verified_at,
	password = :password,
	phone_1 = :phone_1,
	phone_2 = :phone_2,
	phone_number_verified_at = :phone_number_verified_at,
	phone_number_verified_code = :phone_number_verified_code,
	date_of_birth = :date_of_birth,
	gender = :gender,
	fb_id = :fb_id,
	photo = :photo,
	nid_front_photo =:nid_front_photo,
	nid_back_photo = :nid_back_photo,
	nid_number = :nid_number,
	cv_pdf = :cv_pdf,
	present_address = :present_address,
	permanent_address = :permanent_address,
	reference = :reference,
	remember_token = :remember_token,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateUserInfo(ctx context.Context, usr storage.Users) (*storage.Users, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateUsr)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&usr, usr); err != nil {
		return nil, fmt.Errorf("executing delivry charge update: %w", err)
	}

	return &usr, nil
}

const upUsts = `
	UPDATE 
		users 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateUserStatus(ctx context.Context, usr storage.Users) (*storage.Users, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, upUsts)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&usr, usr); err != nil {
		return nil, fmt.Errorf("executing users status: %w", err)
	}

	return &usr, nil
}

const delUsers = `
UPDATE
	users
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteUsers(ctx context.Context, userid string, duid string) error {
	logger.Info("delete users")
	_, err := s.db.Exec(delUsers, duid, userid)
	if err != nil {
		logger.Error("delete users")
		return err
	}
	return nil
}
