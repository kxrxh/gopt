// Package gopt provides an Option type for Go 1.21+.
// Option represents an optional value: either Some(T) or None.
package gopt

// Option is a generic container for an optional value of type T.
// It is either Some (value present) or None (no value).
// Create options using Some, None, FromPtr, FromTuple, or Try.
//
// Example:
//
//	o := Some(42)
//	if o.IsSome() { v := o.UnwrapOr(0) }
type Option[T any] struct {
	value T
	ok    bool
}

// IsSome returns true if the option contains a value.
//
// Example:
//
//	Some(42).IsSome()  // true
//	None[int]().IsSome()  // false
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone returns true if the option does not contain a value.
//
// Example:
//
//	None[int]().IsNone()  // true
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// Get returns the contained value and a boolean indicating whether a value was present.
// If the option is None, the value is the zero value of T and ok is false.
//
// Example:
//
//	v, ok := Some(42).Get()   // v=42, ok=true
//	v, ok := None[int]().Get()  // v=0, ok=false
func (o Option[T]) Get() (T, bool) {
	return o.value, o.ok
}

// Unwrap returns the contained value. It panics if the option is None.
// Prefer UnwrapOr, UnwrapOrElse, or Match when a default or explicit handling is needed.
//
// Example:
//
//	v := Some(42).Unwrap()  // v=42
//	v := None[int]().Unwrap()  // panics
func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("gopt: Unwrap called on None")
	}
	return o.value
}

// UnwrapOr returns the contained value if Some, otherwise returns defaultVal.
//
// Example:
//
//	Some(42).UnwrapOr(0)   // 42
//	None[int]().UnwrapOr(0)  // 0
func (o Option[T]) UnwrapOr(defaultVal T) T {
	if o.ok {
		return o.value
	}
	return defaultVal
}

// UnwrapOrElse returns the contained value if Some, otherwise returns the result of calling fn.
//
// Example:
//
//	None[int]().UnwrapOrElse(func() int { return 99 })  // 99
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.ok {
		return o.value
	}
	return fn()
}

// Expect returns the contained value if Some. It panics with the given message if None.
//
// Example:
//
//	v := Some(42).Expect("required")  // v=42
func (o Option[T]) Expect(msg string) T {
	if !o.ok {
		panic("gopt: " + msg)
	}
	return o.value
}

// Filter returns this option if Some and pred(value) is true, otherwise returns None.
//
// Example:
//
//	Some(4).Filter(func(x int) bool { return x%2 == 0 })  // Some(4)
//	Some(3).Filter(func(x int) bool { return x%2 == 0 })   // None[int]()
func (o Option[T]) Filter(pred func(T) bool) Option[T] {
	if !o.ok || !pred(o.value) {
		return None[T]()
	}
	return o
}

// Or returns this option if Some, otherwise returns other.
//
// Example:
//
//	None[int]().Or(Some(99))  // Some(99)
func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.ok {
		return o
	}
	return other
}

// OrElse returns this option if Some, otherwise returns fn().
//
// Example:
//
//	None[int]().OrElse(func() Option[int] { return Some(99) })  // Some(99)
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if o.ok {
		return o
	}
	return fn()
}

// Tap calls fn with the contained value if Some, then returns this option unchanged.
// Useful for side effects (e.g. logging) without changing the option.
//
// Example:
//
//	Some(42).Tap(func(x int) { log.Println(x) })  // logs 42, returns Some(42)
func (o Option[T]) Tap(fn func(T)) Option[T] {
	if o.ok {
		fn(o.value)
	}
	return o
}

// ToPointer returns nil if None, or a pointer to a copy of the value if Some.
// Caller may mutate the returned pointer; the value is a copy.
//
// Example:
//
//	p := Some(42).ToPointer()  // *int pointing to 42
//	p := None[int]().ToPointer()  // nil
func (o Option[T]) ToPointer() *T {
	if !o.ok {
		return nil
	}
	v := new(T)
	*v = o.value
	return v
}
