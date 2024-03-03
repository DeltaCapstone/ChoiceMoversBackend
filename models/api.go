package models

import "github.com/jackc/pgx/v5/pgtype"

// ---------------------------
// EMPLOYEE
// ---------------------------

type CreateEmployeeParams struct {
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType string        `db:"employee_type" json:"employeeType"`
}

type CreateEmployeeRequest struct {
	UserName      string        `db:"username" json:"userName"`
	PasswordPlain string        `db:"password_plain" json:"passwordPlain"`
	FirstName     string        `db:"first_name" json:"firstName"`
	LastName      string        `db:"last_name" json:"lastName"`
	Email         string        `db:"email" json:"email"`
	PhonePrimary  pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther    []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType  string        `db:"employee_type" json:"employeeType"`
}

type EmployeeLoginRequest struct {
	UserName      string `db:"username" json:"userName"`
	PasswordPlain string `db:"password_plain" json:"passwordPlain"`
}

type GetEmployeeResponse struct {
	UserName     string        `db:"username" json:"userName"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType string        `db:"employee_type" json:"employeeType"`
}

// ---------------------------
// CUSTOMER
// ---------------------------

type CreateCustomerParams struct {
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type CreateCustomerRequest struct {
	UserName      string        `db:"username" json:"userName"`
	PasswordPlain string        `db:"password_plain" json:"passwordPlain"`
	FirstName     string        `db:"first_name" json:"firstName"`
	LastName      string        `db:"last_name" json:"lastName"`
	Email         string        `db:"email" json:"email"`
	PhonePrimary  pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther    []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type CustomerLoginRequest struct {
	UserName      string `db:"username" json:"userName"`
	PasswordPlain string `db:"password_plain" json:"passwordPlain"`
}
