package goptions

// Some returns an Option containing the value v.
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
}

// None returns an Option with no value.
func None[T any]() Option[T] {
	return Option[T]{ok: false}
}

// FromPtr returns Some(*p) if p is non-nil, otherwise None.
func FromPtr[T any](p *T) Option[T] {
	if p == nil {
		return Option[T]{ok: false}
	}
	return Option[T]{value: *p, ok: true}
}

// FromTuple builds an Option from a value and a boolean (e.g. comma-ok form).
// If ok is true, returns Some(v); otherwise None.
func FromTuple[T any](v T, ok bool) Option[T] {
	if !ok {
		return Option[T]{ok: false}
	}
	return Option[T]{value: v, ok: true}
}

// Try returns Some(v) if err is nil, otherwise None.
// Useful for converting (T, error) returns into Option[T].
func Try[T any](v T, err error) Option[T] {
	if err != nil {
		return Option[T]{ok: false}
	}
	return Option[T]{value: v, ok: true}
}

// Cond returns Some(v) if ok is true, otherwise None[T]().
func Cond[T any](ok bool, v T) Option[T] {
	if !ok {
		return Option[T]{ok: false}
	}
	return Option[T]{value: v, ok: true}
}
