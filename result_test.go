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

func TestStringOk(t *testing.T) {
	r := Ok(42)
	s := r.String()
	if s != "Ok(42)" {
		t.Fatalf("expected 'Ok(42)', got %q", s)
	}
}

func TestStringErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	s := r.String()
	if s != "Err(fail)" {
		t.Fatalf("expected 'Err(fail)', got %q", s)
	}
}

func TestExpectOk(t *testing.T) {
	r := Ok(42)
	val := r.Expect("should not panic")
	if val != 42 {
		t.Fatalf("expected 42, got %d", val)
	}
}

func TestExpectPanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic from Expect on Err")
		}
		msg := fmt.Sprintf("%v", r)
		if msg != "config error: fail" {
			t.Fatalf("expected 'config error: fail', got %q", msg)
		}
	}()
	r := Err[int](errors.New("fail"))
	r.Expect("config error")
}

func TestOrOk(t *testing.T) {
	r := Ok(42)
	fallback := Ok(99)
	got := r.Or(fallback)
	if got.Unwrap() != 42 {
		t.Fatalf("expected 42 from Or on Ok, got %d", got.Unwrap())
	}
}

func TestOrErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	fallback := Ok(99)
	got := r.Or(fallback)
	if got.Unwrap() != 99 {
		t.Fatalf("expected 99 from Or on Err, got %d", got.Unwrap())
	}
}

func TestOrBothErr(t *testing.T) {
	r := Err[int](errors.New("first"))
	fallback := Err[int](errors.New("second"))
	got := r.Or(fallback)
	if !got.IsErr() {
		t.Fatal("expected Err when both are Err")
	}
	if got.Error().Error() != "second" {
		t.Fatalf("expected 'second' error, got %q", got.Error().Error())
	}
}

func TestOrElseOk(t *testing.T) {
	r := Ok(42)
	got := r.OrElse(func(err error) Result[int] {
		t.Fatal("function should not be called for Ok")
		return Ok(0)
	})
	if got.Unwrap() != 42 {
		t.Fatalf("expected 42, got %d", got.Unwrap())
	}
}

func TestOrElseErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	got := r.OrElse(func(err error) Result[int] {
		if err.Error() != "fail" {
			t.Fatalf("expected 'fail' error in callback, got %v", err)
		}
		return Ok(99)
	})
	if got.Unwrap() != 99 {
		t.Fatalf("expected 99, got %d", got.Unwrap())
	}
}

func TestOrElseErrToErr(t *testing.T) {
	r := Err[int](errors.New("first"))
	got := r.OrElse(func(err error) Result[int] {
		return Err[int](errors.New("second"))
	})
	if !got.IsErr() {
		t.Fatal("expected Err")
	}
	if got.Error().Error() != "second" {
		t.Fatalf("expected 'second', got %q", got.Error().Error())
	}
}

func TestFilterOkPass(t *testing.T) {
	r := Ok(42)
	got := r.Filter(
		func(v int) bool { return v > 0 },
		func(v int) error { return fmt.Errorf("expected positive, got %d", v) },
	)
	if !got.IsOk() || got.Unwrap() != 42 {
		t.Fatalf("expected Ok(42), got %v", got)
	}
}

func TestFilterOkFail(t *testing.T) {
	r := Ok(-5)
	got := r.Filter(
		func(v int) bool { return v > 0 },
		func(v int) error { return fmt.Errorf("expected positive, got %d", v) },
	)
	if !got.IsErr() {
		t.Fatal("expected Err when predicate fails")
	}
	if got.Error().Error() != "expected positive, got -5" {
		t.Fatalf("unexpected error: %v", got.Error())
	}
}

func TestFilterErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	got := r.Filter(
		func(v int) bool { return true },
		func(v int) error { return errors.New("should not run") },
	)
	if !got.IsErr() || got.Error().Error() != "fail" {
		t.Fatalf("expected original Err, got %v", got)
	}
}

func TestIsOkAndTrue(t *testing.T) {
	r := Ok(42)
	if !r.IsOkAnd(func(v int) bool { return v == 42 }) {
		t.Fatal("expected true")
	}
}

func TestIsOkAndFalse(t *testing.T) {
	r := Ok(42)
	if r.IsOkAnd(func(v int) bool { return v == 0 }) {
		t.Fatal("expected false when predicate fails")
	}
}

func TestIsOkAndErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	if r.IsOkAnd(func(v int) bool { return true }) {
		t.Fatal("expected false for Err")
	}
}

func TestIsErrAndTrue(t *testing.T) {
	r := Err[int](errors.New("not found"))
	if !r.IsErrAnd(func(err error) bool { return err.Error() == "not found" }) {
		t.Fatal("expected true")
	}
}

func TestIsErrAndFalse(t *testing.T) {
	r := Err[int](errors.New("not found"))
	if r.IsErrAnd(func(err error) bool { return err.Error() == "timeout" }) {
		t.Fatal("expected false when predicate fails")
	}
}

func TestIsErrAndOk(t *testing.T) {
	r := Ok(42)
	if r.IsErrAnd(func(err error) bool { return true }) {
		t.Fatal("expected false for Ok")
	}
}

func TestTapOk(t *testing.T) {
	r := Ok(42)
	var captured int
	got := r.Tap(func(v int) { captured = v })
	if captured != 42 {
		t.Fatalf("expected captured 42, got %d", captured)
	}
	if !got.IsOk() || got.Unwrap() != 42 {
		t.Fatalf("expected original Ok(42) returned, got %v", got)
	}
}

func TestTapErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	called := false
	got := r.Tap(func(v int) { called = true })
	if called {
		t.Fatal("Tap should not call fn for Err")
	}
	if !got.IsErr() {
		t.Fatal("expected Err returned")
	}
}

func TestTapErrOnErr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	var captured error
	got := r.TapErr(func(err error) { captured = err })
	if captured == nil || captured.Error() != "fail" {
		t.Fatalf("expected captured 'fail' error, got %v", captured)
	}
	if !got.IsErr() {
		t.Fatal("expected Err returned")
	}
}

func TestTapErrOnOk(t *testing.T) {
	r := Ok(42)
	called := false
	got := r.TapErr(func(err error) { called = true })
	if called {
		t.Fatal("TapErr should not call fn for Ok")
	}
	if !got.IsOk() || got.Unwrap() != 42 {
		t.Fatal("expected original Ok(42) returned")
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
