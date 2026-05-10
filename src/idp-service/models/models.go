package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int       `db:"id"`
	UserUid      uuid.UUID `db:"user_uid"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullName"`
}

type UserResponse struct {
	UserUid  uuid.UUID `json:"userUid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
	Role     string    `json:"role"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" binding:"required"`
	Code         string `json:"code" form:"code"`
	RedirectUri  string `json:"redirect_uri" form:"redirect_uri"`
	ClientId     string `json:"client_id" form:"client_id" binding:"required"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	Username     string `json:"username" form:"username"`
	Password     string `json:"password" form:"password"`
	Scope        string `json:"scope" form:"scope"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
}

type AuthCode struct {
	ID          int       `db:"id"`
	Code        string    `db:"code"`
	UserUid     uuid.UUID `db:"user_uid"`
	ClientId    string    `db:"client_id"`
	RedirectUri string    `db:"redirect_uri"`
	Scope       string    `db:"scope"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	Token     string    `db:"token"`
	UserUid   uuid.UUID `db:"user_uid"`
	ClientId  string    `db:"client_id"`
	Scope     string    `db:"scope"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type JWTClaims struct {
	UserUid  uuid.UUID `json:"sub"`
	Username string    `json:"preferred_username"`
	Email    string    `json:"email"`
	FullName string    `json:"name"`
	Role     string    `json:"role"`
	Scope    string    `json:"scope"`
	Iat      int64     `json:"iat"`
	Exp      int64     `json:"exp"`
	Iss      string    `json:"iss"`
	Aud      string    `json:"aud"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type AuthorizeRequest struct {
	ClientId     string `form:"client_id" binding:"required"`
	RedirectUri  string `form:"redirect_uri" binding:"required"`
	ResponseType string `form:"response_type" binding:"required"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
	Nonce        string `form:"nonce"`
}
