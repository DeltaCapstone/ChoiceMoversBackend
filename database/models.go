package DB

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Old
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

// Old
type Employee struct {
	ID           int           `db:"employee_id,omitempty" json:"employeeId,omitempty"`
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType string        `db:"employee_type" json:"employeeType"`
}

type Job struct {
	ID         int               `db:"job_id"`
	CustomerId int               `db:"customer_id"`
	LoadAddr   int               `db:"load_addr"`
	UnloadAddr int               `db:"unload_addr"`
	StartTime  pgtype.Timestamp  `db:"start_time"`
	HoursLabor pgtype.Interval   `db:"hours_labor"`
	Finalized  bool              `db:"finalized"`
	Rooms      pgtype.JSONBCodec `db:"rooms"`
	Pack       bool              `db:"pack"`
	Unpack     bool              `db:"unpack"`
	Load       bool              `db:"load"`
	Unload     bool              `db:"unload"`
	Clean      bool              `db:"clean"`
	Milage     int               `db:"milage"`
	Cost       float64           `db:"cost"`
}
