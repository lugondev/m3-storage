package dto

// HealthCheckRequest represents the request for checking storage provider health
type HealthCheckRequest struct {
	ProviderType string `json:"provider_type" validate:"required" example:"s3"`
}

// HealthCheckResponse represents the response for a single provider health check
type HealthCheckResponse struct {
	Status  string `json:"status" example:"healthy"`
	Message string `json:"message,omitempty" example:""`
}

// HealthCheckAllResponse represents the response for all providers health check
type HealthCheckAllResponse struct {
	Providers map[string]HealthCheckResponse `json:"providers"`
}

// ProviderInfo represents information about a storage provider
type ProviderInfo struct {
	Type        string `json:"type" example:"s3"`
	Name        string `json:"name" example:"Amazon S3"`
	Description string `json:"description" example:"Amazon Simple Storage Service"`
}

// ListProvidersResponse represents the response for listing available providers
type ListProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}
