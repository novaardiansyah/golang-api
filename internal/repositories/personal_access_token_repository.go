package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PersonalAccessTokenRepository struct {
	db *gorm.DB
}

func NewPersonalAccessTokenRepository(db *gorm.DB) *PersonalAccessTokenRepository {
	return &PersonalAccessTokenRepository{db: db}
}

type TokenWithUser struct {
	models.PersonalAccessToken
	UserName string
}

func (repo PersonalAccessTokenRepository) FindByIDAndHashedToken(id uint64, hashedToken string) (*models.PersonalAccessToken, error) {
	var token models.PersonalAccessToken

	result := repo.db.Where("id = ? AND token = ?", id, hashedToken).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}

	return &token, nil
}

func (repo PersonalAccessTokenRepository) FindByIDAndHashedTokenWithUser(id uint64, hashedToken string) (*TokenWithUser, error) {
	var result TokenWithUser

	err := repo.db.Table("personal_access_tokens").
		Select("personal_access_tokens.*, users.name as user_name").
		Joins("INNER JOIN users ON users.id = personal_access_tokens.tokenable_id").
		Where("personal_access_tokens.id = ? AND personal_access_tokens.token = ?", id, hashedToken).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo PersonalAccessTokenRepository) Delete(token *models.PersonalAccessToken) error {
	return repo.db.Delete(token).Error
}

func (repo PersonalAccessTokenRepository) Create(token *models.PersonalAccessToken) error {
	return repo.db.Create(token).Error
}

func (repo PersonalAccessTokenRepository) DeleteByUserID(userID uint) error {
	return repo.db.Where("tokenable_type = ? AND tokenable_id = ?", "App\\Models\\User", userID).Delete(&models.PersonalAccessToken{}).Error
}

func (repo PersonalAccessTokenRepository) UpdateFields(token *models.PersonalAccessToken, fields map[string]interface{}) error {
	return repo.db.Model(token).Updates(fields).Error
}
