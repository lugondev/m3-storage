package domain

// MediaType represents the type of media.
type MediaType string

const (
	// Image types
	MediaTypeJPEG MediaType = "image/jpeg"
	MediaTypeJPG  MediaType = "image/jpg" // Often same as jpeg, but good to have
	MediaTypePNG  MediaType = "image/png"
	MediaTypeGIF  MediaType = "image/gif"
	MediaTypeWEBP MediaType = "image/webp"

	// Video types
	MediaTypeMP4  MediaType = "video/mp4"
	MediaTypeAVI  MediaType = "video/avi" // More precisely video/x-msvideo
	MediaTypeMOV  MediaType = "video/quicktime"
	MediaTypeWEBM MediaType = "video/webm"

	// Audio types
	MediaTypeMP3  MediaType = "audio/mpeg"
	MediaTypeWAV  MediaType = "audio/wav"
	MediaTypeOGG  MediaType = "audio/ogg"
	MediaTypeFLAC MediaType = "audio/flac"

	// Text/Document types
	MediaTypeTXT  MediaType = "text/plain"
	MediaTypeMD   MediaType = "text/markdown"
	MediaTypePDF  MediaType = "application/pdf"
	MediaTypeDOC  MediaType = "application/msword"
	MediaTypeDOCX MediaType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
)

// SupportedImageExtensions lists all supported image file extensions.
var SupportedImageExtensions = map[string]MediaType{
	".jpg":  MediaTypeJPEG,
	".jpeg": MediaTypeJPEG,
	".png":  MediaTypePNG,
	".gif":  MediaTypeGIF,
	".webp": MediaTypeWEBP,
}

// SupportedVideoExtensions lists all supported video file extensions.
var SupportedVideoExtensions = map[string]MediaType{
	".mp4":  MediaTypeMP4,
	".avi":  MediaTypeAVI,
	".mov":  MediaTypeMOV,
	".webm": MediaTypeWEBM,
}

// SupportedAudioExtensions lists all supported audio file extensions.
var SupportedAudioExtensions = map[string]MediaType{
	".mp3":  MediaTypeMP3,
	".wav":  MediaTypeWAV,
	".ogg":  MediaTypeOGG,
	".flac": MediaTypeFLAC,
}

// SupportedDocumentExtensions lists all supported document file extensions.
var SupportedDocumentExtensions = map[string]MediaType{
	".txt":  MediaTypeTXT,
	".md":   MediaTypeMD,
	".pdf":  MediaTypePDF,
	".doc":  MediaTypeDOC,
	".docx": MediaTypeDOCX,
}

// IsSupportedExtension checks if the given file extension is supported.
func IsSupportedExtension(extension string) bool {
	if _, ok := SupportedImageExtensions[extension]; ok {
		return true
	}
	if _, ok := SupportedVideoExtensions[extension]; ok {
		return true
	}
	if _, ok := SupportedAudioExtensions[extension]; ok {
		return true
	}
	if _, ok := SupportedDocumentExtensions[extension]; ok {
		return true
	}
	return false
}

// GetMediaTypeFromExtension returns the MediaType for a given extension.
// Returns empty string if not found.
func GetMediaTypeFromExtension(extension string) MediaType {
	if mt, ok := SupportedImageExtensions[extension]; ok {
		return mt
	}
	if mt, ok := SupportedVideoExtensions[extension]; ok {
		return mt
	}
	if mt, ok := SupportedAudioExtensions[extension]; ok {
		return mt
	}
	if mt, ok := SupportedDocumentExtensions[extension]; ok {
		return mt
	}
	return ""
}
