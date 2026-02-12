package repositories

import (
	"golang-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type UptimeMonitorRepository struct {
	db *gorm.DB
}

func NewUptimeMonitorRepository(db *gorm.DB) *UptimeMonitorRepository {
	return &UptimeMonitorRepository{db: db}
}

func (r *UptimeMonitorRepository) Count(search string) (int64, error) {
	var count int64
	err := r.db.Model(&models.UptimeMonitor{}).
		Where("name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&count).Error
	return count, err
}

func (r *UptimeMonitorRepository) FindAllPaginated(page int, limit int, search string) ([]models.UptimeMonitor, error) {
	var monitors []models.UptimeMonitor
	offset := (page - 1) * limit

	query := r.db.Model(&models.UptimeMonitor{})
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(&monitors).Error
	return monitors, err
}

func (r *UptimeMonitorRepository) FindByID(id uint) (*models.UptimeMonitor, error) {
	var monitor models.UptimeMonitor
	err := r.db.First(&monitor, id).Error
	return &monitor, err
}

func (r *UptimeMonitorRepository) Store(monitor *models.UptimeMonitor) error {
	return r.db.Create(monitor).Error
}

func (r *UptimeMonitorRepository) Update(monitor *models.UptimeMonitor) error {
	return r.db.Save(monitor).Error
}

func (r *UptimeMonitorRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.UptimeMonitor{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UptimeMonitorRepository) Delete(id uint) error {
	return r.db.Delete(&models.UptimeMonitor{}, id).Error
}

func (r *UptimeMonitorRepository) FindDueForCheck() ([]models.UptimeMonitor, error) {
	var monitors []models.UptimeMonitor
	now := time.Now()
	err := r.db.Where("is_active = ?", true).
		Where("next_check_at IS NULL OR next_check_at <= ?", now).
		Find(&monitors).Error
	return monitors, err
}

func (r *UptimeMonitorRepository) ProcessDueForCheck(batchSize int, callback func(monitors []models.UptimeMonitor) error) error {
	now := time.Now()
	return r.db.Model(&models.UptimeMonitor{}).
		Where("is_active = ?", true).
		Where("next_check_at IS NULL OR next_check_at <= ?", now).
		FindInBatches(&[]models.UptimeMonitor{}, batchSize, func(tx *gorm.DB, batch int) error {
			var monitors []models.UptimeMonitor
			if err := tx.Find(&monitors).Error; err != nil {
				return err
			}
			return callback(monitors)
		}).Error
}
