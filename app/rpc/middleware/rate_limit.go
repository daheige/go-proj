package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Limiter defines the interface to perform request rate limiting.
// If Limit function return true, the request will be rejected.
// Otherwise, the request will pass.
type Limiter interface {
	Limit() bool
}

type RateLimitWare struct {
}

type MockPassLimiter struct{}

func (*MockPassLimiter) Limit() bool {
	return false
}

func Limit(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Limit() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleare, please retry later.", info.FullMethod)
		}

		log.Println("req method: ", info.FullMethod)

		return handler(ctx, req)
	}
}
