package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	logger "github.com/lugondev/go-log" // Import custom logger
	"github.com/lugondev/m3-storage/internal/modules/media/port"
	"go.uber.org/zap"
)

type MediaHandler struct {
	logger       logger.Logger // Keep as *zap.Logger for internal use
	mediaService port.MediaService
}

// NewMediaHandler creates a new MediaHandler.
func NewMediaHandler(appLogger logger.Logger, mediaService port.MediaService) *MediaHandler {
	return &MediaHandler{
		logger:       appLogger,
		mediaService: mediaService,
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
// @Success 200 {object} map[string]interface{} "File uploaded successfully with media details"
// @Failure 400 {object} map[string]string "Bad request - missing file or invalid parameters"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /media/upload [post]
func (h *MediaHandler) UploadFile(c *fiber.Ctx) error {
	// 0. Extract UserID from JWT token
	//claims, ok := c.Locals(constants.UserClaimsKey).(*jwt.JWTClaims)
	//if !ok || claims == nil {
	//	h.logger.Warn(c.Context(), "User claims not found in context or invalid type")
	//	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
	//		"error": "Unauthorized: user claims not found",
	//	})
	//}
	//userID, err := uuid.Parse(claims.Subject)
	//if err != nil {
	//	h.logger.Error(c.Context(), "Failed to parse userID from claims", zap.Error(err), zap.String("claims.Subject", claims.Subject))
	//	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	//		"error": "Internal server error: could not parse user ID",
	//	})
	//}

	//h.logger.Info(c.Context(), "Handling file upload request", zap.String("userID", userID.String()))

	// 1. Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get file from form", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get file from form: " + err.Error(),
		})
	}

	// 2. Get optional provider and media_type from form
	providerName := c.FormValue("provider")    // Empty if not provided, service will use default
	mediaTypeHint := c.FormValue("media_type") // Empty if not provided, service will attempt to determine

	h.logger.Info(c.Context(), "Upload parameters",
		zap.String("fileName", fileHeader.Filename),
		zap.Int64("fileSize", fileHeader.Size),
		zap.String("providerName", providerName),
		zap.String("mediaTypeHint", mediaTypeHint),
	)

	// 3. Call the media service to upload the file
	// Pass c.Context() for the context.Context parameter
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	mediaEntity, err := h.mediaService.UploadFile(c.Context(), userID, fileHeader, providerName, mediaTypeHint)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to upload file via media service", zap.Error(err))
		// Consider more specific error handling based on err type if needed
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to upload file: %v", err),
		})
	}

	h.logger.Info(c.Context(), "File uploaded successfully", zap.String("mediaID", mediaEntity.ID.String()), zap.String("publicURL", mediaEntity.PublicURL))

	// 4. Return the public URL or other relevant metadata
	return c.Status(http.StatusOK).JSON(mediaEntity)
}
