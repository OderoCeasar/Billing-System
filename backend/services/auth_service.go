package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/OderoCeasar/system/config"
	"github.com/OderoCeasar/system/db/models"
	"github.com/OderoCeasar/system/db/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	
)



type AuthService struct {
	userRepo 	*repositories.UserRepository
	cfg			*config.Config		
}


func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg: cfg,
	}
}


type Claims struct {
	userID		uuid.UUID	`json:"user_id"`
	PhoneNumber	string		`json:"phone_number"`
	IsAdmin		bool		`json:"is_admin"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(phoneNumber, password string) (*models.User, error) {
	existingUser, err := s.userRepo.FindByPhoneNumber(phoneNumber)
	if err == nil && existingUser != nil {
		return nil, errors.New("User already exists")
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}


	user := &models.User{
		PhoneNumber: phoneNumber,
		Username: fmt.Sprintf("user_%s", uuid.New().String()[:8]),
		Password: string(hashedPassword),
		IsActive: true,
		IsAdmin: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}


func (s *AuthService) Login(phoneNumber, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return "", nil, errors.New("Invalid crededentials")

	}

	if !user.IsActive {
		return "", nil, errors.New("account is inactive")
	}


	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// generate jwt token
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}


	return token, user, nil
}


func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := Claims{
		userID:                 user.ID,
		PhoneNumber: 			user.PhoneNumber,
		IsAdmin: 				user.IsAdmin,
		RegisteredClaims: 		jwt.RegisteredClaims{
			ExpiresAt: 	jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:	jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}


func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}


	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}


func (s *AuthService) QuickRegiter(phoneNumber string) (*models.User, string, error) {
	existingUser, err := s.userRepo.FindByPhoneNumber(phoneNumber)
	if err == nil && existingUser != nil {
		token, err := s.GenerateToken(existingUser)
		if err != nil {
			return nil, "", err
		}
		return	existingUser, token, nil
	}

	autoPassword := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(autoPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		PhoneNumber: 	phoneNumber,
		Username:		fmt.Sprintf("user_%s", uuid.New().String()[:8]),
		Password: 		string(hashedPassword),
		IsActive: 		true,
		IsAdmin: 		false,

	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}


	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}


