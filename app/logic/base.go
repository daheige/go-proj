package logic

import (
	"context"
)

// Ctx 标准上下文，如果是gin.Context就是ctx.Request.Context(),如果是grpc就是ctx
type BaseLogic struct {
	Ctx context.Context
}

func (b *BaseLogic) SetCtx(ctx context.Context) {
	b.Ctx = ctx
}
