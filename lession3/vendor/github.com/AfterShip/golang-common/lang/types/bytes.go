package types

import (
	"bytes"
	"encoding/json"
	"reflect"
)

var (
	_               WrappedType = Bytes{}
	BytesType                   = reflect.TypeOf(Bytes{})
	nativeBytesType             = reflect.TypeOf(Bytes{}.val)
)

// swagger:type string
type Bytes struct {
	val []byte

	isAssigned bool
	isNull     bool
}

// types.Bytes constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewBytes(val []byte, status ...bool) Bytes {
	switch len(status) {
	case 0:
		return newBytes(val, true, false)
	case 1:
		if !status[0] {
			return newBytes(nil, false, false)
		}
		return newBytes(val, status[0], false)
	default:
		if !status[0] {
			return newBytes(nil, false, false)
		}
		if status[1] {
			return newBytes(nil, true, true)
		}
		return newBytes(val, status[0], status[1])
	}
}

func NewNullBytes() Bytes {
	return newBytes(nil, true, true)
}

func newBytes(val []byte, assigned bool, null bool) Bytes {
	return Bytes{
		val:        val,
		isAssigned: assigned,
		isNull:     null,
	}
}

func (n Bytes) Val() interface{} {
	return n.val
}

func (n Bytes) Bytes() []byte {
	return n.Val().([]byte)
}

func (n Bytes) Assigned() bool {
	return n.isAssigned
}

func (n Bytes) Null() bool {
	return n.isNull
}

func (n Bytes) Type() reflect.Type {
	return BytesType
}

func (n Bytes) NativeType() reflect.Type {
	return nativeBytesType
}

func (n *Bytes) CopyFrom(from Bytes, options ...Option) {
	n.val = from.Bytes()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *Bytes) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.val = []byte{}
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		n.val = data
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n Bytes) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}
		return json.Marshal(string(n.val))
	}
	return json.Marshal("")
}
