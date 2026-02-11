package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type UptimeMonitorLogRepository struct {
	db *gorm.DB
}

func NewUptimeMonitorLogRepository(db *gorm.DB) *UptimeMonitorLogRepository {
	return &UptimeMonitorLogRepository{db: db}
}

func (r *UptimeMonitorLogRepository) Count(monitorID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.UptimeMonitorLog{})
	if monitorID != 0 {
		query = query.Where("uptime_monitor_id = ?", monitorID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *UptimeMonitorLogRepository) FindAllPaginated(page int, limit int, monitorID uint) ([]models.UptimeMonitorLog, error) {
	var logs []models.UptimeMonitorLog
	offset := (page - 1) * limit
	query := r.db.Model(&models.UptimeMonitorLog{})
	if monitorID != 0 {
		query = query.Where("uptime_monitor_id = ?", monitorID)
	}
	err := query.Limit(limit).Offset(offset).Order("checked_at DESC").Find(&logs).Error
	return logs, err
}

func (r *UptimeMonitorLogRepository) FindByID(id uint) (*models.UptimeMonitorLog, error) {
	var log models.UptimeMonitorLog
	err := r.db.First(&log, id).Error
	return &log, err
}

func (r *UptimeMonitorLogRepository) Store(log *models.UptimeMonitorLog) error {
	return r.db.Create(log).Error
}

func (r *UptimeMonitorLogRepository) Delete(id uint) error {
	return r.db.Delete(&models.UptimeMonitorLog{}, id).Error
}
