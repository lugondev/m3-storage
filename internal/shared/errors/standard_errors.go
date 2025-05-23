package errors

import "net/http"

// Code constants
const (
	CodeBadRequest   = http.StatusBadRequest
	CodeUnauthorized = http.StatusUnauthorized
	CodeForbidden    = http.StatusForbidden
	CodeNotFound     = http.StatusNotFound
	CodeConflict     = http.StatusConflict
	CodeInternal     = http.StatusInternalServerError
)

// Standard errors that can be reused across the application
var (
	// Common errors
	ErrInvalidInput   = NewError(http.StatusBadRequest, "invalid_input")
	ErrUnauthorized   = NewError(http.StatusUnauthorized, "unauthorized")
	ErrForbidden      = NewError(http.StatusForbidden, "forbidden")
	ErrNotFound       = NewError(http.StatusNotFound, "not_found")
	ErrConflict       = NewError(http.StatusConflict, "conflict")
	ErrInternalServer = NewError(http.StatusInternalServerError, "internal_server_error")
	ErrValidation     = NewError(http.StatusBadRequest, "validation_error")

	// Tenant-specific errors
	ErrTenantNotFound      = NewError(http.StatusNotFound, "tenant_not_found")
	ErrTenantSlugExists    = NewError(http.StatusConflict, "tenant_slug_exists")
	ErrOwnerNotFound       = NewError(http.StatusNotFound, "owner_not_found")
	ErrNotTenantOwner      = NewError(http.StatusForbidden, "not_tenant_owner")
	ErrOwnerUserIDRequired = NewError(http.StatusBadRequest, "owner_user_id_required")
	ErrNewOwnerNotFound    = NewError(http.StatusNotFound, "new_owner_not_found")
	ErrNewOwnerNotActive   = NewError(http.StatusBadRequest, "new_owner_not_active")
	ErrTransferToSelf      = NewError(http.StatusBadRequest, "transfer_to_self")
	ErrTenantForbidden     = NewError(http.StatusForbidden, "tenant_forbidden")

	// User-related errors
	ErrUserNotFound              = NewError(http.StatusNotFound, "user_not_found")
	ErrUserNotActive             = NewError(http.StatusBadRequest, "user_not_active")
	ErrUserAlreadyMember         = NewError(http.StatusConflict, "user_already_member")
	ErrUserNotMemberOfTenant     = NewError(http.StatusNotFound, "user_not_member")
	ErrUserNotInTenant           = NewError(http.StatusNotFound, "user_not_in_tenant")
	ErrNewOwnerNotActiveInTenant = NewError(http.StatusBadRequest, "new_owner_not_active_in_tenant")

	// Role and permission errors
	ErrRoleNotFound                        = NewError(http.StatusNotFound, "role_not_found")
	ErrInvalidRole                         = NewError(http.StatusBadRequest, "invalid_role")
	ErrCannotRemoveTenantDesignatedOwner   = NewError(http.StatusForbidden, "cannot_remove_tenant_owner")
	ErrCannotRemoveUserWithTenantOwnerRole = NewError(http.StatusForbidden, "cannot_remove_tenant_owner_role")

	// Profile and user status errors
	ErrProfileNotFound      = NewError(http.StatusNotFound, "profile_not_found")
	ErrEmailTaken           = NewError(http.StatusConflict, "email_taken")
	ErrEmailAlreadyVerified = NewError(http.StatusBadRequest, "email_already_verified")
	ErrEmailNotVerified     = NewError(http.StatusBadRequest, "email_not_verified")
	ErrAccountNotActive     = NewError(http.StatusForbidden, "account_not_active")
	ErrInvalidTwoFactorCode = NewError(http.StatusBadRequest, "invalid_2fa_code")

	// Tenant status errors
	ErrTenantInactive             = NewError(http.StatusForbidden, "tenant_inactive")
	ErrUserInactiveInTenant       = NewError(http.StatusForbidden, "user_inactive_in_tenant")
	ErrTenantJoinRequiresApproval = NewError(http.StatusForbidden, "tenant_join_requires_approval")

	// Authentication errors
	ErrInvalidCredentials = NewError(http.StatusUnauthorized, "invalid_credentials")
	ErrTokenExpired       = NewError(http.StatusUnauthorized, "token_expired")
	ErrInvalidToken       = NewError(http.StatusUnauthorized, "invalid_token", "Invalid token")
)
