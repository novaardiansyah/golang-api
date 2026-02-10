package dto

import "encoding/json"

type StoreActivityLogRequest struct {
	LogName        string          `json:"log_name"`
	Description    string          `json:"description"`
	SubjectID      uint            `json:"subject_id"`
	SubjectType    string          `json:"subject_type"`
	Event          string          `json:"event"`
	CauserID       uint            `json:"causer_id"`
	CauserType     string          `json:"causer_type"`
	PrevProperties json.RawMessage `json:"prev_properties"`
	Properties     json.RawMessage `json:"properties"`
	BatchUUID      string          `json:"batch_uuid"`
	IPAddress      string          `json:"ip_address"`
	Country        string          `json:"country"`
	City           string          `json:"city"`
	Region         string          `json:"region"`
	Postal         string          `json:"postal"`
	Geolocation    string          `json:"geolocation"`
	Timezone       string          `json:"timezone"`
	UserAgent      string          `json:"user_agent"`
	Referer        string          `json:"referer"`
}
