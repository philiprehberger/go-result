# go-result

Generic Result type for Go — `Ok[T]` / `Err[T]` with mapping and chaining.

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

### Collect Results

```go
results := []result.Result[int]{result.Ok(1), result.Ok(2), result.Ok(3)}
combined := result.All(results) // Ok([1, 2, 3])
```

## License

MIT
