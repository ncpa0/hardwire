package utils

import (
	Promise "github.com/ncpa0cpl/go_promise"
)

func InParallel[T any, U any](arr []T, processFn func(value T) (U, error)) ([]U, []error) {
	promises := make([]*Promise.Promise[U], len(arr))
	for i, value := range arr {
		promises[i] = Promise.New(func() (U, error) {
			return processFn(value)
		})
	}
	promiseResults := *Promise.AwaitAll(promises)

	values := make([]U, len(arr))
	errors := make([]error, 0)

	for i, promiseResult := range promiseResults {
		if promiseResult.Err != nil {
			errors = append(errors, promiseResult.Err)
		} else {
			values[i] = promiseResult.Result
		}
	}

	return values, errors
}
