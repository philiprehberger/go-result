# Changelog

## 0.3.1

- Add badges and Development section to README

## 0.3.0

- Add `String()` method for human-readable representation (`Ok(42)` / `Err(fail)`)
- Add `Expect(msg)` method that panics with a custom message on Err
- Add `Or(fallback)` method that returns the fallback Result when Err

## 0.2.0

- Complete test coverage for all public functions
- Add tests for `Errf`, `UnwrapOrElse`, `Match`, `Unwrap` panic path
- Add missing branch coverage for `UnwrapOr`, `Error`, `Map`, `FlatMap`

## 0.1.0

- Initial release
