package repositories

import (
	"time"

	"github.com/OderoCeasar/system/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindByID(id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := r.db.Preload("User").Preload("Package").Where("id = ?", id).First(&session).Error
	return &session, err
}

func (r *SessionRepository) FindBySessionID(sessionID string) (*models.Session, error) {
	var session models.Session
	err := r.db.Preload("User").Preload("Package").Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

func (r *SessionRepository) FindActiveByUser(userID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := r.db.Preload("Package").
		Where("user_id = ? AND status = ?", userID, models.SessionStatusActive).
		First(&session).Error
	return &session, err
}

func (r *SessionRepository) Update(session *models.Session) error {
	return r.db.Save(session).Error
}

func (r *SessionRepository) ListActive() ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.Preload("User").Preload("Package").
		Where("status = ?", models.SessionStatusActive).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) ListExpired() ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.Preload("User").Preload("Package").
		Where("status = ? AND expires_at < ?", models.SessionStatusActive, time.Now()).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) ListByUser(userID uuid.UUID, limit, offset int) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.Preload("Package").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) List(limit, offset int) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.Preload("User").Preload("Package").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&sessions).Error
	return sessions, err
}