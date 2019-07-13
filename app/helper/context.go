package helper

/**
在请求上下文上设置和获取指定的key/val
*/

import (
	"context"
	"net/http"
)

//从请求上下文获取指定的key值
func ContextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

//将指定的key/val设置到上下文中
func ContextSet(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}

func GetStringByCtx(ctx context.Context, key string) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}

	str, ok := val.(string)
	if !ok {
		return ""
	}

	return str
}
