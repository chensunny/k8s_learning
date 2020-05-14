package types

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"cloud.google.com/go/civil"
)

var (
	_              WrappedType = Date{}
	DateType                   = reflect.TypeOf(Date{})
	nativeDateType             = reflect.TypeOf(Date{}.val)
)

// swagger:type date
type Date struct {
	val        civil.Date
	isAssigned bool
	isNull     bool
}

// types.Date constructor.
// status[0]: isAssigned, whether field exists in JSON payloads or not.
// status[1]: isNull, whether field is JSON-Null / Spanner-Null / any other null value or not
func NewDate(val civil.Date, status ...bool) Date {
	switch len(status) {
	case 0:
		return newDate(val, true, false)
	case 1:
		if !status[0] {
			return newDate(civil.Date{}, false, false)
		}
		return newDate(val, status[0], false)
	default:
		if !status[0] {
			return newDate(civil.Date{}, false, false)
		}
		if status[1] {
			return newDate(civil.Date{}, true, true)
		}
		return newDate(val, status[0], status[1])
	}
}

func NewNullDate() Date {
	return newDate(civil.Date{}, true, true)
}

func newDate(val civil.Date, assigned bool, null bool) Date {
	return Date{
		val:        val,
		isAssigned: assigned,
		isNull:     null,
	}
}

func (n Date) Val() interface{} {
	return n.val
}

func (n Date) Date() civil.Date {
	return n.Val().(civil.Date)
}

func (n Date) Assigned() bool {
	return n.isAssigned
}

func (n Date) Null() bool {
	return n.isNull
}

func (n Date) Type() reflect.Type {
	return DateType
}

func (n Date) NativeType() reflect.Type {
	return nativeDateType
}

func (n *Date) CopyFrom(from Date, options ...Option) {
	n.val = from.Date()
	n.isAssigned = from.Assigned()
	n.isNull = from.Null()

	if len(options) > 0 {
		for _, opt := range options {
			opt(n)
		}
	}
}

func (n *Date) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		n.isAssigned = true
		n.isNull = true
		return nil
	}
	if len(data) > 0 {
		timeStr, err := strconv.Unquote(string(data))
		if err != nil {
			return err
		}
		n.val, err = civil.ParseDate(timeStr)
		if err != nil {
			return err
		}
		n.isAssigned = true
		n.isNull = false
	}
	return nil
}

func (n Date) MarshalJSON() ([]byte, error) {
	if n.isAssigned {
		if n.isNull {
			return jsonNull, nil
		}

		return json.Marshal(n.val.String())
	}
	return json.Marshal("")
}
