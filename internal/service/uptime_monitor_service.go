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
	s.updateMonitorStats(monitor, result)
	s.monitorRepo.Update(monitor)

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

func (s *UptimeMonitorService) updateMonitorStats(monitor *models.UptimeMonitor, result checkResult) {
	now := time.Now()

	monitor.Status = result.Status
	monitor.TotalChecks = monitor.TotalChecks + 1
	monitor.LastCheckedAt = &now

	nextCheck := now.Add(time.Duration(monitor.Interval) * time.Second)
	monitor.NextCheckAt = &nextCheck

	if result.IsHealthy {
		monitor.HealthyChecks = monitor.HealthyChecks + 1
		monitor.LastHealthyAt = &now
	} else {
		monitor.UnhealthyChecks = monitor.UnhealthyChecks + 1
		monitor.LastUnhealthyAt = &now
	}

	if result.Status == StatusSlow {
		monitor.LastUnhealthyAt = &now
	}
}

func (s *UptimeMonitorService) RunScheduledChecks() map[string]int {
	results := map[string]int{
		"total":     0,
		"healthy":   0,
		"unhealthy": 0,
	}

	monitors, err := s.monitorRepo.FindDueForCheck()
	if err != nil {
		return results
	}

	for i := range monitors {
		results["total"]++
		if s.Check(&monitors[i]) {
			results["healthy"]++
		} else {
			results["unhealthy"]++
		}
	}

	return results
}
