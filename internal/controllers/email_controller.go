package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type EmailController struct {
	repo *repositories.EmailRepository
}

func NewEmailController(repo *repositories.EmailRepository) *EmailController {
	return &EmailController{repo: repo}
}

type EmailAttachmentResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	FileName    string `json:"file_name"`
	FileSize    string `json:"file_size"`
	DownloadURL string `json:"download_url"`
}

// GetAttachments godoc
// @Summary Get email attachments
// @Description Get list of attachments for a specific email by UID
// @Tags emails
// @Accept json
// @Produce json
// @Param uid path string true "Email UID"
// @Success 200 {object} utils.Response{data=[]EmailAttachmentSwagger}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /emails/{uid}/attachments [get]
// @Security BearerAuth
func (ctrl *EmailController) GetAttachments(c *fiber.Ctx) error {
	uid := c.Params("uid")

	email, err := ctrl.repo.GetByUID(uid)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Email not found")
	}

	attachments := make([]EmailAttachmentResponse, 0, len(email.Files))
	for _, file := range email.Files {
		attachments = append(attachments, EmailAttachmentResponse{
			ID:          file.ID,
			Code:        file.Code,
			FileName:    file.FileName,
			FileSize:    utils.FormatFileSize(file.FileSize),
			DownloadURL: file.DownloadURL,
		})
	}

	return utils.SuccessResponse(c, "Email attachments retrieved successfully", attachments)
}
