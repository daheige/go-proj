package task

import "context"

type BaseTask struct {
	ctx context.Context
}

func (b *BaseTask) SetCtx(ctx context.Context) {
	b.ctx = ctx
}
