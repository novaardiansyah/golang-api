package repositories

import (
	"golang-api/internal/models"

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

func (r *UptimeMonitorRepository) Delete(id uint) error {
	return r.db.Delete(&models.UptimeMonitor{}, id).Error
}
