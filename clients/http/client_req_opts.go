package http

import (
	"context"
)

// ReqOptions presents request options.
type ReqOptions struct {
	Header map[string]string
}

var reqOptsCtxKey = new(int)

func WithReqOpts(ctx context.Context, reqOpts ReqOptions) context.Context {
	if defaults, ok := ctx.Value(reqOptsCtxKey).(ReqOptions); ok {
		reqOpts = (&reqOpts).WithDefaults(defaults)
	}
	return context.WithValue(ctx, reqOptsCtxKey, reqOpts)
}

func GetReqOpts(ctx context.Context) ReqOptions {
	if reqOpts, ok := ctx.Value(reqOptsCtxKey).(ReqOptions); ok {
		return reqOpts
	}
	return ReqOptions{}
}

func WithClearReqOpts(ctx context.Context) context.Context {
	return context.WithValue(ctx, reqOptsCtxKey, ReqOptions{})
}

func (opts *ReqOptions) WithDefaults(defaults ReqOptions) ReqOptions {
	if opts == nil {
		return defaults
	}
	if opts.Header != nil && len(opts.Header) > 0 {
		defaults.Header = opts.Header
	}
	return defaults
}
