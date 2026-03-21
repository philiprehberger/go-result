# go-result

[![CI](https://github.com/philiprehberger/go-result/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-result/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-result.svg)](https://pkg.go.dev/github.com/philiprehberger/go-result) [![License](https://img.shields.io/github/license/philiprehberger/go-result)](LICENSE)

Generic Result type for Go — `Ok[T]` / `Err[T]` with mapping and chaining

## Installation

```bash
go get github.com/philiprehberger/go-result
```

## Usage

### Basic Result

```go
import "github.com/philiprehberger/go-result"

r := result.Ok(42)
fmt.Println(r.Unwrap()) // 42

r2 := result.Err[int](errors.New("not found"))
fmt.Println(r2.UnwrapOr(0)) // 0
```

### UnwrapOrElse — Compute Default from Error

```go
val := r.UnwrapOrElse(func(err error) int {
    log.Printf("falling back due to: %v", err)
    return -1
})
```

### Expect — Unwrap with Custom Panic Message

```go
cfg := r.Expect("config must be valid")
// panics with "config must be valid: <error>" if Err
```

### Or — Fallback Result

```go
r := result.Err[int](errors.New("fail"))
val := r.Or(result.Ok(42)) // Ok(42)
```

### String Representation

```go
fmt.Println(result.Ok(42))              // Ok(42)
fmt.Println(result.Err[int](err))       // Err(something failed)
```

> **Note:** `Unwrap()` and `Expect()` panic if the Result is an error. Use `UnwrapOr` or `UnwrapOrElse` for safe extraction.

### Try — Wrap (T, error)

```go
r := result.Try(func() (int, error) {
    return strconv.Atoi("42")
})
// Ok(42)
```

### Map and FlatMap

```go
r := result.Ok(10)
doubled := result.Map(r, func(v int) int { return v * 2 })
// Ok(20)

chained := result.FlatMap(r, func(v int) result.Result[string] {
    return result.Ok(fmt.Sprintf("value: %d", v))
})
```

### Match

```go
msg := result.Match(r,
    func(v int) string { return fmt.Sprintf("got %d", v) },
    func(err error) string { return fmt.Sprintf("error: %v", err) },
)
```

### Error Recovery

```go
r := result.Err[int](errors.New("primary failed"))
val := r.OrElse(func(err error) result.Result[int] {
    log.Printf("recovering from: %v", err)
    return result.Ok(fallbackValue())
})
```

### Filtering

```go
r := result.Ok(age)
valid := r.Filter(
    func(v int) bool { return v >= 18 },
    func(v int) error { return fmt.Errorf("age %d is below minimum 18", v) },
)
```

### Predicate Checks

```go
r := result.Ok(42)
r.IsOkAnd(func(v int) bool { return v > 0 })   // true
r.IsErrAnd(func(err error) bool { return true }) // false
```

### Side Effects

```go
r := result.Ok(42)
r.Tap(func(v int) { log.Printf("got value: %d", v) }).
    TapErr(func(err error) { log.Printf("got error: %v", err) })
```

### Collect Results

```go
results := []result.Result[int]{result.Ok(1), result.Ok(2), result.Ok(3)}
combined := result.All(results) // Ok([1, 2, 3])
```

## API

| Function / Method | Description |
|---|---|
| `Result[T]` | Generic type representing either a success or error value |
| `Ok[T](value T) Result[T]` | Create a successful Result |
| `Err[T](err error) Result[T]` | Create a failed Result |
| `Errf[T](format string, args ...any) Result[T]` | Create a failed Result with formatted error |
| `Try[T](fn func() (T, error)) Result[T]` | Wrap a (T, error) call into a Result |
| `Map[T, U](r Result[T], fn func(T) U) Result[U]` | Transform the success value |
| `FlatMap[T, U](r Result[T], fn func(T) Result[U]) Result[U]` | Chain Results with a function returning Result |
| `All[T](results []Result[T]) Result[[]T]` | Collect a slice of Results into a single Result |
| `Match[T, U](r Result[T], onOk func(T) U, onErr func(error) U) U` | Pattern match on Ok or Err |
| `(Result[T]) IsOk() bool` | True if the Result is a success |
| `(Result[T]) IsErr() bool` | True if the Result is an error |
| `(Result[T]) IsOkAnd(fn func(T) bool) bool` | True if Ok and value matches predicate |
| `(Result[T]) IsErrAnd(fn func(error) bool) bool` | True if Err and error matches predicate |
| `(Result[T]) Unwrap() T` | Return success value or panic |
| `(Result[T]) Expect(msg string) T` | Return success value or panic with message |
| `(Result[T]) UnwrapOr(def T) T` | Return success value or the provided default |
| `(Result[T]) UnwrapOrElse(fn func(error) T) T` | Return success value or compute from error |
| `(Result[T]) Or(other Result[T]) Result[T]` | Return self if Ok, otherwise the fallback |
| `(Result[T]) OrElse(fn func(error) Result[T]) Result[T]` | Return self if Ok, otherwise compute fallback |
| `(Result[T]) Filter(pred func(T) bool, errFn func(T) error) Result[T]` | Convert to Err if value fails predicate |
| `(Result[T]) Tap(fn func(T)) Result[T]` | Run side effect on Ok value, return unchanged |
| `(Result[T]) TapErr(fn func(error)) Result[T]` | Run side effect on Err value, return unchanged |
| `(Result[T]) Error() error` | Return the error or nil |
| `(Result[T]) String() string` | Human-readable representation |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
