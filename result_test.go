package result

import (
	"errors"
	"testing"
)

func TestOk(t *testing.T) {
	r := Ok(42)
	if !r.IsOk() {
		t.Fatal("expected Ok")
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
	if r.Error() == nil {
		t.Fatal("expected error")
	}
}

func TestUnwrapOr(t *testing.T) {
	r := Err[int](errors.New("fail"))
	if r.UnwrapOr(99) != 99 {
		t.Fatal("expected default")
	}
}

func TestMap(t *testing.T) {
	r := Ok(10)
	mapped := Map(r, func(v int) string { return "ok" })
	if !mapped.IsOk() || mapped.Unwrap() != "ok" {
		t.Fatal("expected mapped Ok")
	}
}

func TestFlatMap(t *testing.T) {
	r := Ok(10)
	chained := FlatMap(r, func(v int) Result[int] {
		return Ok(v * 2)
	})
	if chained.Unwrap() != 20 {
		t.Fatal("expected 20")
	}
}

func TestTry(t *testing.T) {
	r := Try(func() (int, error) { return 42, nil })
	if !r.IsOk() || r.Unwrap() != 42 {
		t.Fatal("expected Ok(42)")
	}

	r2 := Try(func() (int, error) { return 0, errors.New("fail") })
	if !r2.IsErr() {
		t.Fatal("expected Err")
	}
}

func TestAll(t *testing.T) {
	results := []Result[int]{Ok(1), Ok(2), Ok(3)}
	r := All(results)
	if !r.IsOk() {
		t.Fatal("expected Ok")
	}

	results2 := []Result[int]{Ok(1), Err[int](errors.New("fail")), Ok(3)}
	r2 := All(results2)
	if !r2.IsErr() {
		t.Fatal("expected Err")
	}
}
