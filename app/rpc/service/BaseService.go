package service

import "context"

type BaseService struct {
	Ctx context.Context
}

func (s *BaseService) SetCtx(ctx context.Context) {
	s.Ctx = ctx
}
