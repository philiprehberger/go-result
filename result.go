// Package result provides a generic Result type for Go inspired by Rust's Result<T, E>.
package result

import "fmt"

// Result represents either a success value (Ok) or an error value (Err).
type Result[T any] struct {
	value T
	err   error
	ok    bool
}

// Ok creates a successful Result containing the given value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

// Err creates a failed Result containing the given error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// Errf creates a failed Result with a formatted error message.
func Errf[T any](format string, args ...any) Result[T] {
	return Result[T]{err: fmt.Errorf(format, args...), ok: false}
}

// IsOk returns true if the Result is a success.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr returns true if the Result is an error.
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// String returns a human-readable representation of the Result.
func (r Result[T]) String() string {
	if r.ok {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Err(%v)", r.err)
}

// Unwrap returns the success value or panics if the Result is an error.
func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic(fmt.Sprintf("called Unwrap on an Err value: %v", r.err))
	}
	return r.value
}

// Expect returns the success value or panics with the given message if the Result is an error.
func (r Result[T]) Expect(msg string) T {
	if !r.ok {
		panic(fmt.Sprintf("%s: %v", msg, r.err))
	}
	return r.value
}

// Or returns the Result if it is Ok, otherwise returns the provided fallback Result.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.ok {
		return r
	}
	return other
}

// UnwrapOr returns the success value or the provided default.
func (r Result[T]) UnwrapOr(def T) T {
	if r.ok {
		return r.value
	}
	return def
}

// UnwrapOrElse returns the success value or calls the provided function.
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.ok {
		return r.value
	}
	return fn(r.err)
}

// Error returns the error value or nil if the Result is Ok.
func (r Result[T]) Error() error {
	if r.ok {
		return nil
	}
	return r.err
}

// OrElse returns r if Ok, otherwise calls fn with the error to produce an alternative Result.
func (r Result[T]) OrElse(fn func(error) Result[T]) Result[T] {
	if r.ok {
		return r
	}
	return fn(r.err)
}

// Filter returns Err if the Ok value doesn't match the predicate. errFn produces the error from the value.
func (r Result[T]) Filter(predicate func(T) bool, errFn func(T) error) Result[T] {
	if !r.ok {
		return r
	}
	if predicate(r.value) {
		return r
	}
	return Err[T](errFn(r.value))
}

// IsOkAnd returns true if the Result is Ok and the value matches the predicate.
func (r Result[T]) IsOkAnd(predicate func(T) bool) bool {
	return r.ok && predicate(r.value)
}

// IsErrAnd returns true if the Result is Err and the error matches the predicate.
func (r Result[T]) IsErrAnd(predicate func(error) bool) bool {
	return !r.ok && predicate(r.err)
}

// Tap calls fn with the Ok value for side effects, then returns the original Result unchanged.
func (r Result[T]) Tap(fn func(T)) Result[T] {
	if r.ok {
		fn(r.value)
	}
	return r
}

// TapErr calls fn with the error for side effects, then returns the original Result unchanged.
func (r Result[T]) TapErr(fn func(error)) Result[T] {
	if !r.ok {
		fn(r.err)
	}
	return r
}

// Map transforms the success value using the given function.
func Map[T any, U any](r Result[T], fn func(T) U) Result[U] {
	if r.ok {
		return Ok[U](fn(r.value))
	}
	return Err[U](r.err)
}

// FlatMap transforms the success value using a function that returns a Result.
func FlatMap[T any, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.ok {
		return fn(r.value)
	}
	return Err[U](r.err)
}

// Try wraps a function call that returns (T, error) into a Result.
func Try[T any](fn func() (T, error)) Result[T] {
	value, err := fn()
	if err != nil {
		return Err[T](err)
	}
	return Ok[T](value)
}

// All collects a slice of Results into a Result containing a slice of values.
// Returns the first error encountered.
func All[T any](results []Result[T]) Result[[]T] {
	values := make([]T, 0, len(results))
	for _, r := range results {
		if r.IsErr() {
			return Err[[]T](r.err)
		}
		values = append(values, r.value)
	}
	return Ok[[]T](values)
}

// Match applies one of two functions depending on whether the Result is Ok or Err.
func Match[T any, U any](r Result[T], onOk func(T) U, onErr func(error) U) U {
	if r.ok {
		return onOk(r.value)
	}
	return onErr(r.err)
}
