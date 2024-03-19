package DB

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestNewPG(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *postgres
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPG(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPG() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanStructfromRows(t *testing.T) {
	type args struct {
		rows pgx.Rows
		dest interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := scanStructfromRows(tt.args.rows, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("scanStructfromRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scanStruct(t *testing.T) {
	type args struct {
		row  pgx.Row
		dest interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := scanStruct(tt.args.row, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("scanStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
