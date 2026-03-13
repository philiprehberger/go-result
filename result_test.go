package result

import (
	"errors"
	"fmt"
	"testing"
)

func TestOk(t *testing.T) {
	r := Ok(42)
	if !r.IsOk() {
		t.Fatal("expected Ok")
	}
	if r.IsErr() {
		t.Fatal("expected not Err")
	}
	if r.Unwrap() != 42 {
		t.Fatalf("expected 42, got %d", r.Unwrap())
	}
}

func TestErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	if !r.IsErr() {
		t.Fatal("expected Err")
	}
	if r.IsOk() {
		t.Fatal("expected not Ok")
	}
	if r.Error() == nil {
		t.Fatal("expected error")
	}
	if r.Error().Error() != "fail" {
		t.Fatalf("expected 'fail', got %q", r.Error().Error())
	}
}

func TestErrf(t *testing.T) {
	r := Errf[string]("code %d: %s", 404, "not found")
	if !r.IsErr() {
		t.Fatal("expected Err")
	}
	expected := "code 404: not found"
	if r.Error().Error() != expected {
		t.Fatalf("expected %q, got %q", expected, r.Error().Error())
	}
}

func TestUnwrapOk(t *testing.T) {
	r := Ok("hello")
	if r.Unwrap() != "hello" {
		t.Fatalf("expected 'hello', got %q", r.Unwrap())
	}
}

func TestUnwrapPanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic from Unwrap on Err")
		}
		msg := fmt.Sprintf("%v", r)
		if msg == "" {
			t.Fatal("expected non-empty panic message")
		}
	}()
	r := Err[int](errors.New("fail"))
	r.Unwrap() // should panic
}

func TestUnwrapOrOk(t *testing.T) {
	r := Ok(42)
	if r.UnwrapOr(99) != 42 {
		t.Fatal("expected Ok value, not default")
	}
}

func TestUnwrapOrErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	if r.UnwrapOr(99) != 99 {
		t.Fatal("expected default")
	}
}

func TestUnwrapOrElseOk(t *testing.T) {
	r := Ok(42)
	val := r.UnwrapOrElse(func(err error) int {
		t.Fatal("function should not be called for Ok")
		return -1
	})
	if val != 42 {
		t.Fatalf("expected 42, got %d", val)
	}
}

func TestUnwrapOrElseErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	var receivedErr error
	val := r.UnwrapOrElse(func(err error) int {
		receivedErr = err
		return -1
	})
	if val != -1 {
		t.Fatalf("expected -1, got %d", val)
	}
	if receivedErr == nil || receivedErr.Error() != "fail" {
		t.Fatalf("expected 'fail' error in callback, got %v", receivedErr)
	}
}

func TestErrorOk(t *testing.T) {
	r := Ok(42)
	if r.Error() != nil {
		t.Fatal("expected nil error for Ok")
	}
}

func TestErrorErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	if r.Error() == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestMapOk(t *testing.T) {
	r := Ok(10)
	mapped := Map(r, func(v int) string { return fmt.Sprintf("val=%d", v) })
	if !mapped.IsOk() {
		t.Fatal("expected Ok")
	}
	if mapped.Unwrap() != "val=10" {
		t.Fatalf("expected 'val=10', got %q", mapped.Unwrap())
	}
}

func TestMapErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	mapped := Map(r, func(v int) string { return "should not run" })
	if !mapped.IsErr() {
		t.Fatal("expected Err to propagate through Map")
	}
	if mapped.Error().Error() != "fail" {
		t.Fatalf("expected original error, got %v", mapped.Error())
	}
}

func TestFlatMapOk(t *testing.T) {
	r := Ok(10)
	chained := FlatMap(r, func(v int) Result[int] {
		return Ok(v * 2)
	})
	if chained.Unwrap() != 20 {
		t.Fatal("expected 20")
	}
}

func TestFlatMapErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	chained := FlatMap(r, func(v int) Result[int] {
		t.Fatal("function should not be called for Err")
		return Ok(0)
	})
	if !chained.IsErr() {
		t.Fatal("expected Err to propagate through FlatMap")
	}
}

func TestTryOk(t *testing.T) {
	r := Try(func() (int, error) { return 42, nil })
	if !r.IsOk() || r.Unwrap() != 42 {
		t.Fatal("expected Ok(42)")
	}
}

func TestTryErr(t *testing.T) {
	r := Try(func() (int, error) { return 0, errors.New("fail") })
	if !r.IsErr() {
		t.Fatal("expected Err")
	}
}

func TestAllOk(t *testing.T) {
	results := []Result[int]{Ok(1), Ok(2), Ok(3)}
	r := All(results)
	if !r.IsOk() {
		t.Fatal("expected Ok")
	}
	vals := r.Unwrap()
	if len(vals) != 3 || vals[0] != 1 || vals[1] != 2 || vals[2] != 3 {
		t.Fatalf("expected [1,2,3], got %v", vals)
	}
}

func TestAllWithErr(t *testing.T) {
	results := []Result[int]{Ok(1), Err[int](errors.New("fail")), Ok(3)}
	r := All(results)
	if !r.IsErr() {
		t.Fatal("expected Err")
	}
	if r.Error().Error() != "fail" {
		t.Fatalf("expected first error 'fail', got %v", r.Error())
	}
}

func TestAllEmpty(t *testing.T) {
	results := []Result[int]{}
	r := All(results)
	if !r.IsOk() {
		t.Fatal("expected Ok for empty slice")
	}
	if len(r.Unwrap()) != 0 {
		t.Fatal("expected empty slice")
	}
}

func TestMatchOk(t *testing.T) {
	r := Ok(42)
	msg := Match(r,
		func(v int) string { return fmt.Sprintf("got %d", v) },
		func(err error) string { return "error" },
	)
	if msg != "got 42" {
		t.Fatalf("expected 'got 42', got %q", msg)
	}
}

func TestMatchErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	msg := Match(r,
		func(v int) string { return "ok" },
		func(err error) string { return fmt.Sprintf("err: %v", err) },
	)
	if msg != "err: fail" {
		t.Fatalf("expected 'err: fail', got %q", msg)
	}
}

func TestMatchReturnType(t *testing.T) {
	r := Ok(10)
	val := Match(r,
		func(v int) int { return v * 2 },
		func(err error) int { return -1 },
	)
	if val != 20 {
		t.Fatalf("expected 20, got %d", val)
	}
}
