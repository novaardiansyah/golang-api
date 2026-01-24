package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, error)
	ChangePassword(user *models.User, currentPassword, newPassword string) error
}

type authService struct {
	UserRepo  *repositories.UserRepository
	TokenRepo *repositories.PersonalAccessTokenRepository
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		UserRepo:  repositories.NewUserRepository(db),
		TokenRepo: repositories.NewPersonalAccessTokenRepository(db),
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid_credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid_credentials")
	}

	expireDays := 7
	hashedToken, plainToken := generateToken(40)
	expiration := time.Now().AddDate(0, 0, expireDays)

	token := models.PersonalAccessToken{
		TokenableType: "App\\Models\\User",
		TokenableID:   user.ID,
		Name:          "auth_token",
		Token:         hashedToken,
		Abilities:     "[\"*\"]",
		ExpiresAt:     &expiration,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.TokenRepo.Create(&token); err != nil {
		return "", errors.New("token_creation_failed")
	}

	fullToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

	return fullToken, nil
}

func (s *authService) ChangePassword(user *models.User, currentPassword, newPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return errors.New("current_password_incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	user.Password = strings.Replace(string(hashed), "$2a$", "$2y$", 1)

	return s.UserRepo.Update(user)
}

func generateToken(length int) (string, string) {
	bytes := make([]byte, length)
	rand.Read(bytes)
	plainToken := hex.EncodeToString(bytes)[:length]

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	return hashedToken, plainToken
}
