package ctx_tools

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type loggerKey struct{}
type requestIDKey struct{}
type userAgentKey struct{}

type enrichLoggerFunc func(entry *log.Entry) *log.Entry

func PutLogger(ctx context.Context, f enrichLoggerFunc) context.Context {
	return context.WithValue(ctx, loggerKey{}, f)
}

func GetLogger(ctx context.Context, entry *log.Entry) *log.Entry {
	f, ok := ctx.Value(loggerKey{}).(enrichLoggerFunc)
	if !ok {
		return log.Info().Str("logger_origin", "newly created")
	}
	return f(entry)
}

func PutRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		return strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	return requestID
}
