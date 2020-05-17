package logger

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

/*
 * trace id相关的常量
 */
const (
	MaxTraceID = 100000000
)

type traceIDKey struct{}

// GetTraceID 获取trace id
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		traceID = GenTraceID()
	}
	return traceID
}

// GenTraceID 生成trace id
func GenTraceID() string {
	now := time.Now()
	traceID := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%08d", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), rand.Int31n(MaxTraceID))
	return traceID
}

// WithTraceID Context传递trace id
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}
