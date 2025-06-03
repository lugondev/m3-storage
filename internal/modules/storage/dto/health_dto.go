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
