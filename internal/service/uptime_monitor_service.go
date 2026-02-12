package service

import (
	"fmt"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const (
	SlowResponseThresholdMs = 300
	HttpTimeoutSeconds      = 30

	StatusUp   = "up"
	StatusDown = "down"
	StatusSlow = "slow"
)

type UptimeMonitorService struct {
	monitorRepo *repositories.UptimeMonitorRepository
	logRepo     *repositories.UptimeMonitorLogRepository
}

func NewUptimeMonitorService(db *gorm.DB) *UptimeMonitorService {
	return &UptimeMonitorService{
		monitorRepo: repositories.NewUptimeMonitorRepository(db),
		logRepo:     repositories.NewUptimeMonitorLogRepository(db),
	}
}

type checkResult struct {
	StatusCode     int
	ResponseTimeMs int
	IsHealthy      bool
	Status         string
	ErrorMessage   string
}

func (s *UptimeMonitorService) Check(monitor *models.UptimeMonitor) bool {
	result := s.performHttpCheck(monitor.URL)
	result = s.evaluateSlowResponse(result)
	result = s.generateErrorMessage(result)

	s.createLog(monitor, result)
	fields := s.updateMonitorStats(monitor, result)
	s.monitorRepo.UpdateFields(monitor.ID, fields)

	return result.IsHealthy
}

func (s *UptimeMonitorService) performHttpCheck(url string) checkResult {
	var statusCode int
	var responseTimeMs int
	isHealthy := false

	client := &http.Client{
		Timeout: HttpTimeoutSeconds * time.Second,
	}

	startTime := time.Now()

	resp, err := client.Get(url)
	responseTimeMs = int(time.Since(startTime).Milliseconds())

	if err != nil {
		statusCode = 408
	} else {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
		isHealthy = statusCode >= 200 && statusCode < 300
	}

	return checkResult{
		StatusCode:     statusCode,
		ResponseTimeMs: responseTimeMs,
		IsHealthy:      isHealthy,
	}
}

func (s *UptimeMonitorService) evaluateSlowResponse(result checkResult) checkResult {
	status := StatusUp

	if result.IsHealthy && result.ResponseTimeMs > SlowResponseThresholdMs {
		status = StatusSlow
	} else if !result.IsHealthy {
		status = StatusDown
	}

	result.Status = status
	return result
}

func (s *UptimeMonitorService) generateErrorMessage(result checkResult) checkResult {
	if result.IsHealthy {
		result.ErrorMessage = ""
		return result
	}

	result.ErrorMessage = fmt.Sprintf("HTTP %d", result.StatusCode)
	return result
}

func (s *UptimeMonitorService) createLog(monitor *models.UptimeMonitor, result checkResult) {
	now := time.Now()
	log := &models.UptimeMonitorLog{
		UptimeMonitorID: monitor.ID,
		StatusCode:      result.StatusCode,
		ResponseTimeMs:  result.ResponseTimeMs,
		IsHealthy:       result.IsHealthy,
		ErrorMessage:    result.ErrorMessage,
		CheckedAt:       now,
	}
	s.logRepo.Store(log)
}

func (s *UptimeMonitorService) updateMonitorStats(monitor *models.UptimeMonitor, result checkResult) map[string]interface{} {
	now := time.Now()
	nextCheck := now.Add(time.Duration(monitor.Interval) * time.Second)

	fields := map[string]interface{}{
		"status":          result.Status,
		"total_checks":    monitor.TotalChecks + 1,
		"last_checked_at": &now,
		"next_check_at":   &nextCheck,
	}

	if result.IsHealthy {
		fields["healthy_checks"] = monitor.HealthyChecks + 1
		fields["last_healthy_at"] = &now
	} else {
		fields["unhealthy_checks"] = monitor.UnhealthyChecks + 1
		fields["last_unhealthy_at"] = &now
	}

	if result.Status == StatusSlow {
		fields["last_unhealthy_at"] = &now
	}

	monitor.Status = result.Status
	monitor.TotalChecks++
	monitor.LastCheckedAt = &now
	monitor.NextCheckAt = &nextCheck

	if result.IsHealthy {
		monitor.HealthyChecks++
		monitor.LastHealthyAt = &now
	} else {
		monitor.UnhealthyChecks++
		monitor.LastUnhealthyAt = &now
	}

	return fields
}

func (s *UptimeMonitorService) RunScheduledChecks() map[string]int {
	results := map[string]int{
		"total":     0,
		"healthy":   0,
		"unhealthy": 0,
	}

	err := s.monitorRepo.ProcessDueForCheck(50, func(monitors []models.UptimeMonitor) error {
		for i := range monitors {
			results["total"]++
			if s.Check(&monitors[i]) {
				results["healthy"]++
			} else {
				results["unhealthy"]++
			}
		}
		return nil
	})

	if err != nil {
		// Log error if needed
	}

	return results
}
