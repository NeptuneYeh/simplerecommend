package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has invalid claims: token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p Payload) Valid() error {
	now := time.Now()
	if now.After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (p Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p Payload) GetNotBefore() (*jwt.NumericDate, error) {
	// Assuming the token is valid immediately upon issue
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p Payload) GetIssuer() (string, error) {
	// If your application requires an issuer, set it here
	// For simplicity, this example does not set an issuer
	return "", nil
}

func (p Payload) GetSubject() (string, error) {
	// Here we use the token ID as the JWT subject
	return p.ID.String(), nil
}

func (p Payload) GetAudience() (jwt.ClaimStrings, error) {
	// If your application requires an audience, set it here
	// For simplicity, this example does not set an audience
	return nil, nil
}
