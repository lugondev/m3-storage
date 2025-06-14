package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	logger "github.com/lugondev/go-log" // Import custom logger
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/media/domain"
	"github.com/lugondev/m3-storage/internal/modules/media/port"
	"github.com/lugondev/m3-storage/internal/presentation/http/fiber/middleware"
	"github.com/lugondev/m3-storage/internal/shared/errors"
	"github.com/lugondev/m3-storage/internal/shared/utils"
)

var _ domain.Media

type MediaHandler struct {
	logger       logger.Logger
	mediaService port.MediaService
	config       *config.Config
}

// NewMediaHandler creates a new MediaHandler.
func NewMediaHandler(appLogger logger.Logger, mediaService port.MediaService, cfg *config.Config) *MediaHandler {
	return &MediaHandler{
		logger:       appLogger.WithFields(map[string]any{"component": "MediaHandler"}),
		mediaService: mediaService,
		config:       cfg,
	}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to the specified provider with optional media type hint
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param provider formData string false "Storage provider (e.g., s3, azure, firebase, discord). If not specified, default provider will be used."
// @Param media_type formData string false "Media type hint (e.g., image/jpeg, video/mp4). If not specified, it will be determined from the file."
// @Failure default {object} errors.Error
// @Router /media/upload [post]
func (h *MediaHandler) UploadFile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get userID from claims", map[string]any{"error": err})
		return err
	}

	h.logger.Info(c.Context(), "Handling file upload request", map[string]any{"userID": userID.String()})

	// 1. Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get file from form", map[string]any{"error": err})
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get file from form: " + err.Error(),
		})
	}

	// 2. Get optional provider and media_type from form
	providerName := c.FormValue("provider")    // Empty if not provided, service will use default
	mediaTypeHint := c.FormValue("media_type") // Empty if not provided, service will attempt to determine

	h.logger.Info(c.Context(), "Upload parameters", map[string]any{
		"fileName":      fileHeader.Filename,
		"fileSize":      fileHeader.Size,
		"providerName":  providerName,
		"mediaTypeHint": mediaTypeHint,
	})

	// 3. Call the media service to upload the file
	// Pass c.Context() for the context.Context parameter
	mediaEntity, err := h.mediaService.UploadFile(c.Context(), userID, fileHeader, providerName, mediaTypeHint)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to upload file via media service", map[string]any{"error": err})
		// Consider more specific error handling based on err type if needed
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to upload file: %v", err),
		})
	}

	h.logger.Info(c.Context(), "File uploaded successfully", map[string]any{"mediaID": mediaEntity.ID.String(), "publicURL": mediaEntity.PublicURL})

	// 4. Return the public URL or other relevant metadata
	return c.Status(http.StatusOK).JSON(mediaEntity)
}

// ListMedia godoc
// @Summary List media files for the authenticated user with pagination
// @Description Get a paginated list of media files owned by the authenticated user
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Paginated list of media files"
// @Failure default {object} errors.Error
// @Router /media [get]
func (h *MediaHandler) ListMedia(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get userID from claims", map[string]any{"error": err})
		return err
	}

	// Parse pagination query
	paginationQuery := &utils.PaginationQuery{}
	if err = c.QueryParser(paginationQuery); err != nil {
		h.logger.Warn(c.Context(), "Failed to parse pagination query", map[string]any{"error": err})
		return errors.ErrInvalidInput
	}

	// Get paginated media files
	pagination, mediaFiles, err := h.mediaService.ListMedia(c.Context(), userID, paginationQuery)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to list media files", map[string]any{"error": err})
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to list media files: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"pagination": pagination,
		"data":       mediaFiles,
	})
}

// GetMedia godoc
// @Summary Get a specific media file
// @Description Get details of a specific media file by ID
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Param id path string true "Media ID"
// @Success 200 {object} domain.Media "Media file details"
// @Failure default {object} errors.Error
// @Router /media/{id} [get]
func (h *MediaHandler) GetMedia(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get userID from claims", map[string]any{"error": err})
		return err
	}

	// Parse media ID from path parameter
	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Warn(c.Context(), "Invalid media ID format", map[string]any{"mediaID": c.Params("id")})
		return errors.ErrInvalidInput
	}

	// Get the media file
	media, err := h.mediaService.GetMedia(c.Context(), userID, mediaID)
	if err != nil {
		if err.Error() == "media file not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Media file not found",
			})
		}
		h.logger.Error(c.Context(), "Failed to get media file", map[string]any{"error": err})
		return err
	}

	return c.Status(http.StatusOK).JSON(media)
}

// DeleteMedia godoc
// @Summary Delete a specific media file
// @Description Delete a specific media file by ID
// @Tags Media
// @Security BearerAuth
// @Param id path string true "Media ID"
// @Success 200 {object} map[string]string "Media file deleted successfully"
// @Failure default {object} errors.Error
// @Router /media/{id} [delete]
func (h *MediaHandler) DeleteMedia(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get userID from claims", map[string]any{"error": err})
		return err
	}

	// Parse media ID from path parameter
	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Warn(c.Context(), "Invalid media ID format", map[string]any{"mediaID": c.Params("id")})
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid media ID format",
		})
	}

	// Delete the media file
	if err := h.mediaService.DeleteMedia(c.Context(), userID, mediaID); err != nil {
		if err.Error() == "media file not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Media file not found",
			})
		}
		h.logger.Error(c.Context(), "Failed to delete media file", map[string]any{"error": err})
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete media file: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Media file deleted successfully",
	})
}

// ServeLocalFile godoc
// @Summary Serve a local media file
// @Description Serve a local media file by ID for authenticated users
// @Tags Media
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "Media ID"
// @Success 200 {file} file "Media file content"
// @Failure 404 {object} fiber.Map "Media file not found"
// @Failure 403 {object} fiber.Map "Access denied"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /media/{id}/file [get]
func (h *MediaHandler) ServeLocalFile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get userID from claims", map[string]any{"error": err})
		return err
	}

	// Parse media ID from path parameter
	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Warn(c.Context(), "Invalid media ID format", map[string]any{"mediaID": c.Params("id")})
		return errors.ErrInvalidInput
	}

	// Get the media file
	media, err := h.mediaService.GetMedia(c.Context(), userID, mediaID)
	if err != nil {
		if err.Error() == "media file not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Media file not found",
			})
		}
		h.logger.Error(c.Context(), "Failed to get media file", map[string]any{"error": err})
		return err
	}

	// Check if this is a local storage file
	if media.Provider != "local" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "This endpoint only serves local storage files",
		})
	}

	// Serve the file directly from the local file system
	return c.SendFile(h.config.LocalStorage.Path + "/" + media.FilePath)
}

// ServePublicLocalFile godoc
// @Summary Serve a public local media file
// @Description Serve a local media file without authentication
// @Tags Media
// @Produce application/octet-stream
// @Param id path string true "Media ID"
// @Success 200 {file} file "Media file content"
// @Failure 404 {object} fiber.Map "Media file not found"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /media/public/{id}/file [get]
func (h *MediaHandler) ServePublicLocalFile(c *fiber.Ctx) error {
	// Parse media ID from path parameter
	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Warn(c.Context(), "Invalid media ID format", map[string]any{"mediaID": c.Params("id")})
		return errors.ErrInvalidInput
	}

	// Get the media file without user authentication
	media, err := h.mediaService.GetPublicMedia(c.Context(), mediaID)
	if err != nil {
		if err.Error() == "media file not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Media file not found",
			})
		}
		h.logger.Error(c.Context(), "Failed to get public media file", map[string]any{"error": err})
		return err
	}

	// Check if this is a local storage file
	if media.Provider != "local" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "This endpoint only serves local storage files",
		})
	}

	// Serve the file directly from the local file system
	return c.SendFile(h.config.LocalStorage.Path + "/" + media.FilePath)
}
