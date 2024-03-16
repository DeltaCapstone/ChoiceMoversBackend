package DB

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type postgres struct {
	db *pgxpool.Pool
}

var ErrRecordNotFound = pgx.ErrNoRows

var (
	PgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context) (*postgres, error) {
	pgOnce.Do(func() {
		config, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s dbname=%s host=db port=5432 sslmode=disable", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE")))
		if err != nil {
			log.Fatalf("unable to parse PostgreSQL configuration: %v", err)
		}

		db, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			// Check if the error is a PgError
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Fatalf("PostgreSQL error - Code: %s, Message: %s", pgErr.Code, pgErr.Message)
			} else {
				// If it's not a PgError, log the general error
				log.Fatalf("unable to create connection pool: %v", err)
			}
		}
		PgInstance = &postgres{db}
	})

	return PgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func scanStructfromRows(rows pgx.Rows, dest interface{}) error {
	columns := make([]interface{}, 0)
	columnsMap := make(map[string]interface{})
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	for i := 0; i < destValue.NumField(); i++ {
		field := destType.Field(i)
		column := field.Tag.Get("db")
		if column == "" {
			column = strings.ToLower(field.Name)
		}
		value := destValue.Field(i).Addr().Interface()
		columns = append(columns, value)
		columnsMap[column] = value
	}

	if rows.Next() {
		if err := rows.Scan(columns...); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("no rows found")
}

func scanStruct(row pgx.Row, dest interface{}) error {
	columns := make([]interface{}, 0)
	columnsMap := make(map[string]interface{})
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	for i := 0; i < destValue.NumField(); i++ {
		field := destType.Field(i)
		column := field.Tag.Get("db")
		if column == "" {
			column = strings.ToLower(field.Name)
		}
		value := destValue.Field(i).Addr().Interface()
		columns = append(columns, value)
		columnsMap[column] = value
	}

	if err := row.Scan(columns...); err != nil {
		zap.L().Sugar().Errorf(err.Error())
		return err
	}
	return nil
}
