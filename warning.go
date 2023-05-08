package wch_otd_api

import "fmt"

/**
 *   Warning
 * This type is for handling the situation where an error might occur which is
 * not fatal within the same function call as an error which is fatal.
 * Specifically, I want to inform the caller when there was an issue with the
 * connection to the cache, but still return the data from the database
 */
type Warning interface {
	error
	Warning() (string, bool)
	HasError() bool
}
type CacheMiss struct{}

func (e CacheMiss) Error() string {
	return "cache miss"
}

type ErrorWithCacheFailure struct {
	err     error
	warning error
}

func (e ErrorWithCacheFailure) Error() string {
	if e.HasError() {
		return e.err.Error()
	} else {
		return fmt.Sprintf("warning: %e (nil error for this warning)", e.warning)
	}
}

func (e ErrorWithCacheFailure) HasError() bool {
	return e.err != nil
}

func (e ErrorWithCacheFailure) Warning() (string, bool) {
	if e.warning == nil {
		return "", false
	} else {
		return e.warning.Error(), true
	}
}

func CreateCacheMiss() *ErrorWithCacheFailure {
	return &ErrorWithCacheFailure{
		warning: CacheMiss{},
	}
}

func (e *ErrorWithCacheFailure) AndAnErrorToo(err error) *ErrorWithCacheFailure {
	e.err = err
	return e
}

func ErrorConnectingToCache(err error) *ErrorWithCacheFailure {
	return &ErrorWithCacheFailure{warning: err}
}
