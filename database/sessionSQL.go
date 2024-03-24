package DB

import (
	"context"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const createSession = ` INSERT INTO sessions (
	id, username,role ,refresh_token, user_agent, client_ip, is_blocked, expires_at ) VALUES (
	@id, @username, @role, @refresh_token, @user_agent, @client_ip, @is_blocked, @expires_at) RETURNING id`

func (pg *postgres) CreateSession(ctx context.Context, s models.CreateSessionParams) (uuid.UUID, error) {
	row := pg.db.QueryRow(ctx, createSession, pgx.NamedArgs(utils.StructToMap(s, "db")))
	var sid uuid.UUID
	err := row.Scan(&sid)
	return sid, err
}

const getSession = `SELECT id, username, role, refresh_token, user_agent, 
client_ip, is_blocked, expires_at, created_at FROM sessions WHERE id = $1`

func (pg *postgres) GetSession(ctx context.Context, sid uuid.UUID) (models.Session, error) {
	var s models.Session
	row := pg.db.QueryRow(ctx, getSession, sid)
	err := scanStruct(row, &s)
	return s, err
}
