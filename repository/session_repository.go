package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const (
	createUserSession = ` INSERT INTO session (uuid, user_id, expires_at) VALUES(?,?,?)`
	deleteUserSession = ` DELETE session WHERE uuid = ?`
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(cx context.Context, userID int) (*models.Session, error) {

	session := models.Session{
		SesssionId: uuid.New().String(),
		UserID:     userID,
		ExpiresAt:  time.Now().Add(30 * (24 * time.Hour)),
	}
	_, postErr := r.db.ExecContext(cx, createUserSession, session.SesssionId, session.UserID, session.ExpiresAt)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return nil, customerrors.ErrDuplicateEntry
		}
		return nil, customerrors.ErrInternalError
	}
	return &session, nil
}

// what happens if the session table is deleted
// the pereiodic clean up should fix this but that takes time
// should the user sessoin validity be checked as well what is the best way to handle this the proper way
func (r *SessionRepository) DeleteSession(cx context.Context, sessionId string) error {
	_, deleteErr := r.db.ExecContext(cx, createUserSession, sessionId)
	if deleteErr != nil {
		if deleteErr == sql.ErrNoRows {
			return customerrors.ErrNotFound
		}
		return customerrors.ErrInternalError
	}
	return nil
}
