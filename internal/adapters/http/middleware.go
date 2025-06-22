package http

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(sec string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(sec),
			JWTAlg: jwt.SigningMethodHS256.Name,
		},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	})
}

func RequestIdMiddleware() fiber.Handler {
	return requestid.New()
}

func LoggerMiddleware() fiber.Handler {
	return fiberlogger.New(fiberlogger.Config{
		Format: "${time} ${pid} ${locals:requestid} ${status} - ${method} ${path} | ${latency}\n",
	})
}
