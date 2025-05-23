package middleware

import (
	"fmt"
	"strings"
	"time"

	customLogger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/shared/errors"

	"golang.org/x/text/language"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/gofiber/contrib/otelfiber/v2"
)

// SetupMiddleware configures all middleware for Fiber
func SetupMiddleware(app *fiber.App, cfg config.Config, i18nBundle *i18n.Bundle, log customLogger.Logger) {
	// Request ID middleware
	app.Use(requestid.New())
	app.Use(otelfiber.Middleware())

	// i18n middleware
	app.Use(I18nMiddleware(i18nBundle))

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Ho_Chi_Minh",
	}))

	// OpenTelemetry tracing middleware
	app.Use(func(c *fiber.Ctx) error {
		if desired, _, _ := language.ParseAcceptLanguage(c.Get(fiber.HeaderAcceptLanguage)); len(desired) > 0 {
			if len(desired) > 1 {
				c.Locals("lang", desired[1].String())
			}
			c.Locals("lang", strings.Split(desired[0].String(), "-")[0])
		}

		// Skip logging for OPTIONS requests
		if c.Method() != fiber.MethodOptions {
			log.Info(c.UserContext(), "Request processed", map[string]any{
				"request_id": c.Get(fiber.HeaderXRequestID),
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     c.Response().StatusCode(),
				"ip":         c.IP(),
				"user_agent": c.Get(fiber.HeaderUserAgent),
			})
		}

		return c.Next()
	})

	// Recovery middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS middleware

	corsCfg := cors.Config{
		AllowOrigins:     cfg.App.Origins,
		AllowCredentials: cfg.App.Origins != "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID, traceparent",
		ExposeHeaders:    "Authorization",
	}
	app.Use(cors.New(corsCfg))

	// Rate limiter middleware using values from config
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimiter.Max,
		Expiration: time.Duration(cfg.RateLimiter.ExpirationSeconds) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":  "Too many requests",
				"status": "error",
			})
		},
	}))
}

// ErrorHandler returns a custom error handler for Fiber
func ErrorHandler(log customLogger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = TranslatorTranslate(c, e.Message)
		}

		if e, ok := err.(*errors.Error); ok {
			code = e.StatusCode // Use StatusCode instead of Code
			message = TranslatorTranslate(c, fmt.Sprintf("error_%s", e.Code), e.Message)
		}

		log.Error(c.UserContext(), "Request error", map[string]any{
			"error":      err.Error(),
			"request_id": c.Get(fiber.HeaderXRequestID),
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     code,
			"ip":         c.IP(),
			"user_agent": c.Get(fiber.HeaderUserAgent),
		})

		// Build error response
		response := fiber.Map{
			"status":     "error",
			"message":    message,
			"code":       code,
			"request_id": c.Get(fiber.HeaderXRequestID),
		}

		return c.Status(code).JSON(response)
	}
}
