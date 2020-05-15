package types

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

var (
	_                WrappedType = String{}
	StringType                   = reflect.TypeOf(String{})
	nativeStringType             = reflect.TypeOf(String{}.val)
)

// swagger:type string
type String struct {
	val string

	isAssigned bool
	isNull     bool
}

// types.String constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewString(val string, status ...bool) String {
	switch len(status) {
	case 0:
		return newString(val, true, false)
	case 1:
		if !status[0] {
			return newString("", false, false)
		}
		return newString(val, status[0], false)
	default:
		if !status[0] {
			return newString("", false, false)
		}
		if status[1] {
			return newString("", true, true)
		}
		return newString(val, status[0], status[1])
	}
}

func NewNullString() String {
	return newString("", true, true)
}

func newString(val string, assigned bool, null bool) String {
	return String{
		val:        val,
		isAssigned: assigned,
		isNull:     null,
	}
}

func (n String) Val() interface{} {
	return n.val
}

func (n String) String() string {
	return n.Val().(string)
}

func (n String) Assigned() bool {
	return n.isAssigned
}

func (n String) Null() bool {
	return n.isNull
}

func (n String) Type() reflect.Type {
	return StringType
}

func (n String) NativeType() reflect.Type {
	return nativeStringType
}

func (n *String) CopyFrom(from String, options ...Option) {
	n.val = from.String()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		var err error
		n.val, err = strconv.Unquote(string(data))
		if err != nil {
			data = []byte(strings.Replace(string(data), `\/`, "/", -1))
			n.val, err = strconv.Unquote(string(data))
			if err != nil {
				return err
			}
		}
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n String) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}
		return json.Marshal(n.val)
	}
	return json.Marshal("")
}
