package job

import "context"

type BaseJob struct {
	ctx context.Context
}

func (j *BaseJob) SetCtx(ctx context.Context) {
	j.ctx = ctx
}
