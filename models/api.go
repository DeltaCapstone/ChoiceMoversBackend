package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------
// EMPLOYEE
// ---------------------------

type CreateEmployeeParams struct {
	UserName         string        `db:"username" json:"userName"`
	PasswordHash     string        `db:"password_hash" json:"passwordHash"`
	FirstName        string        `db:"first_name" json:"firstName"`
	LastName         string        `db:"last_name" json:"lastName"`
	Email            string        `db:"email" json:"email"`
	PhonePrimary     pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther       []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType     string        `db:"employee_type" json:"employeeType"`
	EmployeePriority int           `db:"employee_priority" json:"employeePriority"`
}

type EmployeeSignup struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Email       string    `db:"email" json:"email"`
	SignupToken string    `db:"signup_token" json:"signupToken"`
	ExpiresAt   time.Time `db:"expires_at"`
	Used        bool      `db:"used"`
}

type CreateEmployeeRequest struct {
	UserName         string        `db:"username" json:"userName"`
	PasswordPlain    string        `db:"password_plain" json:"passwordPlain"`
	FirstName        string        `db:"first_name" json:"firstName"`
	LastName         string        `db:"last_name" json:"lastName"`
	Email            string        `db:"email" json:"email"`
	PhonePrimary     pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther       []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType     string        `db:"employee_type" json:"employeeType"`
	EmployeePriority int           `db:"employee_priority" json:"employeePriority"`
}

type UpdateEmployeeParams struct {
	UserName     string        `db:"username" json:"userName"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type UpdateEmployeeTypePriorityParams struct {
	UserName         string `db:"username" json:"userName"`
	EmployeeType     string `db:"employee_type" json:"employeeType"`
	EmployeePriority int    `db:"employee_priority" json:"employeePriority"`
}

type UpdateEmployeePasswordRequest struct {
	UserName    string `json:"userName"`
	PasswordOld string `json:"passwordOld"`
	PasswordNew string `json:"passwordNew"`
}

type EmployeeLoginRequest struct {
	UserName      string `db:"username" json:"userName"`
	PasswordPlain string `db:"password_plain" json:"passwordPlain"`
}

type GetEmployeeResponse struct {
	UserName         string        `db:"username" json:"userName"`
	FirstName        string        `db:"first_name" json:"firstName"`
	LastName         string        `db:"last_name" json:"lastName"`
	Email            string        `db:"email" json:"email"`
	PhonePrimary     pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther       []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType     string        `db:"employee_type" json:"employeeType"`
	EmployeePriority int           `db:"employee_priority" json:"employeePriority"`
}

type EmployeeLoginResponse struct {
	SessionId             uuid.UUID `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Username              string    `json:"username"`
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

type UpdateCustomerParams struct {
	UserName     string        `db:"username" json:"userName"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type UpdateCustomerPassowrdRequest struct {
	UserName    string `json:"userName"`
	PasswordOld string `json:"passwordOld"`
	PasswordNew string `json:"passwordNew"`
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

type GetCustomerResponse struct {
	UserName     string        `db:"username" json:"userName"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type CustomerLoginRequest struct {
	UserName      string `db:"username" json:"userName"`
	PasswordPlain string `db:"password_plain" json:"passwordPlain"`
}

type CustomerLoginResponse struct {
	SessionId             uuid.UUID `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Username              string    `json:"username"`
}

// ---------------------------
// JOBS
// ---------------------------

type JobsDisplayRequest struct {
	Status    string `json:"status"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type JobResponse struct {
	ID          int                    `db:"job_id" json:"id"`
	Customer    GetCustomerResponse    `json:"customer"`
	LoadAddr    Address                `json:"loadAddr"`
	UnloadAddr  Address                `json:"unloadAddr"`
	StartTime   pgtype.Timestamp       `db:"start_time" json:"startTime"`
	HoursLabor  pgtype.Interval        `db:"hours_labor" json:"hoursLabor"`
	Finalized   bool                   `db:"finalized" json:"finalized"`
	Rooms       map[string]interface{} `db:"rooms" json:"rooms"`
	Pack        bool                   `db:"pack" json:"pack"`
	Unpack      bool                   `db:"unpack" json:"unpack"`
	Load        bool                   `db:"load" json:"load"`
	Unload      bool                   `db:"unload" json:"unload"`
	Clean       bool                   `db:"clean" json:"clean"`
	Milage      int                    `db:"milage" json:"milage"`
	Cost        string                 `db:"cost" json:"cost"`
	Notes       pgtype.Text            `db:"notes" json:"notes"`
	AssignedEmp []GetEmployeeResponse  `json:"assignedEmployees"`
}
