package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

func (lr *LoginRequest) Validate() error {
	if lr.Username != "" {
		if lr.Password == "" {
			return fmt.Errorf("password is required when username is provided")
		}
		if lr.Token != "" {
			return fmt.Errorf("token must not be provided when username is provided")
		}
		return nil
	}

	if lr.Token != "" {
		if lr.Username != "" || lr.Password != "" {
			return fmt.Errorf("username and password must not be provided when token is provided")
		}
		return nil
	}

	return fmt.Errorf("either username/password or token must be provided")
}

type AuthResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  User      `json:"user"`
}

type RenewAcessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAcessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
