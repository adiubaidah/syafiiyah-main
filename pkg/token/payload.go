package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Payload defines the structure for JWT payload
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
	Audience  string    `json:"aud"`
	Issuer    string    `json:"iss"`
	NotBefore time.Time `json:"nbf"`
	Subject   string    `json:"sub"`
}

// NewPayload creates a new payload
func NewPayload(username, role string, duration time.Duration) (*Payload, error) {
	now := time.Now()
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
		Audience:  "post-service",
		Issuer:    "auth.service",
		NotBefore: now,
		Subject:   "access",
	}
	return payload, nil
}

// GetAudience returns the audience claim
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{payload.Audience}, nil
}

// GetExpirationTime returns the expiration time claim
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.ExpiredAt), nil
}

// GetIssuedAt returns the issued at claim
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}

// GetIssuer returns the issuer claim
func (payload *Payload) GetIssuer() (string, error) {
	return payload.Issuer, nil
}

// GetNotBefore returns the not before claim
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.NotBefore), nil
}

// GetSubject returns the subject claim
func (payload *Payload) GetSubject() (string, error) {
	return payload.Subject, nil
}

// Valid validates the payload claims
func (payload *Payload) Valid() error {
	if err := payload.ValidateAudience(); err != nil {
		return err
	}
	if err := payload.ValidateIssuer(); err != nil {
		return err
	}
	if err := payload.ValidateExpiry(); err != nil {
		return err
	}
	return nil
}

// ValidateAudience checks the audience claim
func (payload *Payload) ValidateAudience() error {
	if len(payload.Audience) == 0 {
		return fmt.Errorf("audience is missing")
	}
	return nil
}

// ValidateIssuer checks the issuer claim
func (payload *Payload) ValidateIssuer() error {
	if len(payload.Issuer) == 0 {
		return fmt.Errorf("issuer is missing")
	}
	return nil
}

// ValidateExpiry checks the expiration time
func (payload *Payload) ValidateExpiry() error {
	if time.Now().After(payload.ExpiredAt) {
		return fmt.Errorf("token has expired")
	}
	return nil
}
