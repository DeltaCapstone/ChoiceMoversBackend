package models

import (
	"fmt"
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
	Id               uuid.UUID    `db:"id" json:"id"`
	Email            string       `db:"email" json:"email"`
	EmployeeType     EmployeeType `db:"employee_type" json:"employeeType"`
	EmployeePriority int          `db:"employee_priority" json:"employeePriority"`
	SignupToken      string       `db:"signup_token" json:"signupToken"`
	ExpiresAt        time.Time    `db:"expires_at"`
	Used             bool         `db:"used"`
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

	ManHours string  `db:"man_hours" json:"ManHours"`
	Rate     float64 `db:"rate" json:"Rate"`
	Cost     float64 `db:"cost" json:"Cost"`

	Finalized      bool    `db:"finalized" json:"finalized"` //meaning customer agrees to all the job parameters
	ActualManHours string  `db:"actual_man_hours" json:"actualManHours"`
	FinalCost      float64 `db:"final_cost" json:"finalCost"`
	AmountPaid     float64 `db:"ammount_payed" json:"ammountPaid"`

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

	EstimateManHours string  `db:"estimated_man_hours" json:"estimatedManHours"`
	EstimateRate     float64 `db:"estimated_rate" json:"estimatedRate"`
	EstimateCost     float64 `db:"estimated_cost" json:"estimatedCost"`
}

func intervalToISO(i pgtype.Interval) string {
	us := i.Microseconds
	d := i.Days
	temp, _ := time.ParseDuration(fmt.Sprintf("%vus", us))
	val := fmt.Sprintf("P%vT", d) + temp.String()
	return val
}

func (jr *JobResponse) MakeFromJoin(ej EstimateJobJoin) {
	jr.JobID = ej.JobID
	jr.ManHours = intervalToISO(ej.ManHours)
	jr.Rate = ej.Rate
	jr.Cost = ej.Cost
	jr.Finalized = ej.Finalized
	jr.ActualManHours = intervalToISO(ej.ActualManHours)
	jr.FinalCost = ej.FinalCost
	jr.AmountPaid = ej.AmountPaid
	jr.Notes = ej.Notes

}

func (er *EstimateResponse) MakeFromJoin(ej EstimateJobJoin) {
	er.EstimateID = ej.EstimateID
	er.StartTime = ej.StartTime
	er.EndTime = ej.EndTime
	er.Rooms = ej.Rooms
	er.Special = ej.Special
	er.Small = ej.Small
	er.Medium = ej.Medium
	er.Large = ej.Large
	er.Boxes = ej.Boxes
	er.ItemLoad = ej.ItemLoad
	er.FlightMult = ej.FlightMult
	er.Pack = ej.Pack
	er.Unpack = ej.Unpack
	er.Load = ej.Load
	er.Unload = ej.Unload
	er.Clean = ej.Clean
	er.NeedTruck = ej.NeedTruck
	er.NumberWorkers = ej.NumberWorkers
	er.DistToJob = ej.DistToJob
	er.DistMove = ej.DistMove
	er.EstimateManHours = intervalToISO(ej.EstimateManHours)
	er.EstimateRate = ej.EstimateRate
	er.EstimateCost = ej.EstimateCost
}
