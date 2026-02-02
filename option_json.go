package goptions

import (
	"bytes"
	"encoding/json"
)

// MarshalOption marshals o using the given marshal function. None becomes "null";
// Some(v) becomes marshal(v). Use with any JSON lib (stdlib, sonic, etc.).
func MarshalOption[T any](o Option[T], marshal func(T) ([]byte, error)) ([]byte, error) {
	if !o.ok {
		return []byte("null"), nil
	}
	return marshal(o.value)
}

// UnmarshalOption unmarshals data into Option[T] using the given unmarshal function.
// Null or empty input becomes None; otherwise unmarshal into a new T and return Some(t).
// Use with any JSON lib (stdlib, sonic, etc.).
func UnmarshalOption[T any](data []byte, unmarshal func([]byte, *T) error) (Option[T], error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return None[T](), nil
	}
	var t T
	if err := unmarshal(data, &t); err != nil {
		return None[T](), err
	}
	return Some(t), nil
}

// MarshalJSON implements encoding/json.Marshaler. None encodes as null; Some(v) encodes as v.
// T must be JSON-marshalable.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.ok {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements encoding/json.Unmarshaler. Null decodes as None; otherwise decodes into Some(v).
// T must be JSON-unmarshalable.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*o = Option[T]{ok: false}
		return nil
	}
	o.ok = true
	return json.Unmarshal(data, &o.value)
}
