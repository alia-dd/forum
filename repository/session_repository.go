package repository

import (
	"context"
	"database/sql"
	"time"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

const (
	createUserSession      = ` INSERT INTO session (uuid, user_id, expires_at) VALUES(?,?,?)`
	GetUserSession         = ` SELECT user_id, expires_at FROM session WHERE uuid = ?`
	deleteUserSession      = ` DELETE FROM session WHERE uuid = ?`
	deleteSessionsByUserId = ` DELETE FROM session WHERE user_id = ?`
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(cx context.Context, userID int) (*models.Session, error) {

	// New uuid func retund a uuid v7
	uuid, uuidErr := utils.NewUuid()
	if uuidErr != nil {
		return nil, uuidErr
	}
	session := models.Session{
		SesssionId: uuid,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(30 * (24 * time.Hour)), // the expire is currently set to 30 days
	}
	_, postErr := r.db.ExecContext(cx, createUserSession, session.SesssionId, session.UserID, session.ExpiresAt)
	if postErr != nil {
		return nil, customerrors.MapSQLError(postErr)
	}
	return &session, nil
}

func (r *SessionRepository) GetSessionwithSessionId(cx context.Context, sessionId string) (*models.Session, error) {

	session := models.Session{
		SesssionId: sessionId,
	}
	fetchErr := r.db.QueryRowContext(cx, GetUserSession, sessionId).Scan(&session.UserID, &session.ExpiresAt)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows || session.ExpiresAt.Before(time.Now()) {
			r.DeleteSession(cx, sessionId)
			return nil, customerrors.ErrNotFound
		}
		return nil, customerrors.ErrInternalError
	}
	return &session, nil
}

// what happens if the session table is deleted
// the pereiodic clean up should fix this but that takes time
// should the user sessoin validity be checked as well what is the best way to handle this the proper way
func (r *SessionRepository) DeleteSession(cx context.Context, sessionId string) error {
	resp, deleteErr := r.db.ExecContext(cx, deleteUserSession, sessionId)
	if deleteErr != nil {
		return customerrors.MapSQLError(deleteErr)
	}
	rows, _ := resp.RowsAffected()
	if rows == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// this delets the session by userId
func (r *SessionRepository) DeleteSessionsByUserId(cx context.Context, userID int) error {
	_, deleteErr := r.db.ExecContext(cx, deleteSessionsByUserId, userID)
	if deleteErr != nil {
		return customerrors.MapSQLError(deleteErr)
	}
	return nil
}
