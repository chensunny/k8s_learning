package types

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/AfterShip/golang-common/errors"
)

var (
	_              WrappedType = Bool{}
	BoolType                   = reflect.TypeOf(Bool{})
	nativeBoolType             = reflect.TypeOf(Bool{}.val)
)

// swagger:type bool
type Bool struct {
	val        bool
	isAssigned bool
	isNull     bool
}

// types.Bool constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewBool(val bool, status ...bool) Bool {
	switch len(status) {
	case 0:
		return newBool(val, true, false)
	case 1:
		if !status[0] {
			return newBool(false, false, false)
		}
		return newBool(val, status[0], false)
	default:
		if !status[0] {
			return newBool(false, false, false)
		}
		if status[1] {
			return newBool(false, true, true)
		}
		return newBool(val, status[0], status[1])
	}
}

func NewNullBool() Bool {
	return newBool(false, true, true)
}

func newBool(val bool, isAssigned, isNull bool) Bool {
	return Bool{val: val, isAssigned: isAssigned, isNull: isNull}
}

func (n Bool) Val() interface{} {
	return n.val
}

func (n Bool) Bool() bool {
	return n.Val().(bool)
}

func (n Bool) Assigned() bool {
	return n.isAssigned
}

func (n Bool) Null() bool {
	return n.isNull
}

func (n Bool) Type() reflect.Type {
	return BoolType
}

func (n Bool) NativeType() reflect.Type {
	return nativeBoolType
}

func (n *Bool) CopyFrom(from Bool, options ...Option) {
	n.val = from.Bool()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		if bytes.Equal(data, jsonTrue) {
			n.val = true
		} else if bytes.Equal(data, jsonFalse) {
			n.val = false
		} else {
			return errors.WithScene(errors.ErrUnprocessableEntity)
		}
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n Bool) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}
		return json.Marshal(n.val)
	}
	return json.Marshal(false)
}
