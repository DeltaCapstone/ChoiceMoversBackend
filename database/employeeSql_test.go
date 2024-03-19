package DB

import (
	"context"
	"reflect"
	"testing"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/google/uuid"
)

func Test_postgres_GetEmployeeCredentials(t *testing.T) {
	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pg.GetEmployeeCredentials(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetEmployeeCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("postgres.GetEmployeeCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_GetEmployeeByUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		want    models.GetEmployeeResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pg.GetEmployeeByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetEmployeeByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgres.GetEmployeeByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_GetEmployeeRole(t *testing.T) {
	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pg.GetEmployeeRole(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetEmployeeRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("postgres.GetEmployeeRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_DeleteEmployeeByUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pg.DeleteEmployeeByUsername(tt.args.ctx, tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("postgres.DeleteEmployeeByUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgres_GetEmployeeList(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		want    []models.GetEmployeeResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pg.GetEmployeeList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetEmployeeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgres.GetEmployeeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_AddEmployeeSignup(t *testing.T) {
	type args struct {
		ctx               context.Context
		newEmployeeSignUp models.EmployeeSignup
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pg.AddEmployeeSignup(tt.args.ctx, tt.args.newEmployeeSignUp); (err != nil) != tt.wantErr {
				t.Errorf("postgres.AddEmployeeSignup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgres_GetEmployeeSignup(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		want    models.EmployeeSignup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pg.GetEmployeeSignup(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetEmployeeSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgres.GetEmployeeSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_UseEmployeeSignup(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pg.UseEmployeeSignup(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("postgres.UseEmployeeSignup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgres_CreateEmployee(t *testing.T) {
	type args struct {
		ctx         context.Context
		newEmployee models.CreateEmployeeParams
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pg.CreateEmployee(tt.args.ctx, tt.args.newEmployee); (err != nil) != tt.wantErr {
				t.Errorf("postgres.CreateEmployee() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgres_UpdateEmployee(t *testing.T) {
	type args struct {
		ctx             context.Context
		updatedEmployee models.UpdateEmployeeParams
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pg.UpdateEmployee(tt.args.ctx, tt.args.updatedEmployee); (err != nil) != tt.wantErr {
				t.Errorf("postgres.UpdateEmployee() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
