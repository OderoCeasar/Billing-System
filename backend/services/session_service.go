package services

import (
	"errors"
	"fmt"
	"time"
	
    "github.com/google/uuid"
	"github.com/OderoCeasar/system/db/models"
	"github.com/OderoCeasar/system/db/repositories"
)


type SessionService struct {
	sessionRepo *repositories.SessionRepository
	packageRepo *repositories.PackageRepository
	paymentRepo *repositories.PaymentRepository
}

func NewSessionService (
	packageRepo   *repositories.PackageRepository,
	sessionRepo  *repositories.SessionRepository,
	paymentRepo  *repositories.PaymentRepository,
) *SessionService {
	return &SessionService{
		packageRepo: packageRepo,
		sessionRepo: sessionRepo,
		paymentRepo: paymentRepo,
	}
}

// create session after successful payment
func (s *SessionService) CreateSession(userID, packageID, paymentID uuid.UUID) (*models.Session, error) {
	pkg, err := s.packageRepo.FindByID(packageID)
	if err != nil {
		return nil, errors.New("package not found")
	}

	activeSession, err := s.sessionRepo.FindActiveByUser(userID)
	if err == nil && activeSession != nil {
		if !activeSession.IsExpired() {
			return nil, errors.New("user already has an active session")
		}
	}

	var expiresAt time.Time
	if pkg.PackageType == models.PackageTypeTime {
		expiresAt = time.Now().Add(time.Duration(pkg.DurationMinutes) * time.Minute)

	} else {
		expiresAt = time.Now().Add(time.Duration(pkg.ValidityDays) * 24 * time.Hour)
	}

	session := &models.Session{
		UserID: 		userID,
		PackageID: 		packageID,
		PaymentID: 		paymentID,
		Status:			models.SessionStatusActive,
		StartTime: 		time.Now(),
		ExpiresAt: 		expiresAt,
		DataUsedBytes:  0,
		TimeUsedMinutes: 0,
		LastUpdateTime:  time.Now(),
	}

	if pkg.PackageType == models.PackageTypeTime {
		session.TimeLimitMinutes = pkg.DurationMinutes
	} else if pkg.PackageType == models.PackageTypeData {
		session.DataLimitBytes = pkg.DataLimitMB * 1024 * 1024
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}


// Active session for a user
func (s *SessionService) GetActiveSession(userID uuid.UUID) (*models.Session, error) {
	return s.sessionRepo.FindActiveByUser(userID)
}

// Disconnect manually 
func (s *SessionService) DisconnectSession(sessionID uuid.UUID) error {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return err
	}

	if session.Status != models.SessionStatusActive {
		return errors.New("session is not active")
	}

	session.Status = models.SessionStatusDisconnected
	now := time.Now()
	session.EndTime = &now

	return s.sessionRepo.Update(session)
}


// expired old session
func (s *SessionService) ExpireOldSessions() (int, error) {
	expiredSessions, err := s.sessionRepo.ListExpired()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, session := range expiredSessions {
		session.Status = models.SessionStatusExpired
		now := time.Now()
		if session.EndTime == nil {
			session.EndTime = &now
		}
		if err := s.sessionRepo.Update(&session); err == nil {
			count ++
		}
	}
	return count, nil
}

// Session statistics
func (s *SessionService) GetSessionStats(sessionID uuid.UUID) (map[string]interface{}, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, err
	}

	var remainingTime int
	var remainingData int64
	var percentUsed float64

	if session.TimeLimitMinutes > 0 {
		remainingTime = session.TimeLimitMinutes - session.TimeUsedMinutes
		if remainingTime < 0 {
			remainingTime = 0
		}
		percentUsed = float64(session.TimeUsedMinutes) / float64(session.TimeLimitMinutes) * 100

	}

	if session.DataLimitBytes > 0 {
		remainingData = session.DataLimitBytes - session.DataUsedBytes
		if remainingData < 0 {
			remainingData = 0
		}
		percentUsed = float64(session.DataUsedBytes) / float64(session.DataLimitBytes) * 100
	}

	stats := map[string]interface{}{
		"session_id":		session.ID,
		"status":			session.Status,
		"start_time":		session.StartTime,
		"expires_at":		session.ExpiresAt,
		"time_used_minutes":	session.TimeUsedMinutes,
		"time_limit_minutes": 	session.TimeLimitMinutes,
		"remaining_minutes": 	remainingTime,
		"data_used_md":		float64(session.DataUsedBytes) / (1024 * 1024),
		"data_limit_mb":	float64(session.DataLimitBytes) /(1024 * 1024),
		"remaining_data_mb": float64(remainingData) / (1024 * 1024),
		"percent_used":		fmt.Sprintf("%.2f", percentUsed),
		"is_active": 	session.Status == models.SessionStatusActive,
	}

	return stats, nil
}


// list all active sessions
func (s *SessionService) ListActiveSessions() ([]models.Session, error) {
	return s.sessionRepo.ListActive()
}

// list user sessions
func (s *SessionService) ListUserSession(userID uuid.UUID, limit, offset int) ([]models.Session, error) {
	return s.sessionRepo.ListByUser(userID, limit, offset)
}