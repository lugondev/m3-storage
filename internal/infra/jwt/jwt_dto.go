package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the standard JWT claims plus custom ones.
// Uses jwt.RegisteredClaims for standard fields.
type JWTClaims struct {
	Email       string   `json:"email"`
	Roles       []string `json:"roles,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}
