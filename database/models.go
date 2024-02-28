package DB

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	ID           int           `db:"customer_id,omitempty" json:"customer,omitempty"`
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
	PartTime EmployeeType = "Part-Time"
	Manager  EmployeeType = "Manager"
	Admin    EmployeeType = "Admin"
)

type Employee struct {
	ID           int           `db:"employee_id,omitempty" json:"employeeId,omitempty"`
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType EmployeeType  `db:"employee_type" json:"employeeType"`
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
	AddressID int           `db:"address_id,omitempty" json:"addressId"`
	Street    string        `db:"street" json:"street"`
	City      string        `db:"city" json:"city"`
	State     string        `db:"state" json:"state"`
	Zip       string        `db:"zip" json:"zip"`
	ResType   ResidenceType `db:"res_type" json:"resType"`
	Flights   int           `db:"flights" json:"flights"`
	AptNum    string        `db:"apt_num" json:"aptNum"`
}

type Job struct {
	ID         int                    `db:"job_id" json:"id"`
	CustomerID int                    `db:"customer_id" json:"customerId"`
	LoadAddr   Address                `db:"load_addr" json:"loadAddr"`
	UnloadAddr Address                `db:"unload_addr" json:"unloadAddr"`
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
	Cost       string                 `db:"cost" json:"cost"`
}
