package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is not valid")
)

type JWTMaker struct {
	Secret string
}

type JWTPayload struct {
	Username  string    `json:"usernamer"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	ID        uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func NewJWTPayload(username string, duration time.Duration) (*JWTPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &JWTPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
func NewJWTMaker(secret string) (JWTMaker, error) {
	if len(secret) < minSecretKeySize {
		return JWTMaker{}, fmt.Errorf("secret key must be of min size %v", minSecretKeySize)
	}
	return JWTMaker{secret}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewJWTPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.Secret))
}

func (maker *JWTMaker) VerifyToken(token string) (*JWTPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.Secret), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {

		return nil, err
	}
	if !jwtToken.Valid {
		return nil, ErrExpiredToken
	}
	claims, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
func (payload *JWTPayload) Validate() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
