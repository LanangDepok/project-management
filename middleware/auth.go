package middleware

import (
	"strings"

	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/utils"
	"github.com/gofiber/fiber/v3"
	jwt "github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c, "Unauthorized", "missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return utils.Unauthorized(c, "Unauthorized", "invalid format, expected: Bearer <token>")
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.Unauthorized(c, "Unauthorized", "invalid or expired token")
		}

		c.Locals("user", token)
		return c.Next()
	}
}
