package DB

import (
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

func structToMap(data interface{}, tag string) map[string]interface{} {
	result := make(map[string]interface{})
	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tag) // You can customize this tag based on your needs

		// If the tag is not empty, use it as the key in the map
		if tag != "" {
			result[tag] = value.Field(i).Interface()
		}
	}

	return result
}

type Customer struct {
	ID           int           `db:"customer_id,omitempty" json:"CustomerId,omitempty"`
	UserName     string        `db:"username"`
	PasswordHash string        `db:"password_hash" json:"Password"`
	FirstName    string        `db:"first_name"`
	LastName     string        `db:"last_name"`
	Email        string        `db:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary"`
	PhoneOther   []pgtype.Text `db:"phone_other"`
}

type Employee struct {
	ID           int           `db:"employee_id,omitempty" json:"EmployeeId,omitempty"`
	UserName     string        `db:"username"`
	PasswordHash string        `db:"password_hash" json:"Password"`
	FirstName    string        `db:"first_name"`
	LastName     string        `db:"last_name"`
	Email        string        `db:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary"`
	PhoneOther   []pgtype.Text `db:"phone_other"`
	EmployeeType string        `db:"employee_type"`
}
