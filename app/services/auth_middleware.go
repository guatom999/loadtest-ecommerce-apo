package services

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/guatom999/ecommerce-product-api/app/utils"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func AuthMiddleware(maker *utils.JWTMaker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing bearer token"})
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			tok, err := maker.Parse(tokenStr)
			if err != nil || !tok.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}
			claims, ok := tok.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid claims"})
			}
			sub, _ := claims["sub"].(string)
			if sub == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing sub"})
			}
			c.Set(string(UserIDKey), sub)
			return next(c)
		}
	}
}

func UserIDFromCtx(c echo.Context) string {
	if v := c.Get(string(UserIDKey)); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
