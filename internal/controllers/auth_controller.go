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
	"gorm.io/gorm"
)

type AuthController struct {
	UserRepo  *repositories.UserRepository
	TokenRepo *repositories.PersonalAccessTokenRepository
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		TokenRepo: repositories.NewPersonalAccessTokenRepository(db),
		UserRepo:  repositories.NewUserRepository(db),
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Authenticate a user
// @Description Login with email and password to receive a personal access token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response{data=LoginResponse}
// @Failure 401 {object} utils.Response
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	data := make(map[string]interface{})

	rules := govalidator.MapData{
		"email":    []string{"required", "email"},
		"password": []string{"required", "min:6"},
	}

	errs := utils.ValidateJSON(c, &data, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, err := ctrl.UserRepo.FindByEmail(data["email"].(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"].(string)))
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

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Locals("token").(models.PersonalAccessToken)
	ctrl.TokenRepo.Delete(&token)

	return utils.SuccessResponse(c, "Logout successful. Current access token has been revoked.", nil)
}

type ValidateTokenUserResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ValidateTokenResponse struct {
	User ValidateTokenUserResponse `json:"user"`
}

// ValidateToken godoc
// @Summary Validate authentication token
// @Description Validate the personal access token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=ValidateTokenResponse}
// @Failure 401 {object} utils.Response
// @Router /auth/validate-token [get]
func (ctrl *AuthController) ValidateToken(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	return utils.SuccessResponse(c, "Token is valid", ValidateTokenResponse{
		User: ValidateTokenUserResponse{
			ID:   user.ID,
			Code: user.Code,
			Name: user.Name,
		},
	})
}

func generateRandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}
