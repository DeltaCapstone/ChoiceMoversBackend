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

func (pg *postgres) CreateSession(c context.Context, s models.CreateSessionParams) (uuid.UUID, error) {
	row := pg.db.QueryRow(c, createSession, pgx.NamedArgs(utils.StructToMap(s, "db")))
	var sid uuid.UUID
	err := row.Scan(&sid)
	return sid, err
}
