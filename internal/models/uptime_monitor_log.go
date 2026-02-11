/*
 * Project Name: models
 * File: uptime_monitor_log.go
 * Created Date: Wednesday February 11th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package models

import (
	"time"

	"gorm.io/gorm"
)

type UptimeMonitorLog struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UptimeMonitorID uint           `json:"uptime_monitor_id"`
	StatusCode      int            `json:"status_code"`
	ResponseTimeMs  int            `json:"response_time_ms"`
	IsHealthy       bool           `json:"is_healthy"`
	ErrorMessage    string         `json:"error_message"`
	CheckedAt       time.Time      `json:"checked_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
	UptimeMonitor   *UptimeMonitor `gorm:"foreignKey:UptimeMonitorID" json:"uptime_monitor,omitempty"`
}

func (UptimeMonitorLog) TableName() string {
	return "uptime_monitor_logs"
}
