// Package goptions provides a Rust-inspired Option type for Go 1.21+.
// Option represents an optional value: either Some(T) or None.
package goptions

// Option is a generic container for an optional value of type T.
// It is either Some (value present) or None (no value).
// Create options using Some, None, FromPtr, FromTuple, or Try.
type Option[T any] struct {
	value T
	ok    bool
}

// IsSome returns true if the option contains a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone returns true if the option does not contain a value.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// Get returns the contained value and a boolean indicating whether a value was present.
// If the option is None, the value is the zero value of T and ok is false.
func (o Option[T]) Get() (T, bool) {
	return o.value, o.ok
}

// Unwrap returns the contained value. It panics if the option is None.
// Prefer UnwrapOr, UnwrapOrElse, or Match when a default or explicit handling is needed.
func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("goptions: Unwrap called on None")
	}
	return o.value
}

// UnwrapOr returns the contained value if Some, otherwise returns defaultVal.
func (o Option[T]) UnwrapOr(defaultVal T) T {
	if o.ok {
		return o.value
	}
	return defaultVal
}

// UnwrapOrElse returns the contained value if Some, otherwise returns the result of calling fn.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.ok {
		return o.value
	}
	return fn()
}

// Expect returns the contained value if Some. It panics with the given message if None.
func (o Option[T]) Expect(msg string) T {
	if !o.ok {
		panic("goptions: " + msg)
	}
	return o.value
}

// Filter returns this option if Some and pred(value) is true, otherwise returns None.
func (o Option[T]) Filter(pred func(T) bool) Option[T] {
	if !o.ok || !pred(o.value) {
		return None[T]()
	}
	return o
}

// Or returns this option if Some, otherwise returns other.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.ok {
		return o
	}
	return other
}

// OrElse returns this option if Some, otherwise returns fn().
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if o.ok {
		return o
	}
	return fn()
}

// Tap calls fn with the contained value if Some, then returns this option unchanged.
// Useful for side effects (e.g. logging) without changing the option.
func (o Option[T]) Tap(fn func(T)) Option[T] {
	if o.ok {
		fn(o.value)
	}
	return o
}

// ToPointer returns nil if None, or a pointer to a copy of the value if Some.
// Caller may mutate the returned pointer; the value is a copy.
func (o Option[T]) ToPointer() *T {
	if !o.ok {
		return nil
	}
	v := new(T)
	*v = o.value
	return v
}
