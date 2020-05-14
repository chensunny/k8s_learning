package types

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

var (
	_                 WrappedType = Float64{}
	Float64Type                   = reflect.TypeOf(Float64{})
	nativeFloat64Type             = reflect.TypeOf(Float64{}.val)
)

// swagger:type float64
type Float64 struct {
	val        float64
	isAssigned bool
	isNull     bool
}

// types.Float64 constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewFloat64(val float64, status ...bool) Float64 {
	switch len(status) {
	case 0:
		return newFloat64(val, true, false)
	case 1:
		if !status[0] {
			return newFloat64(0, false, false)
		}
		return newFloat64(val, status[0], false)
	default:
		if !status[0] {
			return newFloat64(0, false, false)
		}
		if status[1] {
			return newFloat64(0, true, true)
		}
		return newFloat64(val, status[0], status[1])
	}
}

func NewNullFloat64() Float64 {
	return newFloat64(0, true, true)
}

func newFloat64(val float64, assigned bool, null bool) Float64 {
	return Float64{
		val:        val,
		isAssigned: assigned,
		isNull:     null,
	}
}

func (n Float64) Val() interface{} {
	return n.val
}

func (n Float64) Float64() float64 {
	return n.Val().(float64)
}

func (n Float64) Assigned() bool {
	return n.isAssigned
}

func (n Float64) Null() bool {
	return n.isNull
}

func (n Float64) Type() reflect.Type {
	return Float64Type
}

func (n Float64) NativeType() reflect.Type {
	return nativeFloat64Type
}

func (n *Float64) CopyFrom(from Float64, options ...Option) {
	n.val = from.Float64()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *Float64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		tmp, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		n.val = tmp
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n Float64) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}
		return json.Marshal(n.val)
	}
	return json.Marshal(0.0)
}
