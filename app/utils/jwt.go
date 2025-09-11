package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTMaker() *JWTMaker {
	secret := []byte(MustGetenv("JWT_SECRET"))
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "hexshop"
	}
	hours := 24
	if v := os.Getenv("JWT_EXPIRE_HOURS"); v != "" {
		// ignore parse error, keep default
		fmtHours, _ := time.ParseDuration(v + "h")
		if fmtHours > 0 {
			hours = int(fmtHours.Hours())
		}
	}
	return &JWTMaker{secret: secret, issuer: issuer, ttl: time.Duration(hours) * time.Hour}
}

func (m *JWTMaker) Create(userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": m.issuer,
		"iat": now.Unix(),
		"exp": now.Add(m.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTMaker) Parse(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	})
}
