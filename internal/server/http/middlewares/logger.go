package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/mars-terminal/mechta/internal/shared/ctx_tools"
)

func mapStatusToMessage(status int) string {
	result := http.StatusText(status)
	if result != "" {
		return strings.ToLower(result)
	}
	return strconv.Itoa(status)
}

func NewLogger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()
		path, method := ctx.Path(), ctx.Method()

		entry := log.Info().
			Str("path", path).
			Str("method", method).
			Time("start_at", start).
			Str("request_id", ctx_tools.GetRequestID(ctx.UserContext()))

		ctx.SetUserContext(ctx_tools.PutLogger(
			ctx.UserContext(),
			func(entry *log.Entry) *log.Entry {
				return entry.
					Str("path", path).
					Str("method", method).
					Time("start_at", start).
					Str("request_id", ctx_tools.GetRequestID(ctx.UserContext()))
			},
		))
		ctx.Set("REQUEST-ID", ctx_tools.GetRequestID(ctx.UserContext()))

		if err := ctx.Next(); err != nil {
			if err := ctx.App().ErrorHandler(ctx, err); err != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
			}
			entry.Err(err)
		}

		entry.
			Str("user_agent", string(ctx.Request().Header.UserAgent())).
			Str("ip", ctx.IP()).
			Str("latency", time.Now().Sub(start).String()).
			Int("status", ctx.Response().StatusCode()).
			Msg(mapStatusToMessage(ctx.Response().StatusCode()))

		return nil
	}
}
