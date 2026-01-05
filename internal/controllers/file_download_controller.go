package controllers

import (
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type FileDownloadController struct {
	repo *repositories.FileDownloadRepository
}

func NewFileDownloadController(repo *repositories.FileDownloadRepository) *FileDownloadController {
	return &FileDownloadController{repo: repo}
}

type FileResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	FileName    string `json:"file_name"`
	FileSize    string `json:"file_size"`
	DownloadURL string `json:"download_url"`
}

// GetFiles godoc
// @Summary Get files by file download UID
// @Description Get list of files for a specific file download by UID
// @Tags files
// @Accept json
// @Produce json
// @Param uid path string true "File Download UID"
// @Success 200 {object} utils.Response{data=[]FileSwagger}
// @Failure 404 {object} utils.Response
// @Router /files/d/{uid} [get]
// @Security BearerAuth
func (ctrl *FileDownloadController) GetFiles(c *fiber.Ctx) error {
	uid := c.Params("uid")

	fileDownload, err := ctrl.repo.GetByUID(uid)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "File download not found")
	}

	files := make([]FileResponse, 0, len(fileDownload.Files))
	for _, file := range fileDownload.Files {
		files = append(files, FileResponse{
			ID:          file.ID,
			Code:        file.Code,
			FileName:    file.FileName,
			FileSize:    utils.FormatFileSize(file.FileSize),
			DownloadURL: file.DownloadURL,
		})
	}

	return utils.SuccessResponse(c, "Files retrieved successfully", files)
}
