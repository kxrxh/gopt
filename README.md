# gopt

Option type for Go 1.21+. Generics, zero deps, no allocs.

```bash
go get github.com/kxrxh/gopt
```

---

## How to use

```go
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/kxrxh/gopt"
)

func main() {
	// Create
	o := gopt.Some(42)
	o = gopt.None[int]()
	o = gopt.FromPtr(&x)       // nil -> None
	o = gopt.FromTuple(v, ok)  // comma-ok
	o = gopt.Try(val, err)     // (T, error) -> Option[T]
	o = gopt.Cond(ok, v)       // if ok then Some(v) else None

	// Inspect / unwrap
	if o.IsSome() { /* ... */ }
	v, ok := o.Get()
	v = o.UnwrapOr(0)
	v = o.UnwrapOrElse(func() int { return 0 })
	v = o.Expect("must have value")  // panics if None
	p := o.ToPointer()               // nil if None; else *T (copy)

	// Transform (methods + package funcs)
	o2 := gopt.Map(o, func(x int) string { return strconv.Itoa(x) })
	o2 = gopt.AndThen(o, func(x int) gopt.Option[string] { return gopt.Some(strconv.Itoa(x)) })
	o = o.Filter(func(x int) bool { return x > 0 })
	o = o.Or(gopt.Some(99))
	o = o.OrElse(func() gopt.Option[int] { return gopt.Some(99) })
	o = gopt.Flatten(gopt.Some(gopt.Some(1)))  // -> Some(1)
	o.Tap(func(x int) { log.Println(x) })
	result := gopt.Match(o, func(x int) string { return fmt.Sprint(x) }, func() string { return "none" })

	// Map with default
	v = gopt.MapOr(o, 0, func(x int) int { return x * 2 })
	v = gopt.MapOrElse(o, func() int { return 0 }, func(x int) int { return x * 2 })

	// Error in map (e.g. parse string -> int)
	oStr := gopt.Some("42")
	o2, err := gopt.TryMap(oStr, strconv.Atoi)
	_, _ = o2, err

	// Compare / zip
	eq := gopt.Equals(o, other)
	pair := gopt.Zip(o1, o2)  // Option[Pair[T,U]]; pair.First, pair.Second
	_, _ = eq, pair

	// JSON (pluggable — use sonic, stdlib, any lib)
	b, _ := gopt.MarshalOption(o, json.Marshal)
	o, _ = gopt.UnmarshalOption[int](b, func(data []byte, p *int) error { return json.Unmarshal(data, p) })
	// Or with stdlib directly (Option implements json.Marshaler/Unmarshaler):
	b, _ = json.Marshal(o)
	json.Unmarshal(b, &o)
}
```

---

## API reference

**Constructors**

| API | Description |
|-----|-------------|
| `Some(v)` | Option with value. |
| `None[T]()` | Empty option. |
| `FromPtr(p)` | None if p is nil, else Some(*p). |
| `FromTuple(v, ok)` | Some(v) if ok, else None. |
| `Try(v, err)` | Some(v) if err is nil, else None. |
| `Cond(ok, v)` | Some(v) if ok, else None. |

**Inspection**

| API | Description |
|-----|-------------|
| `IsSome()` | True if value present. |
| `IsNone()` | True if empty. |
| `Get()` | (value, ok). |

**Unwrap**

| API | Description |
|-----|-------------|
| `Unwrap()` | Value or panic. |
| `UnwrapOr(default)` | Value or default. |
| `UnwrapOrElse(fn)` | Value or fn(). |
| `Expect(msg)` | Value or panic with msg. |
| `ToPointer()` | nil if None; else *T (copy). |

**Transform** (package funcs; methods exist where types allow)

| API | Description |
|-----|-------------|
| `Map(o, fn)` | Some(fn(v)) or None. |
| `MapOr(o, default, fn)` | fn(v) or default. |
| `MapOrElse(o, defaultFn, fn)` | fn(v) or defaultFn(). |
| `AndThen(o, fn)` | fn(v) (Option) or None. |
| `Filter(o, pred)` | o if pred(v), else None. |
| `Or(o, other)` | o if Some, else other. |
| `OrElse(o, fn)` | o if Some, else fn(). |
| `Flatten(o)` | Option[Option[T]] -> Option[T]. |
| `Tap(o, fn)` | Call fn(v) if Some; return o. |
| `Match(o, onSome, onNone)` | onSome(v) or onNone(). |
| `TryMap(o, fn)` | Map with (U, error); (Option[U], error). |
| `Equals(a, b)` | a == b (both None or both Some with same value). |
| `Zip(a, b)` | Some(Pair{a,b}) if both Some, else None. |

**Pair** (from Zip): `First`, `Second` fields.

**JSON**

| API | Description |
|-----|-------------|
| `MarshalOption(o, marshal)` | None -> "null"; Some(v) -> marshal(v). Use any lib (sonic, json). |
| `UnmarshalOption(data, unmarshal)` | null/empty -> None; else unmarshal into T. |
| `Option` implements `json.Marshaler` / `Unmarshaler` | Works with encoding/json directly. |

---

[pkg.go.dev/github.com/kxrxh/gopt](https://pkg.go.dev/github.com/kxrxh/gopt) · MIT
