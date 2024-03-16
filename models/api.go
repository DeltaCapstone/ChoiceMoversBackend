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
	JobID int `db:"job_id" json:"jobId"`
	EstimateResponse

	ManHours pgtype.Interval `db:"man_hours" json:"ManHours"`
	Rate     float64         `db:"rate" json:"Rate"`
	Cost     float64         `db:"cost" json:"Cost"`

	Finalized      bool            `db:"finalized" json:"finalized"` //meaning customer agrees to all the job parameters
	ActualManHours pgtype.Interval `db:"actual_man_hours" json:"actualManHours"`
	FinalCost      float64         `db:"final_cost" json:"finalCost"`
	AmmountPaid    float64         `db:"ammount_payed" json:"ammountPaid"`

	Notes       pgtype.Text           `db:"notes" json:"notes"`
	AssignedEmp []GetEmployeeResponse `json:"assignedEmployees"`
}

type EstimateResponse struct {
	EstimateID int                 `db:"estimate_id" json:"estimateId"`
	Customer   GetCustomerResponse `json:"customer"`
	LoadAddr   Address             `json:"loadAddr"`
	UnloadAddr Address             `json:"unloadAddr"`
	StartTime  pgtype.Timestamp    `db:"start_time" json:"startTime"`
	EndTime    pgtype.Timestamp    `db:"end_time" json:"endTime"`

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
}

func (jr *JobResponse) MakeFromJoin(ej EstimateJobJoin) {
	jr.JobID = ej.JobID
	jr.ManHours = ej.ManHours
	jr.Rate = ej.Rate
	jr.Cost = ej.Cost
	jr.Finalized = ej.Finalized
	jr.ActualManHours = ej.ActualManHours
	jr.FinalCost = ej.FinalCost
	jr.AmmountPaid = ej.AmmountPaid
	jr.Notes = ej.Notes
}

func (er *EstimateResponse) MakeFromEstimate(e Estimate) {
	er.EstimateID = e.EstimateID
	er.StartTime = e.StartTime
	er.EndTime = e.EndTime
	er.Rooms = e.Rooms
	er.Special = e.Special
	er.Small = e.Small
	er.Medium = e.Medium
	er.Large = e.Large
	er.Boxes = e.Boxes
	er.ItemLoad = e.ItemLoad
	er.FlightMult = e.FlightMult
	er.Pack = e.Pack
	er.Unpack = e.Unpack
	er.Load = e.Load
	er.Unload = e.Unload
	er.Clean = e.Clean
	er.NeedTruck = e.NeedTruck
	er.NumberWorkers = e.NumberWorkers
	er.DistToJob = e.DistToJob
	er.DistMove = e.DistMove
	er.EstimateManHours = e.EstimateManHours
	er.EstimateRate = e.EstimateRate
	er.EstimateCost = e.EstimateCost
}
