package goptions

// Some returns an Option containing the value v.
//
// Example:
//
//	o := Some(42)
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
}

// None returns an Option with no value.
//
// Example:
//
//	o := None[int]()
func None[T any]() Option[T] {
	return Option[T]{ok: false}
}

// FromPtr returns Some(*p) if p is non-nil, otherwise None.
//
// Example:
//
//	x := 7
//	o := FromPtr(&x)   // Some(7)
//	o := FromPtr[int](nil)  // None[int]()
func FromPtr[T any](p *T) Option[T] {
	if p == nil {
		return Option[T]{ok: false}
	}
	return Option[T]{value: *p, ok: true}
}

// FromTuple builds an Option from a value and a boolean (e.g. comma-ok form).
// If ok is true, returns Some(v); otherwise None.
//
// Example:
//
//	v, ok := m["key"]
//	o := FromTuple(v, ok)
func FromTuple[T any](v T, ok bool) Option[T] {
	if !ok {
		return Option[T]{ok: false}
	}
	return Option[T]{value: v, ok: true}
}

// Try returns Some(v) if err is nil, otherwise None.
// Useful for converting (T, error) returns into Option[T].
//
// Example:
//
//	n, err := strconv.Atoi(s)
//	o := Try(n, err)
func Try[T any](v T, err error) Option[T] {
	if err != nil {
		return Option[T]{ok: false}
	}
	return Option[T]{value: v, ok: true}
}

// Cond returns Some(v) if ok is true, otherwise None[T]().
//
// Example:
//
//	o := Cond(x != nil, *x)
func Cond[T any](ok bool, v T) Option[T] {
	if !ok {
		return Option[T]{ok: false}
	}
	return Option[T]{value: v, ok: true}
}
