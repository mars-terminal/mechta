package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/mars-terminal/mechta/internal/shared/ctx_tools"
)

func NewRequestIDInjector() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(
			ctx_tools.PutRequestID(
				ctx.UserContext(),
				strings.ReplaceAll(uuid.NewString(), "-", ""),
			),
		)
		return ctx.Next()
	}
}
