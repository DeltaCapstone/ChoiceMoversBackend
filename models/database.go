package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	UserName     string `db:"username" json:"userName"`
	PasswordHash string `db:"password_hash" json:"passwordHash"`
	FirstName    string `db:"first_name" json:"firstName"`
	LastName     string `db:"last_name" json:"lastName"`
	Email        string `db:"email" json:"email"`
	PhonePrimary string `db:"phone_primary" json:"phonePrimary"`
	PhoneOther1  string `db:"phone_other1" json:"phoneOther1"`
}

type EmployeeType string

const (
	FullTime EmployeeType = "Full-time"
	PartTime EmployeeType = "Part-time"
	Manager  EmployeeType = "Manager"
	Admin    EmployeeType = "Admin"
)

func IsValidEmployeeType(s string) (EmployeeType, bool) {
	for _, et := range []EmployeeType{FullTime, PartTime, Manager, Admin} {
		if string(et) == s {
			return et, true
		}
	}
	return "", false
}

type Employee struct {
	UserName         string       `db:"username" json:"userName"`
	PasswordHash     string       `db:"password_hash" json:"passwordHash"`
	FirstName        string       `db:"first_name" json:"firstName"`
	LastName         string       `db:"last_name" json:"lastName"`
	Email            string       `db:"email" json:"email"`
	PhonePrimary     string       `db:"phone_primary" json:"phonePrimary"`
	PhoneOther1      string       `db:"phone_other1" json:"phoneOther1"`
	EmployeeType     EmployeeType `db:"employee_type" json:"employeeType"`
	EmployeePriority int          `db:"employee_priority" json:"employeePriority"`
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
	JobID      int `db:"job_id" json:"jobId"`
	EstimateID int `db:"estimate_id" json:"estimateId"`

	ManHours pgtype.Interval `db:"man_hours" json:"ManHours"`
	Rate     float64         `db:"rate" json:"Rate"`
	Cost     float64         `db:"cost" json:"Cost"`

	Finalized      bool            `db:"finalized" json:"finalized"` //meaning customer agrees to all the job parameters
	ActualManHours pgtype.Interval `db:"actual_man_hours" json:"actualManHours"`
	FinalCost      float64         `db:"final_cost" json:"finalCost"`
	AmountPaid     float64         `db:"amount_payed" json:"amountPaid"`

	Notes string `db:"notes" json:"notes"`
}

// allows saving estimates with them beeing treated as jobs
type Estimate struct {
	EstimateID       int              `db:"estimate_id" json:"estimateId"`
	CustomerUsername string           `db:"customer_username" json:"customerUsername"`
	LoadAddrID       int              `db:"load_addr_id" json:"loadAddrID"`
	UnloadAddrID     int              `db:"unload_addr_id" json:"unloadAddrID"`
	StartTime        pgtype.Timestamp `db:"start_time" json:"startTime"`
	EndTime          pgtype.Timestamp `db:"end_time" json:"endTime"`

	Rooms      []Room         `db:"rooms" json:"rooms"`
	Special    map[string]int `db:"special" json:"special"`
	Small      int            `db:"small_items" json:"smallItems"`
	Medium     int            `db:"medium_items" json:"mediumItems"`
	Large      int            `db:"large_items" json:"largeItems"`
	Boxes      int            `db:"boxes" json:"boxes"`
	ItemLoad   float64        `db:"item_load" json:"itemLoad"`
	FlightMult float64        `db:"flight_mult" json:"flightMult"`

	Pack   bool `db:"pack" json:"pack"`
	Unpack bool `db:"unpack" json:"unpack"`
	Load   bool `db:"load" json:"load"`
	Unload bool `db:"unload" json:"unload"`

	Clean bool `db:"clean" json:"clean"`

	NeedTruck     bool `db:"need_truck" json:"needTruck"`
	NumberWorkers int  `db:"number_workers" json:"numberWorkers"`
	DistToJob     int  `db:"dist_to_job" json:"distToJob"`
	DistMove      int  `db:"dist_move" json:"distMove"`

	EstimateManHours pgtype.Interval `db:"estimated_man_hours" json:"estimatedManHours"`
	EstimateRate     float64         `db:"estimated_rate" json:"estimatedRate"`
	EstimateCost     float64         `db:"estimated_cost" json:"estimatedCost"`
}

type EstimateJobJoin struct {
	EstimateID       int              `db:"estimate_id" json:"estimateId"`
	CustomerUsername string           `db:"customer_username" json:"customerUsername"`
	LoadAddrID       int              `db:"load_addr_id" json:"loadAddrID"`
	UnloadAddrID     int              `db:"unload_addr_id" json:"unloadAddrID"`
	StartTime        pgtype.Timestamp `db:"start_time" json:"startTime"`
	EndTime          pgtype.Timestamp `db:"end_time" json:"endTime"`

	Rooms      map[string]interface{} `db:"rooms" json:"rooms"`
	Special    map[string]interface{} `db:"special" json:"special"`
	Small      int                    `db:"small_items" json:"smallItems"`
	Medium     int                    `db:"medium_items" json:"mediumItems"`
	Large      int                    `db:"large_items" json:"largeItems"`
	Boxes      int                    `db:"boxes" json:"boxes"`
	ItemLoad   int                    `db:"item_load" json:"itemLoad"`
	FlightMult float64                `db:"flight_mult" json:"flightMult"`

	Pack   bool `db:"pack" json:"pack"`
	Unpack bool `db:"unpack" json:"unpack"`
	Load   bool `db:"load" json:"load"`
	Unload bool `db:"unload" json:"unload"`

	Clean bool `db:"clean" json:"clean"`

	NeedTruck     bool `db:"need_truck" json:"needTruck"`
	NumberWorkers int  `db:"number_workers" json:"numberWorkers"`
	DistToJob     int  `db:"dist_to_job" json:"distToJob"`
	DistMove      int  `db:"dist_move" json:"distMove"`

	EstimateManHours pgtype.Interval `db:"estimated_man_hours" json:"estimatedManHours"`
	EstimateRate     float64         `db:"estimated_rate" json:"estimatedRate"`
	EstimateCost     float64         `db:"estimated_cost" json:"estimatedCost"`

	JobID int `db:"job_id" json:"jobId"`

	ManHours pgtype.Interval `db:"man_hours" json:"ManHours"`
	Rate     float64         `db:"rate" json:"Rate"`
	Cost     float64         `db:"cost" json:"Cost"`

	Finalized      bool            `db:"finalized" json:"finalized"` //meaning customer agrees to all the job parameters
	ActualManHours pgtype.Interval `db:"actual_man_hours" json:"actualManHours"`
	FinalCost      float64         `db:"final_cost" json:"finalCost"`
	AmountPaid     float64         `db:"amount_payed" json:"amountPaid"`

	Notes string `db:"notes" json:"notes"`
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

// /////////////////////////////////////////////////////////////////
// Password Reset

type PasswordReset struct {
	Code      string    `db:"code"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	ExpiresAt time.Time `db:"expires_at"`
}
