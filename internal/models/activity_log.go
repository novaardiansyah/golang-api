/*
 * Project Name: models
 * File: activity_log.go
 * Created Date: Tuesday February 10th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type ActivityLog struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	LogName        string           `json:"log_name"`
	Description    string           `json:"description"`
	SubjectID      *uint            `json:"subject_id"`
	SubjectType    *string          `json:"subject_type"`
	Event          string           `json:"event"`
	CauserID       uint             `json:"causer_id"`
	CauserType     string           `json:"causer_type"`
	PrevProperties *json.RawMessage `gorm:"type:json" json:"prev_properties"`
	Properties     json.RawMessage  `gorm:"type:json" json:"properties"`
	BatchUUID      *string          `json:"batch_uuid"`
	IPAddress      *string          `json:"ip_address"`
	Country        *string          `json:"country"`
	City           *string          `json:"city"`
	Region         *string          `json:"region"`
	Postal         *string          `json:"postal"`
	Geolocation    *string          `json:"geolocation"`
	Timezone       *string          `json:"timezone"`
	UserAgent      *string          `json:"user_agent"`
	Referer        *string          `json:"referer"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
