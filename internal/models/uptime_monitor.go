/*
 * Project Name: models
 * File: uptime_monitor.go
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

type UptimeMonitor struct {
	ID              uint               `gorm:"primaryKey" json:"id"`
	Code            string             `json:"code"`
	URL             string             `json:"url"`
	Name            string             `json:"name"`
	Interval        int                `json:"interval"`
	IsActive        bool               `json:"is_active"`
	Status          string             `json:"status"`
	LastCheckedAt   *time.Time         `json:"last_checked_at"`
	LastHealthyAt   *time.Time         `json:"last_healthy_at"`
	LastUnhealthyAt *time.Time         `json:"last_unhealthy_at"`
	TotalChecks     int                `json:"total_checks"`
	HealthyChecks   int                `json:"healthy_checks"`
	UnhealthyChecks int                `json:"unhealthy_checks"`
	NextCheckAt     *time.Time         `json:"next_check_at"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	DeletedAt       gorm.DeletedAt     `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
	Logs            []UptimeMonitorLog `gorm:"foreignKey:UptimeMonitorID" json:"logs,omitempty"`
}

func (UptimeMonitor) TableName() string {
	return "uptime_monitors"
}
