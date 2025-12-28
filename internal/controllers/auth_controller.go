package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang-api/internal/config"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userRepo repositories.UserRepository
}

func NewAuthController(userRepo repositories.UserRepository) *AuthController {
	return &AuthController{userRepo: userRepo}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest

	rules := govalidator.MapData{
		"email":    []string{"required", "email"},
		"password": []string{"required", "min:6"},
	}

	errs := utils.ValidateJSONStruct(c, &req, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, err := ctrl.userRepo.FindByEmail(req.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	plainToken := generateRandomToken(40)

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	expiration := time.Now().AddDate(0, 0, 7)

	db := config.GetDB()
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

	if err := db.Create(&token).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create token")
	}

	fullToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

	return utils.SuccessResponse(c, "Login successful", LoginResponse{
		Token: fullToken,
	})
}

func generateRandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}
