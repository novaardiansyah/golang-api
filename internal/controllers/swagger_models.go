package controllers

import "time"

type UserSwagger struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	HasAllowNotification *bool     `json:"has_allow_notification"`
	NotificationToken    *string   `json:"notification_token,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeletedAt            *string   `json:"deleted_at,omitempty"`
}

type AccountInfoSwagger struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}

type PaymentSwagger struct {
	ID                 uint                `json:"id"`
	UserID             uint                `json:"user_id"`
	Code               string              `json:"code"`
	Name               string              `json:"name"`
	Date               string              `json:"date"`
	Amount             int64               `json:"amount"`
	HasItems           bool                `json:"has_items"`
	IsScheduled        bool                `json:"is_scheduled"`
	IsDraft            bool                `json:"is_draft"`
	TypeID             uint                `json:"type_id"`
	PaymentAccountID   *uint               `json:"payment_account_id"`
	PaymentAccountToID *uint               `json:"payment_account_to_id"`
	UpdatedAt          time.Time           `json:"updated_at"`
	Type               string              `json:"type"`
	FormattedAmount    string              `json:"formatted_amount"`
	FormattedDate      string              `json:"formatted_date"`
	FormattedUpdatedAt string              `json:"formatted_updated_at"`
	AttachmentsCount   int                 `json:"attachments_count"`
	ItemsCount         int                 `json:"items_count"`
	Account            *AccountInfoSwagger `json:"account"`
	AccountTo          *AccountInfoSwagger `json:"account_to"`
}

type FileSwagger struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	FileName    string `json:"file_name"`
	FileSize    string `json:"file_size"`
	FileAlias   string `json:"file_alias"`
	DownloadURL string `json:"download_url"`
}

type FormattedGoalSwagger struct {
	Amount       string `json:"amount"`
	TargetAmount string `json:"target_amount"`
	Progress     string `json:"progress"`
	StartDate    string `json:"start_date"`
	TargetDate   string `json:"target_date"`
}

type PaymentGoalSwagger struct {
	ID              uint                 `json:"id"`
	UserID          uint                 `json:"user_id"`
	StatusID        uint                 `json:"status_id"`
	Code            string               `json:"code"`
	Name            string               `json:"name"`
	Description     *string              `json:"description"`
	Amount          int64                `json:"amount"`
	TargetAmount    int64                `json:"target_amount"`
	ProgressPercent int                  `json:"progress_percent"`
	StartDate       string               `json:"start_date"`
	TargetDate      string               `json:"target_date"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Status          string               `json:"status"`
	Formatted       FormattedGoalSwagger `json:"formatted"`
}

type FormattedAccountSwagger struct {
	Deposit string `json:"deposit"`
	Logo    string `json:"logo"`
}

type PaymentAccountSwagger struct {
	ID        uint                    `json:"id"`
	Name      string                  `json:"name"`
	Deposit   int64                   `json:"deposit"`
	Logo      string                  `json:"logo"`
	Formatted FormattedAccountSwagger `json:"formatted"`
}

type ActivityLogSwagger struct {
	ID             uint   `json:"id"`
	LogName        string `json:"log_name"`
	Description    string `json:"description"`
	SubjectID      uint   `json:"subject_id"`
	SubjectType    string `json:"subject_type"`
	Event          string `json:"event"`
	CauserID       uint   `json:"causer_id"`
	CauserType     string `json:"causer_type"`
	PrevProperties string `json:"prev_properties"`
	Properties     string `json:"properties"`
	BatchUUID      string `json:"batch_uuid"`
	IPAddress      string `json:"ip_address"`
	Country        string `json:"country"`
	City           string `json:"city"`
	Region         string `json:"region"`
	Postal         string `json:"postal"`
	Geolocation    string `json:"geolocation"`
	Timezone       string `json:"timezone"`
	UserAgent      string `json:"user_agent"`
	Referer        string `json:"referer"`
}
