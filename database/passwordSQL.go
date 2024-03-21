package DB

import (
	"context"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgx/v5"
)

const createPasswordReset = ` INSERT INTO password_resets (
	code, username,email, role, expires_at ) VALUES (
	@code, @username, @email, @role, @expires_at) RETURNING
	code, username,email, role, expires_at`

func (pg *postgres) CreatePasswordReset(ctx context.Context, s models.PasswordReset) (models.PasswordReset, error) {
	var pwr models.PasswordReset
	row := pg.db.QueryRow(ctx, createPasswordReset, pgx.NamedArgs(utils.StructToMap(s, "db")))
	err := scanStruct(row, &pwr)
	return pwr, err
}

const getPasswordReset = `SELECT code, username, email, role, expires_at FROM password_resets WHERE code = $1`

func (pg *postgres) GetPasswordReset(ctx context.Context, code string) (models.PasswordReset, error) {
	var pwr models.PasswordReset
	row := pg.db.QueryRow(ctx, getPasswordReset, code)
	err := scanStruct(row, &pwr)
	return pwr, err
}

const usePasswordReset = `DELETE FROM password_resets WHERE code = $1`

func (pg *postgres) DeletePasswordReset(ctx context.Context, code string) error {
	_, err := pg.db.Exec(ctx, usePasswordReset, code)
	return err
}
