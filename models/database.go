package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type EmployeeType string

const (
	FullTime EmployeeType = "Full-time"
	PartTime EmployeeType = "Part-time"
	Manager  EmployeeType = "Manager"
	Admin    EmployeeType = "Admin"
)

type Employee struct {
	UserName         string        `db:"username" json:"userName"`
	PasswordHash     string        `db:"password_hash" json:"passwordHash"`
	FirstName        string        `db:"first_name" json:"firstName"`
	LastName         string        `db:"last_name" json:"lastName"`
	Email            string        `db:"email" json:"email"`
	PhonePrimary     pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther       []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType     EmployeeType  `db:"employee_type" json:"employeeType"`
	EmployeePriority int           `db:"employee_priority" json:"employeePriority"`
}

type ResidenceType string

const (
	House     ResidenceType = "House"
	Apartment ResidenceType = "Apartment"
	Condo     ResidenceType = "Condo"
	Business  ResidenceType = "Business"
	Storage   ResidenceType = "Storage Unit"
	Other     ResidenceType = "Other"
)

type Address struct {
	AddressID  int           `db:"address_id,omitempty" json:"addressId"`
	Street     string        `db:"street" json:"street"`
	City       string        `db:"city" json:"city"`
	State      string        `db:"state" json:"state"`
	Zip        string        `db:"zip" json:"zip"`
	ResType    ResidenceType `db:"res_type" json:"resType"`
	SquareFeet int           `db:"square_feet" json:"sqareFeet"`
	Flights    int           `db:"flights" json:"flights"`
	AptNum     string        `db:"apt_num" json:"aptNum"`
}

type Job struct {
	ID         int                    `db:"job_id" json:"id"`
	Customer   string                 `db:"customer_username" json:"customerUsername"`
	LoadAddr   int                    `db:"load_addr" json:"loadAddr"`
	UnloadAddr int                    `db:"unload_addr" json:"unloadAddr"`
	StartTime  pgtype.Timestamp       `db:"start_time" json:"startTime"`
	HoursLabor pgtype.Interval        `db:"hours_labor" json:"hoursLabor"`
	Finalized  bool                   `db:"finalized" json:"finalized"`
	Rooms      map[string]interface{} `db:"rooms" json:"rooms"`
	Pack       bool                   `db:"pack" json:"pack"`
	Unpack     bool                   `db:"unpack" json:"unpack"`
	Load       bool                   `db:"load" json:"load"`
	Unload     bool                   `db:"unload" json:"unload"`
	Clean      bool                   `db:"clean" json:"clean"`
	Milage     int                    `db:"milage" json:"milage"`
	Notes      pgtype.Text            `db:"notes" json:"notes"`
	Cost       string                 `db:"cost" json:"cost"`
}

// /////////////////////////////////////////////////////////////////
// Session
type Session struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Role         string    `db:"role" json:"role"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	ClientIp     string    `db:"client_ip" json:"client_ip"`
	IsBlocked    bool      `db:"is_blocked" json:"is_blocked"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt    time.Time `db:"create_at" json:"created_at"`
}

type CreateSessionParams struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Role         string    `db:"role" json:"role"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	ClientIp     string    `db:"client_ip" json:"client_ip"`
	IsBlocked    bool      `db:"is_blocked" json:"is_blocked"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
}
