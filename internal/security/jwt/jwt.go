package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret          string
	AccessTokenTTL  string // "15m"
	RefreshTokenTTL string // "720h"
	Issuer          string
}

type Service interface {
	GenerateTokens(userID int64, email string) (access, refresh string, err error)
	Parse(token string) (*jwtv5.Token, jwtv5.MapClaims, error)
}

type service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func New(cfg Config) Service {
	acc, _ := time.ParseDuration(cfg.AccessTokenTTL)
	ref, _ := time.ParseDuration(cfg.RefreshTokenTTL)
	return &service{
		secret:     []byte(cfg.Secret),
		accessTTL:  acc,
		refreshTTL: ref,
		issuer:     cfg.Issuer,
	}
}

func (s *service) GenerateTokens(userID int64, email string) (string, string, error) {
	now := time.Now()

	accessClaims := jwtv5.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   s.issuer,
		"iat":   now.Unix(),
		"exp":   now.Add(s.accessTTL).Unix(),
		"type":  "access",
	}
	accessToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, accessClaims)
	accStr, err := accessToken.SignedString(s.secret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwtv5.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   s.issuer,
		"iat":   now.Unix(),
		"exp":   now.Add(s.refreshTTL).Unix(),
		"type":  "refresh",
	}
	refreshToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, refreshClaims)
	refStr, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return "", "", err
	}

	return accStr, refStr, nil
}

func (s *service) Parse(tokenStr string) (*jwtv5.Token, jwtv5.MapClaims, error) {
	parser := jwtv5.NewParser(jwtv5.WithValidMethods([]string{jwtv5.SigningMethodHS256.Alg()}))
	claims := jwtv5.MapClaims{}
	tkn, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwtv5.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !tkn.Valid {
		return nil, nil, errors.New("invalid token")
	}
	return tkn, claims, nil
}
