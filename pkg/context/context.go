package context

import (
	"context"
	"os"
)

type Key interface {
	~struct{}
}

type ContextValue[keyType Key] interface {
	Key() keyType
}

func With[_type ContextValue[keyType], keyType Key](ctx context.Context, _ctx _type) context.Context {
	return context.WithValue(ctx, _ctx.Key(), _ctx)
}

func From[_type ContextValue[keyType], keyType Key](ctx context.Context) (_type, bool) {
	var value _type
	_ctx, ok := ctx.Value(value.Key()).(_type)
	return _ctx, ok
}

type requestContextKey struct{}

type requestContext struct {
	requestID string
}

func (ctx *requestContext) Key() requestContextKey {
	return requestContextKey{}
}

func WithRequestCtx(ctx context.Context, requestID string) context.Context {
	return With(ctx, &requestContext{
		requestID: requestID,
	})
}

func RequestFromCtx(ctx context.Context) (*requestContext, bool) {
	return From[*requestContext](ctx)
}

type stopContextKey struct{}

type stopContext struct {
	Signal os.Signal
}

func (ctx *stopContext) Key() stopContextKey {
	return stopContextKey{}
}

func WithStopCtx(ctx context.Context, signal os.Signal) context.Context {
	return With(ctx, &stopContext{
		Signal: signal,
	})
}

func StopFromCtx(ctx context.Context) (*stopContext, bool) {
	return From[*stopContext](ctx)
}
