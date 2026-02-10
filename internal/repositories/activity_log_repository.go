package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Count(search string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActivityLog{}).Where("log_name LIKE ?", "%"+search+"%").Count(&count).Error
	return count, err
}

func (r *ActivityLogRepository) FindAllPaginated(page int, limit int, search string) ([]models.ActivityLog, error) {
	var activityLogs []models.ActivityLog

	offset := (page - 1) * limit

	if search != "" {
		r.db = r.db.Where("log_name LIKE ?", "%"+search+"%")
	}

	err := r.db.Limit(limit).Offset(offset).Find(&activityLogs).Error

	return activityLogs, err
}

func (r *ActivityLogRepository) Store(activityLog *models.ActivityLog) error {
	return r.db.Create(activityLog).Error
}
