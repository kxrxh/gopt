# goptions

Option type for Go 1.21+. Generics, zero deps, no allocs.

```bash
go get github.com/kxrxh/goptions
```

---

## How to use

```go
import (
 "encoding/json"
 "fmt"
 "log"
 "strconv"

 "github.com/kxrxh/goptions"
)

// Create
o := goptions.Some(42)
o = goptions.None[int]()
o = goptions.FromPtr(&x)           // nil -> None
o = goptions.FromTuple(v, ok)     // comma-ok
o = goptions.Try(val, err)        // (T, error) -> Option[T]
o = goptions.Cond(ok, v)          // if ok then Some(v) else None

// Inspect / unwrap
if o.IsSome() { ... }
v, ok := o.Get()
v = o.UnwrapOr(0)
v = o.UnwrapOrElse(func() int { return 0 })
v = o.Expect("must have value")   // panics if None
p := o.ToPointer()                // nil if None; else *T (copy)

// Transform (methods + package funcs)
o2 := goptions.Map(o, func(x int) string { return strconv.Itoa(x) })
o2 = goptions.AndThen(o, func(x int) goptions.Option[string] { return goptions.Some(strconv.Itoa(x)) })
o = o.Filter(func(x int) bool { return x > 0 })
o = o.Or(goptions.Some(99))
o = o.OrElse(func() goptions.Option[int] { return goptions.Some(99) })
o = goptions.Flatten(goptions.Some(goptions.Some(1)))  // -> Some(1)
o.Tap(func(x int) { log.Println(x) })
result := goptions.Match(o, func(x int) string { return fmt.Sprint(x) }, func() string { return "none" })

// Map with default
v = goptions.MapOr(o, 0, func(x int) int { return x * 2 })
v = goptions.MapOrElse(o, func() int { return 0 }, func(x int) int { return x * 2 })

// Error in map (e.g. parse string -> int)
oStr := goptions.Some("42")
o2, err := goptions.TryMap(oStr, strconv.Atoi)

// Compare / zip
eq := goptions.Equals(o, other)
pair := goptions.Zip(o1, o2)  // Option[Pair[T,U]]; pair.First, pair.Second

// JSON (pluggable — use sonic, stdlib, any lib)
b, _ := goptions.MarshalOption(o, json.Marshal)
o, _ = goptions.UnmarshalOption[int](b, func(data []byte, p *int) error { return json.Unmarshal(data, p) })
// Or with stdlib directly (Option implements json.Marshaler/Unmarshaler):
b, _ = json.Marshal(o)
json.Unmarshal(b, &o)
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

[pkg.go.dev/github.com/kxrxh/goptions](https://pkg.go.dev/github.com/kxrxh/goptions) · MIT
