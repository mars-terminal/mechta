package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/phuslu/log"

	api "github.com/mars-terminal/mechta/api/gen"
	"github.com/mars-terminal/mechta/internal/server/http/middlewares"
	"github.com/mars-terminal/mechta/internal/server/http/shortener"
	"github.com/mars-terminal/mechta/internal/service"
	"github.com/mars-terminal/mechta/internal/shared/ctx_tools"
)

func NewServer(
	service service.Shortener,
) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) && e.Code != http.StatusInternalServerError {
				return ctx.Status(e.Code).JSON(api.InternalServerError{
					Code:    e.Code,
					Message: e.Message,
				})
			}

			ctx_tools.GetLogger(ctx.UserContext(), log.Error()).Err(err).Msg("unexpected error")
			return ctx.Status(http.StatusInternalServerError).JSON(api.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		},
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Use(
		recoverMiddleware.New(),
		middlewares.NewRequestIDInjector(),
		middlewares.NewLogger(),
		cors.New(cors.Config{
			AllowOrigins:     "http://localhost:8080",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowHeaders:     "",
			AllowCredentials: true,
			ExposeHeaders:    "",
			MaxAge:           144000,
		}),
	)

	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to get swagger: %w", err)
	}
	spec, err := swagger.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec into json: %w", err)
	}

	app.Get("/docs", Docs(string(spec)))

	handlers := shortener.NewHandlers(service)
	api.RegisterHandlers(app.Group("/"), api.NewStrictHandler(handlers, nil))

	return app, nil
}
