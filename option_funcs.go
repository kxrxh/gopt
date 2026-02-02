package goptions

// Map transforms the contained value if o is Some by applying fn, otherwise returns None[U].
//
// Example:
//
//	o := Map(Some(21), func(x int) int { return x * 2 })  // Some(42)
func Map[T, U any](o Option[T], fn func(T) U) Option[U] {
	if !o.ok {
		return None[U]()
	}
	return Some(fn(o.value))
}

// MapOr returns fn(o.value) if o is Some, otherwise defaultVal.
//
// Example:
//
//	v := MapOr(Some(10), 0, func(x int) int { return x * 2 })  // 20
//	v := MapOr(None[int](), 99, func(x int) int { return x })  // 99
func MapOr[T, U any](o Option[T], defaultVal U, fn func(T) U) U {
	if !o.ok {
		return defaultVal
	}
	return fn(o.value)
}

// MapOrElse returns fn(o.value) if o is Some, otherwise defaultFn().
//
// Example:
//
//	v := MapOrElse(None[int](), func() int { return 42 }, func(x int) int { return x })
func MapOrElse[T, U any](o Option[T], defaultFn func() U, fn func(T) U) U {
	if !o.ok {
		return defaultFn()
	}
	return fn(o.value)
}

// AndThen returns fn(o.value) if o is Some, otherwise returns None[U].
// Also known as FlatMap.
//
// Example:
//
//	o := AndThen(Some(4), func(x int) Option[int] { return Some(x * x) })  // Some(16)
func AndThen[T, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	if !o.ok {
		return None[U]()
	}
	return fn(o.value)
}

// Filter returns o if it is Some and pred(o.value) is true, otherwise returns None[T].
//
// Example:
//
//	o := Filter(Some(4), func(x int) bool { return x%2 == 0 })  // Some(4)
func Filter[T any](o Option[T], pred func(T) bool) Option[T] {
	return o.Filter(pred)
}

// Or returns o if it is Some, otherwise returns other.
//
// Example:
//
//	o := Or(None[int](), Some(99))  // Some(99)
func Or[T any](o, other Option[T]) Option[T] {
	return o.Or(other)
}

// OrElse returns o if it is Some, otherwise returns fn().
//
// Example:
//
//	o := OrElse(None[int](), func() Option[int] { return Some(99) })
func OrElse[T any](o Option[T], fn func() Option[T]) Option[T] {
	return o.OrElse(fn)
}

// Flatten converts Option[Option[T]] to Option[T]: Some(Some(x)) -> Some(x), otherwise None.
//
// Example:
//
//	o := Flatten(Some(Some(42)))  // Some(42)
func Flatten[T any](o Option[Option[T]]) Option[T] {
	if !o.ok {
		return None[T]()
	}
	return o.value
}

// Tap calls fn with the contained value if o is Some, then returns o unchanged.
//
// Example:
//
//	o := Tap(Some(42), func(x int) { log.Println(x) })
func Tap[T any](o Option[T], fn func(T)) Option[T] {
	return o.Tap(fn)
}

// Match returns onSome(o.value) if o is Some, otherwise returns onNone().
// Exhaustive handling of both branches; returns a single result of type R.
//
// Example:
//
//	s := Match(Some(42), func(x int) string { return fmt.Sprint(x) }, func() string { return "none" })
func Match[T, R any](o Option[T], onSome func(T) R, onNone func() R) R {
	if o.ok {
		return onSome(o.value)
	}
	return onNone()
}

// TryMap applies fn to the contained value if o is Some.
// If o is None, returns (None[U], nil). If fn returns an error, returns (None[U], err).
//
// Example:
//
//	o, err := TryMap(Some("42"), strconv.Atoi)  // Some(42), nil
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
//
// Example:
//
//	Equals(Some(1), Some(1))   // true
//	Equals(None[int](), None[int]())  // true
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
//
// Example:
//
//	p := Zip(Some(1), Some("a")).Unwrap()
//	p.First, p.Second  // 1, "a"
type Pair[A, B any] struct {
	First  A
	Second B
}

// Zip returns Some(Pair{a.value, b.value}) if both a and b are Some, otherwise None.
//
// Example:
//
//	o := Zip(Some(1), Some("a"))  // Some(Pair{1, "a"})
func Zip[T, U any](a Option[T], b Option[U]) Option[Pair[T, U]] {
	if !a.ok || !b.ok {
		return None[Pair[T, U]]()
	}
	return Some(Pair[T, U]{First: a.value, Second: b.value})
}
