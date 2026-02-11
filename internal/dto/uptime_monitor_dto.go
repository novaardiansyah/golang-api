package dto

import "time"

type StoreUptimeMonitorRequest struct {
	Code     string `json:"code" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	Name     string `json:"name" validate:"required"`
	Interval int    `json:"interval" validate:"required,numeric"`
	IsActive bool   `json:"is_active"`
}

type UpdateUptimeMonitorRequest struct {
	Code     string `json:"code"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Interval int    `json:"interval"`
	IsActive bool   `json:"is_active"`
}

type StoreUptimeMonitorLogRequest struct {
	UptimeMonitorID uint      `json:"uptime_monitor_id" validate:"required"`
	StatusCode      int       `json:"status_code" validate:"required"`
	ResponseTimeMs  int       `json:"response_time_ms" validate:"required"`
	IsHealthy       bool      `json:"is_healthy"`
	ErrorMessage    string    `json:"error_message"`
	CheckedAt       time.Time `json:"checked_at" validate:"required"`
}
