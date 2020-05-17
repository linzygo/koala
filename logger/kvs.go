package logger

import (
	"context"
	"sync"
)

// KeyVal k-v, access日志的参数
type KeyVal struct {
	key interface{}
	val interface{}
}

// AccessField 保存access日志的参数
// 字段
//   kvs: 键值对
//   fieldLock: 对kvs加锁
type AccessField struct {
	kvs       []KeyVal
	fieldLock sync.Mutex
}

// AddField 增加键值对
func (af *AccessField) AddField(key, val interface{}) {
	af.fieldLock.Lock()
	af.kvs = append(af.kvs, KeyVal{key: key, val: val})
	af.fieldLock.Unlock()
}

// KvsKey Context key
type KvsKey struct {
}

// WithFieldContext 构造一个包含AccessField的Context
func WithFieldContext(ctx context.Context) context.Context {
	field := getField(ctx)
	// 已经包含了AccessField
	if field != nil {
		return ctx
	}

	field = &AccessField{}
	return context.WithValue(ctx, KvsKey{}, field)
}

// AddField 增加一个k-v
func AddField(ctx context.Context, key interface{}, val interface{}) {
	field := getField(ctx)
	if field == nil {
		return
	}
	field.AddField(key, val)
}

func getField(ctx context.Context) *AccessField {
	filed, ok := ctx.Value(KvsKey{}).(*AccessField)
	if !ok {
		return nil
	}
	return filed
}
