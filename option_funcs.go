package goptions

// Map transforms the contained value if o is Some by applying fn, otherwise returns None[U].
func Map[T, U any](o Option[T], fn func(T) U) Option[U] {
	if !o.ok {
		return None[U]()
	}
	return Some(fn(o.value))
}

// MapOr returns fn(o.value) if o is Some, otherwise defaultVal.
func MapOr[T, U any](o Option[T], defaultVal U, fn func(T) U) U {
	if !o.ok {
		return defaultVal
	}
	return fn(o.value)
}

// MapOrElse returns fn(o.value) if o is Some, otherwise defaultFn().
func MapOrElse[T, U any](o Option[T], defaultFn func() U, fn func(T) U) U {
	if !o.ok {
		return defaultFn()
	}
	return fn(o.value)
}

// AndThen returns fn(o.value) if o is Some, otherwise returns None[U].
// Also known as FlatMap.
func AndThen[T, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	if !o.ok {
		return None[U]()
	}
	return fn(o.value)
}

// Filter returns o if it is Some and pred(o.value) is true, otherwise returns None[T].
func Filter[T any](o Option[T], pred func(T) bool) Option[T] {
	return o.Filter(pred)
}

// Or returns o if it is Some, otherwise returns other.
func Or[T any](o, other Option[T]) Option[T] {
	return o.Or(other)
}

// OrElse returns o if it is Some, otherwise returns fn().
func OrElse[T any](o Option[T], fn func() Option[T]) Option[T] {
	return o.OrElse(fn)
}

// Flatten converts Option[Option[T]] to Option[T]: Some(Some(x)) -> Some(x), otherwise None.
func Flatten[T any](o Option[Option[T]]) Option[T] {
	if !o.ok {
		return None[T]()
	}
	return o.value
}

// Tap calls fn with the contained value if o is Some, then returns o unchanged.
func Tap[T any](o Option[T], fn func(T)) Option[T] {
	return o.Tap(fn)
}

// Match returns onSome(o.value) if o is Some, otherwise returns onNone().
// Exhaustive handling of both branches; returns a single result of type R.
func Match[T, R any](o Option[T], onSome func(T) R, onNone func() R) R {
	if o.ok {
		return onSome(o.value)
	}
	return onNone()
}

// TryMap applies fn to the contained value if o is Some.
// If o is None, returns (None[U], nil). If fn returns an error, returns (None[U], err).
func TryMap[T, U any](o Option[T], fn func(T) (U, error)) (Option[U], error) {
	if !o.ok {
		return None[U](), nil
	}
	u, err := fn(o.value)
	if err != nil {
		return None[U](), err
	}
	return Some(u), nil
}

// Equals returns true if a and b are both None, or both Some with equal values.
// T must be comparable.
func Equals[T comparable](a, b Option[T]) bool {
	if a.ok != b.ok {
		return false
	}
	if !a.ok {
		return true
	}
	return a.value == b.value
}

// Pair holds two values; used by Zip.
type Pair[A, B any] struct {
	First  A
	Second B
}

// Zip returns Some(Pair{a.value, b.value}) if both a and b are Some, otherwise None.
func Zip[T, U any](a Option[T], b Option[U]) Option[Pair[T, U]] {
	if !a.ok || !b.ok {
		return None[Pair[T, U]]()
	}
	return Some(Pair[T, U]{First: a.value, Second: b.value})
}
