package DB

import (
	"fmt"
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
			fmt.Println(tag, result)
		}
	}
	return result
}

type Customer struct {
	ID           int           `db:"customer_id,omitempty" json:"customerId,omitempty"`
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

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
