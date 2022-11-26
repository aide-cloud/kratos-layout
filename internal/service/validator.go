package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

type validator interface {
	Validate() error
}

type CallbackFunc[T any, R any] func(ctx context.Context, req *T) (*R, error)

func Validate[T any, R any](ctx context.Context, req any, callback CallbackFunc[T, R]) (*R, error) {
	if v, ok := req.(validator); ok {
		if err := v.Validate(); err != nil {
			return nil, errors.BadRequest("", err.Error())
		}
	}

	if v, ok := req.(*T); ok {
		return callback(ctx, v)
	}

	return nil, errors.BadRequest("VALIDATOR", "invalid request")
}
