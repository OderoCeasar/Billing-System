package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/OderoCeasar/system/config"
	"github.com/OderoCeasar/system/db/models"
	"github.com/OderoCeasar/system/db/repositories"
	"golang.org/x/crypto/bcrypt"
)


type RADIUSService struct {
	userRepo		*repositories.UserRepository
	sessionRepo  	*repositories.SessionRepository
	cfg				*config.Config
}


func NewRADIUSService (
	userRepo   *repositories.UserRepository,
	sessionRepo  *repositories.SessionRepository,
	cfg		*config.Config,
) *RADIUSService {
	return &RADIUSService{
		userRepo: userRepo,
		sessionRepo: sessionRepo,
		cfg: 	cfg,
	}
}


// AccountingUpdate handles RADIUS accounting interim-update requests
func (s *RADIUSService) AuthenticateUser(username, password string) (bool, *models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		user, err = s.userRepo.FindByPhoneNumber(username)
		if err != nil {
			return false, nil, fmt.Errorf("user not found")
		}
	}

	if !user.IsActive {
		return false, nil, fmt.Errorf("user account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, nil, fmt.Errorf("invalid credentials")
	}

	// check if user has active session
	activeSession, err := s.sessionRepo.FindActiveByUser(user.ID)
	if err == nil && activeSession != nil {
		if !activeSession.IsExpired() && !activeSession.HasExceededDataLimit() && !activeSession.HasExceededTimeLimit() {
			return true, user, nil
		}

		activeSession.Status = models.SessionStatusExpired
		if activeSession.EndTime == nil {
			now := time.Now()
			activeSession.EndTime = &now
		}
		s.sessionRepo.Update(activeSession)
		return false, nil, fmt.Errorf("session expired, please purchase a new package")
	}

	return false, nil, fmt.Errorf("no active package, purchase another")
}


// AccountingStart handles RADIUS accounting start requests
func (s *RADIUSService) AccountingStart(sessionID, username, nasIP, nasPort, userIP, macAddress string) error {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		user, err := s.userRepo.FindByPhoneNumber(username)
		if err != nil {
			return fmt.Errorf("user not found")
		}
	}

	activeSession, err := s.sessionRepo.FindActiveByUser(user.ID)
	if err != nil {
		return fmt.Errorf("no active session found")
	}

	activeSession.SessionID = sessionID
	activeSession.NASIPAddress = nasIP
	activeSession.NASPortID = nasPort
	activeSession.IPAddress = userIP
	activeSession.MACAddress = macAddress
	activeSession.LastUpdateTime = time.Now()

	return s.sessionRepo.Update(activeSession)
}


// AccountingUpdate handles RADIUS accounting interim-update requests
func (s *RADIUSService) AccountingUpdate(sessionID string, sessionTime int64, inputOctets, outputOctets int64) error {
	session, err := s.sessionRepo.FindBySessionID(sessionID)
	if err != nil {
		return fmt.Errorf("session not found")
	}

	session.TimeUsedMinutes = int(sessionTime / 60)
	session.DataUsedBytes = inputOctets + outputOctets
	session.LastUpdateTime = time.Now()

	if session.HasExceededDataLimit() || session.HasExceededTimeLimit() || session.IsExpired() {
		session.Status = models.SessionStatusExpired
		now := time.Now()
		session.EndTime = &now
	}

	return s.sessionRepo.Update(session)
}

// AccountingStop handles RADIUS accounting stop requests
func (s *RADIUSService) AccountingStop(sessionID string, sessionTime int64, inputOctets, outputOctets int64) error {
	session, err := s.sessionRepo.FindBySessionID(sessionID)
	if err != nil {
		return fmt.Errorf("session not found")
	}


	session.TimeUsedMinutes = int(sessionTime / 60)
	session.DataUsedBytes = inputOctets + outputOctets
	session.Status = models.SessionStatusDisconnected
	now := time.Now()
	session.EndTime = &now
	session.LastUpdateTime = now


	return s.sessionRepo.Update(session)
}

// checks if user is authorized for internet access
func (s *RADIUSService) CheckAuthorization(username string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		user, err = s.userRepo.FindByPhoneNumber(username)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}
	}

	activeSession, err := s.sessionRepo.FindActiveByUser(user.ID)
	if err != nil {
		return nil, fmt.Errorf("no active session")
	}

	if activeSession.IsExpired() || activeSession.HasExceededDataLimit() || activeSession.HasExceededTimeLimit() {
		activeSession.Status = models.SessionStatusExpired
		now := time.Now()
		activeSession.EndTime = &now
		s.sessionRepo.Update(activeSession)
		return nil, fmt.Errorf("session expired")
	}

	attributes := map[string]interface{} {
		"session_timeout": activeSession.ExpiresAt.Sub(time.Now()).Seconds(),
		"idle_timeout": 3600,
	}

	if activeSession.Package.SpeedLimitUp > 0 {
		attributes["upload_speed"] = activeSession.Package.SpeedLimitUp

	}

	if activeSession.Package.SpeedLimitDown > 0 {
		attributes["download_speed"] = activeSession.Package.SpeedLimitDown

	}

	return attributes, nil

	
}
