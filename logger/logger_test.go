package logger

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	Start(WithType("file"))
	Info(ctx, "测试呀测试呀")
	Access(ctx, "可以的%s", "说")
	Warn(ctx, "warn....")
	Stop()
}
