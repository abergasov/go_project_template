package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"
)

const DefaultFallbackTries = 5

var ErrorRepeatableFuncResNil = errors.New("RepeatableFunc result is nil")

type RepeatableFunc[T any] func() (T, error)

type FuncRepeater[T any] struct {
	ctx           context.Context
	fn            RepeatableFunc[T]
	fallback      *RepeatableFunc[T]
	maxTries      int
	fallbackTries int
	triesTimeout  *time.Duration
	linearTimeout bool
	errMsg        string
	exitErrors    []error
}

// Run repeats RepeatableFunc until any of these happens:
// - maxTries reached if it is set and fallback is not
// - fallback is set and maxTries + fallbackTries reached
// - ctx.Done()
// - exitErrors encountered
// - non-nil result returned
// it also records metrics to hist
func (r *FuncRepeater[T]) Run() (res T, err error) {
	triesCount := 0

	function := r.fn
	for {
		// if maxTries are set, and we don't have fallback enabled, we return after maxTries attempts
		// if fallback enabled we switch to fallback function and keep trying fallbackTries more
		if r.maxTries > 0 && triesCount >= r.maxTries {
			if r.fallback == nil || triesCount >= r.maxTries+r.fallbackTries {
				return res, fmt.Errorf("maximum number of retries reached: %w", err)
			}

			function = *r.fallback
		}

		if r.triesTimeout != nil {
			timeout := *r.triesTimeout
			if r.linearTimeout {
				timeout *= time.Duration(triesCount)
			}
			time.Sleep(timeout)
		}

		triesCount++

		select {
		case <-r.ctx.Done():
			if err != nil {
				return res, fmt.Errorf("%s: %w", r.errMsg, err)
			}

			return res, errors.New(r.errMsg)
		default:
			// run function
			res, err = function()
			if err != nil {
				for _, e := range r.exitErrors {
					if errors.Is(err, e) {
						return res, err
					}
				}

				continue
			}

			// some pointer magics to check if res is nil
			value := reflect.ValueOf(res)
			nullableKind := value.Kind() == reflect.Chan || value.Kind() == reflect.Func || value.Kind() == reflect.Map ||
				value.Kind() == reflect.Pointer || value.Kind() == reflect.UnsafePointer ||
				value.Kind() == reflect.Interface || value.Kind() == reflect.Slice
			if !value.IsValid() || (nullableKind && value.IsNil()) {
				err = ErrorRepeatableFuncResNil
				continue
			}
			return //nolint:nakedret // it's ok to return here
		}
	}
}

func NewFuncRepeater[T any](fn RepeatableFunc[T]) *FuncRepeater[T] {
	return &FuncRepeater[T]{
		fn:            fn,
		ctx:           context.Background(),
		fallbackTries: DefaultFallbackTries,
		linearTimeout: false,
	}
}

func (r *FuncRepeater[T]) WithCtx(ctx context.Context) *FuncRepeater[T] {
	r.ctx = ctx
	return r
}

func (r *FuncRepeater[T]) WithFallback(fallback RepeatableFunc[T]) *FuncRepeater[T] {
	r.fallback = &fallback
	return r
}

func (r *FuncRepeater[T]) WithMaxTries(num int) *FuncRepeater[T] {
	r.maxTries = num
	return r
}

func (r *FuncRepeater[T]) WithFallbackTries(num int) *FuncRepeater[T] {
	r.fallbackTries = num
	return r
}

func (r *FuncRepeater[T]) WithTriesTimeout(timeout time.Duration) *FuncRepeater[T] {
	r.triesTimeout = &timeout
	return r
}

func (r *FuncRepeater[T]) WithLinearTimeout(linear bool) *FuncRepeater[T] {
	r.linearTimeout = linear
	return r
}

func (r *FuncRepeater[T]) WithErrMsg(errMsg string) *FuncRepeater[T] {
	r.errMsg = errMsg
	return r
}

func (r *FuncRepeater[T]) WithExitErrors(errs ...error) *FuncRepeater[T] {
	r.exitErrors = errs
	return r
}
