package types

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

var (
	_               WrappedType = Int64{}
	Int64Type                   = reflect.TypeOf(Int64{})
	nativeInt64Type             = reflect.TypeOf(Int64{}.val)
)

// swagger:type int64
type Int64 struct {
	val        int64
	isAssigned bool
	isNull     bool
}

// types.Int64 constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewInt64(val int64, status ...bool) Int64 {
	switch len(status) {
	case 0:
		return newInt64(val, true, false)
	case 1:
		if !status[0] {
			return newInt64(0, false, false)
		}
		return newInt64(val, status[0], false)
	default:
		if !status[0] {
			return newInt64(0, false, false)
		}
		if status[1] {
			return newInt64(0, true, true)
		}
		return newInt64(val, status[0], status[1])
	}
}

func NewNullInt64() Int64 {
	return newInt64(0, true, true)
}

func newInt64(val int64, assigned bool, null bool) Int64 {
	return Int64{
		val:        val,
		isAssigned: assigned,
		isNull:     null,
	}
}

func (n Int64) Val() interface{} {
	return n.val
}

func (n Int64) Int64() int64 {
	return n.Val().(int64)
}

func (n Int64) Assigned() bool {
	return n.isAssigned
}

func (n Int64) Null() bool {
	return n.isNull
}

func (n Int64) Type() reflect.Type {
	return Int64Type
}

func (n Int64) NativeType() reflect.Type {
	return nativeInt64Type
}

func (n *Int64) CopyFrom(from Int64, options ...Option) {
	n.val = from.Int64()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *Int64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		tmp, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		n.val = tmp
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n Int64) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}
		return json.Marshal(n.val)
	}
	return json.Marshal(0)
}
