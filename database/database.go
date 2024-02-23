package DB

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

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

// TODO
// will neet to move and probably change
type User struct {
	ID          int
	UserName    string
	AccountType string
	Email       string
}

// basic querry for retrieving id for name
func (pg *postgres) GetUsers(ctx context.Context, accountType string) ([]User, error) {
	var users []User
	var rows pgx.Rows
	var err error
	if accountType != "" {
		rows, err = pg.db.Query(ctx, "select user_id, username, accnt_type, email from users where accnt_type = $1", accountType)
	} else {
		rows, err = pg.db.Query(ctx, "select user_id, username, accnt_type, email from users")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.UserName, &user.AccountType, &user.Email); err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (pg *postgres) CreateUser(ctx context.Context, user User) (int, error) {
	var userID int
	err := pg.db.QueryRow(ctx, "INSERT INTO users (username, accnt_type, email) VALUES ($1, $2, $3) RETURNING user_id", user.UserName, user.AccountType, user.Email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %v", err)
	}
	return userID, nil
}
