package DB

import (
	"context"
	"reflect"
	"testing"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
)

func Test_postgres_GetCustomerCredentials(t *testing.T) {
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
			got, err := tt.pg.GetCustomerCredentials(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetCustomerCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("postgres.GetCustomerCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_GetCustomerByUserName(t *testing.T) {
	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		pg      *postgres
		args    args
		want    models.GetCustomerResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pg.GetCustomerByUserName(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetCustomerByUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgres.GetCustomerByUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_CreateCustomer(t *testing.T) {
	type args struct {
		ctx         context.Context
		newCustomer models.CreateCustomerParams
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
			got, err := tt.pg.CreateCustomer(tt.args.ctx, tt.args.newCustomer)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.CreateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("postgres.CreateCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgres_UpdateCustomer(t *testing.T) {
	type args struct {
		ctx             context.Context
		updatedCustomer models.UpdateCustomerParams
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
			if err := tt.pg.UpdateCustomer(tt.args.ctx, tt.args.updatedCustomer); (err != nil) != tt.wantErr {
				t.Errorf("postgres.UpdateCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
