package gopt

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
)

func TestSome(t *testing.T) {
	o := Some(42)
	if !o.IsSome() || o.IsNone() {
		t.Fatal("Some(42) should be Some")
	}
	if v, ok := o.Get(); !ok || v != 42 {
		t.Fatalf("Get() = %v, %v; want 42, true", v, ok)
	}
	if o.Unwrap() != 42 {
		t.Fatalf("Unwrap() = %v; want 42", o.Unwrap())
	}
}

func TestNone(t *testing.T) {
	o := None[int]()
	if o.IsSome() || !o.IsNone() {
		t.Fatal("None() should be None")
	}
	if _, ok := o.Get(); ok {
		t.Fatal("Get() on None should return false")
	}
}

func TestFromPtr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		o := FromPtr[int](nil)
		if o.IsSome() {
			t.Fatal("FromPtr(nil) should be None")
		}
	})
	t.Run("non-nil", func(t *testing.T) {
		x := 7
		o := FromPtr(&x)
		if !o.IsSome() {
			t.Fatal("FromPtr(&x) should be Some")
		}
		if o.Unwrap() != 7 {
			t.Fatalf("Unwrap() = %v; want 7", o.Unwrap())
		}
	})
}

func TestFromTuple(t *testing.T) {
	o1 := FromTuple(1, true)
	if !o1.IsSome() || o1.Unwrap() != 1 {
		v, ok := o1.Get()
		t.Fatalf("FromTuple(1, true): got %v, %v", v, ok)
	}
	o2 := FromTuple(0, false)
	if o2.IsSome() {
		t.Fatal("FromTuple(0, false) should be None")
	}
}

func TestTry(t *testing.T) {
	o1 := Try(42, nil)
	if !o1.IsSome() || o1.Unwrap() != 42 {
		t.Fatalf("Try(42, nil): want Some(42)")
	}
	o2 := Try(0, errors.New("err"))
	if o2.IsSome() {
		t.Fatal("Try(0, err) should be None")
	}
}

func TestUnwrapOr(t *testing.T) {
	if Some(1).UnwrapOr(99) != 1 {
		t.Fatal("Some(1).UnwrapOr(99) should be 1")
	}
	if None[int]().UnwrapOr(99) != 99 {
		t.Fatal("None().UnwrapOr(99) should be 99")
	}
}

func TestUnwrapOrElse(t *testing.T) {
	if Some(1).UnwrapOrElse(func() int { return 99 }) != 1 {
		t.Fatal("Some(1).UnwrapOrElse(...) should be 1")
	}
	if None[int]().UnwrapOrElse(func() int { return 99 }) != 99 {
		t.Fatal("None().UnwrapOrElse(...) should be 99")
	}
}

func TestExpect(t *testing.T) {
	if Some(1).Expect("x") != 1 {
		t.Fatal("Some(1).Expect(...) should be 1")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("Expect on None should panic")
			}
		}()
		None[int]().Expect("expected panic")
	}()
}

func TestMap(t *testing.T) {
	// Package-level Map
	o := Some(21)
	m := Map(o, func(x int) int { return x * 2 })
	if !m.IsSome() || m.Unwrap() != 42 {
		v, ok := m.Get()
		t.Fatalf("Map(Some(21), *2) = %v, %v; want Some(42)", v, ok)
	}
	none := Map(None[int](), func(x int) string { return strconv.Itoa(x) })
	if none.IsSome() {
		t.Fatal("Map(None, ...) should be None")
	}
	m2 := Map(Some(10), func(x int) string { return strconv.Itoa(x) })
	if !m2.IsSome() || m2.Unwrap() != "10" {
		v, ok := m2.Get()
		t.Fatalf("Map(Some(10), ItoA) = %v, %v; want Some(\"10\")", v, ok)
	}
}

func TestAndThen(t *testing.T) {
	sq := func(x int) Option[int] {
		if x < 0 {
			return None[int]()
		}
		return Some(x * x)
	}
	if v := AndThen(Some(4), sq).Unwrap(); v != 16 {
		t.Fatalf("AndThen(Some(4), sq) = %v; want 16", v)
	}
	if AndThen(Some(-1), sq).IsSome() {
		t.Fatal("AndThen(Some(-1), sq) should be None")
	}
	if AndThen(None[int](), sq).IsSome() {
		t.Fatal("AndThen(None, sq) should be None")
	}
}

func TestFilter(t *testing.T) {
	even := func(x int) bool { return x%2 == 0 }
	if v := Some(4).Filter(even).Unwrap(); v != 4 {
		t.Fatalf("Some(4).Filter(even) = %v; want 4", v)
	}
	if Some(3).Filter(even).IsSome() {
		t.Fatal("Some(3).Filter(even) should be None")
	}
	if None[int]().Filter(even).IsSome() {
		t.Fatal("None().Filter(even) should be None")
	}
	// Package-level
	if v := Filter(Some(2), even).Unwrap(); v != 2 {
		t.Fatalf("Filter(Some(2), even) = %v; want 2", v)
	}
}

func TestOr(t *testing.T) {
	if !Some(1).Or(Some(2)).IsSome() || Some(1).Or(Some(2)).Unwrap() != 1 {
		t.Fatal("Some(1).Or(Some(2)) should be Some(1)")
	}
	if !None[int]().Or(Some(2)).IsSome() || None[int]().Or(Some(2)).Unwrap() != 2 {
		t.Fatal("None().Or(Some(2)) should be Some(2)")
	}
	if None[int]().Or(None[int]()).IsSome() {
		t.Fatal("None().Or(None()) should be None")
	}
}

func TestOrElse(t *testing.T) {
	if !Some(1).OrElse(func() Option[int] { return Some(2) }).IsSome() ||
		Some(1).OrElse(func() Option[int] { return Some(2) }).Unwrap() != 1 {
		t.Fatal("Some(1).OrElse(...) should be Some(1)")
	}
	if !None[int]().OrElse(func() Option[int] { return Some(2) }).IsSome() ||
		None[int]().OrElse(func() Option[int] { return Some(2) }).Unwrap() != 2 {
		t.Fatal("None().OrElse(...) should be Some(2)")
	}
}

func TestFlatten(t *testing.T) {
	if v := Flatten(Some(Some(42))).Unwrap(); v != 42 {
		t.Fatalf("Flatten(Some(Some(42))) = %v; want 42", v)
	}
	if Flatten(None[Option[int]]()).IsSome() {
		t.Fatal("Flatten(None[Option[int]]) should be None")
	}
	if Flatten(Some(None[int]())).IsSome() {
		t.Fatal("Flatten(Some(None())) should be None")
	}
}

func TestTap(t *testing.T) {
	var got int
	o := Some(7).Tap(func(x int) { got = x })
	if got != 7 || !o.IsSome() || o.Unwrap() != 7 {
		v, ok := o.Get()
		t.Fatalf("Tap: got=%v, o=%v,%v; want got=7, Some(7)", got, v, ok)
	}
	None[int]().Tap(func(x int) { got = -1 })
	if got != 7 {
		t.Fatalf("Tap on None should not call fn; got=%v", got)
	}
}

func TestMatch(t *testing.T) {
	r1 := Match(Some(10), func(x int) string { return strconv.Itoa(x) }, func() string { return "none" })
	if r1 != "10" {
		t.Fatalf("Match(Some(10), ...) = %q; want \"10\"", r1)
	}
	r2 := Match(None[int](), func(x int) string { return strconv.Itoa(x) }, func() string { return "none" })
	if r2 != "none" {
		t.Fatalf("Match(None, ...) = %q; want \"none\"", r2)
	}
}

func TestTryMap(t *testing.T) {
	parse := func(s string) (int, error) { return strconv.Atoi(s) }
	o, err := TryMap(Some("42"), parse)
	if err != nil || !o.IsSome() || o.Unwrap() != 42 {
		v, ok := o.Get()
		t.Fatalf("TryMap(Some(\"42\"), Atoi) = %v, %v, %v; want Some(42), nil", v, ok, err)
	}
	o2, err2 := TryMap(Some("x"), parse)
	if err2 == nil || o2.IsSome() {
		v, ok := o2.Get()
		t.Fatalf("TryMap(Some(\"x\"), Atoi) should return error and None; got %v, %v, %v", v, ok, err2)
	}
	o3, err3 := TryMap(None[string](), parse)
	if err3 != nil || o3.IsSome() {
		v, ok := o3.Get()
		t.Fatalf("TryMap(None, ...) = %v, %v, %v; want None, nil", v, ok, err3)
	}
}

func TestEquals(t *testing.T) {
	if !Equals(Some(1), Some(1)) {
		t.Fatal("Equals(Some(1), Some(1)) should be true")
	}
	if Equals(Some(1), Some(2)) {
		t.Fatal("Equals(Some(1), Some(2)) should be false")
	}
	if !Equals(None[int](), None[int]()) {
		t.Fatal("Equals(None, None) should be true")
	}
	if Equals(Some(1), None[int]()) || Equals(None[int](), Some(1)) {
		t.Fatal("Equals(Some, None) should be false")
	}
	// Float NaN: NaN != NaN in Go, so Equals(Some(NaN), Some(NaN)) is false.
	if Equals(Some(math.NaN()), Some(math.NaN())) {
		t.Fatal("Equals(Some(NaN), Some(NaN)) should be false (NaN != NaN)")
	}
}

func TestZip(t *testing.T) {
	p := Zip(Some(1), Some("a"))
	if !p.IsSome() {
		t.Fatal("Zip(Some(1), Some(\"a\")) should be Some")
	}
	pair := p.Unwrap()
	if pair.First != 1 || pair.Second != "a" {
		t.Fatalf("Zip result = %v, %v; want 1, \"a\"", pair.First, pair.Second)
	}
	if Zip(None[int](), Some("a")).IsSome() {
		t.Fatal("Zip(None, Some) should be None")
	}
	if Zip(Some(1), None[string]()).IsSome() {
		t.Fatal("Zip(Some, None) should be None")
	}
	if Zip(None[int](), None[string]()).IsSome() {
		t.Fatal("Zip(None, None) should be None")
	}
}

func TestMapOr(t *testing.T) {
	if v := MapOr(Some(10), 0, func(x int) int { return x * 2 }); v != 20 {
		t.Fatalf("MapOr(Some(10), 0, *2) = %v; want 20", v)
	}
	if v := MapOr(None[int](), 99, func(x int) int { return x * 2 }); v != 99 {
		t.Fatalf("MapOr(None, 99, *2) = %v; want 99", v)
	}
}

func TestMapOrElse(t *testing.T) {
	if v := MapOrElse(Some(5), func() int { return 0 }, func(x int) int { return x + 1 }); v != 6 {
		t.Fatalf("MapOrElse(Some(5), ...) = %v; want 6", v)
	}
	if v := MapOrElse(None[int](), func() int { return 42 }, func(x int) int { return x }); v != 42 {
		t.Fatalf("MapOrElse(None, 42, ...) = %v; want 42", v)
	}
}

func TestCond(t *testing.T) {
	if !Cond(true, 7).IsSome() || Cond(true, 7).Unwrap() != 7 {
		t.Fatal("Cond(true, 7) should be Some(7)")
	}
	if Cond(false, 7).IsSome() {
		t.Fatal("Cond(false, 7) should be None")
	}
}

func TestToPointer(t *testing.T) {
	if p := None[int]().ToPointer(); p != nil {
		t.Fatalf("None().ToPointer() = %v; want nil", p)
	}
	o := Some(11)
	p := o.ToPointer()
	if p == nil || *p != 11 {
		t.Fatalf("Some(11).ToPointer() = %v; want ptr to 11", p)
	}
	*p = 22
	if o.Unwrap() != 11 {
		t.Fatal("mutating ToPointer result should not change original Option")
	}
}

func TestMarshalOption(t *testing.T) {
	// None -> "null"
	b, err := MarshalOption(None[int](), func(v int) ([]byte, error) { return json.Marshal(v) })
	if err != nil || string(b) != "null" {
		t.Fatalf("MarshalOption(None) = %q, %v; want \"null\", nil", b, err)
	}
	// Some -> marshaled value
	b, err = MarshalOption(Some(42), func(v int) ([]byte, error) { return json.Marshal(v) })
	if err != nil || string(b) != "42" {
		t.Fatalf("MarshalOption(Some(42)) = %q, %v; want \"42\", nil", b, err)
	}
}

func TestUnmarshalOption(t *testing.T) {
	unmarshalInt := func(data []byte, p *int) error { return json.Unmarshal(data, p) }
	// "null" -> None
	o, err := UnmarshalOption([]byte("null"), unmarshalInt)
	if err != nil || o.IsSome() {
		v, ok := o.Get()
		t.Fatalf("UnmarshalOption(null) = %v, %v, %v; want None, nil", v, ok, err)
	}
	// "42" -> Some(42)
	o, err = UnmarshalOption([]byte("42"), unmarshalInt)
	if err != nil || !o.IsSome() || o.Unwrap() != 42 {
		v, ok := o.Get()
		t.Fatalf("UnmarshalOption(42) = %v, %v, %v; want Some(42), nil", v, ok, err)
	}
	// Empty string "" for Option[string] -> Some(""), not None
	unmarshalStr := func(data []byte, p *string) error { return json.Unmarshal(data, p) }
	os, err := UnmarshalOption([]byte(`""`), unmarshalStr)
	if err != nil || !os.IsSome() || os.Unwrap() != "" {
		v, ok := os.Get()
		t.Fatalf("UnmarshalOption with empty string = %v, %v, %v; want Some(\"\"), nil", v, ok, err)
	}
	// Invalid JSON for T ([] or {}) returns error
	_, err = UnmarshalOption([]byte("[]"), unmarshalInt)
	if err == nil {
		t.Fatal("UnmarshalOption([]) for Option[int] should return error")
	}
	_, err = UnmarshalOption([]byte("{}"), unmarshalInt)
	if err == nil {
		t.Fatal("UnmarshalOption({}) for Option[int] should return error")
	}
}

func TestOptionMarshalJSON(t *testing.T) {
	o := None[int]()
	b, err := json.Marshal(o)
	if err != nil || string(b) != "null" {
		t.Fatalf("json.Marshal(None) = %q, %v; want \"null\", nil", b, err)
	}
	o2 := Some(99)
	b, err = json.Marshal(o2)
	if err != nil || string(b) != "99" {
		t.Fatalf("json.Marshal(Some(99)) = %q, %v; want \"99\", nil", b, err)
	}
}

func TestOptionUnmarshalJSON(t *testing.T) {
	var o Option[int]
	if err := json.Unmarshal([]byte("null"), &o); err != nil || o.IsSome() {
		v, ok := o.Get()
		t.Fatalf("json.Unmarshal(null) = %v, %v, %v; want None, nil", v, ok, err)
	}
	if err := json.Unmarshal([]byte("100"), &o); err != nil || !o.IsSome() || o.Unwrap() != 100 {
		v, ok := o.Get()
		t.Fatalf("json.Unmarshal(100) = %v, %v, %v; want Some(100), nil", v, ok, err)
	}
	// Invalid JSON for T ([] or {}) returns error
	if err := json.Unmarshal([]byte("[]"), &o); err == nil {
		t.Fatal("json.Unmarshal([]) for Option[int] should return error")
	}
	if err := json.Unmarshal([]byte("{}"), &o); err == nil {
		t.Fatal("json.Unmarshal({}) for Option[int] should return error")
	}
}

func TestUnwrapPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Unwrap on None should panic")
		}
	}()
	None[int]().Unwrap()
}
