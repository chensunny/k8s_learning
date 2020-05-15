package types

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/civil"
)

var (
	jsonNull  = []byte("null")
	jsonTrue  = []byte("true")
	jsonFalse = []byte("false")

	WrappedTypeInterface = reflect.TypeOf((*WrappedType)(nil)).Elem()
	Null                 = nullValue{}
)

type InvalidTypeValue struct {
	typ reflect.Type
	val interface{}
}

func (v InvalidTypeValue) Error() string {
	return fmt.Sprintf("wrong type of input value, expected: %s, actual: %+v", v.typ.Name(), v.val)
}

type nullValue struct{}

type WrappedType interface {
	Val() interface{}
	Assigned() bool
	Null() bool
	Type() reflect.Type
	NativeType() reflect.Type
}

type Option func(class WrappedType)

func NewWrappedType(typ reflect.Type, val interface{}, status ...bool) (WrappedType, error) {
	switch typ {
	case StringType:
		if len(status) > 1 {
			if status[1] {
				val = ""
			}
		}
		if _, ok := val.(string); ok {
			return NewString(val.(string), status...), nil
		}
	case Int64Type:
		if len(status) > 1 {
			if status[1] {
				val = 0
			}
		}
		switch t := val.(type) {
		case uint64:
			val = int64(t)
		case int32:
			val = int64(t)
		case uint32:
			val = int64(t)
		case int:
			val = int64(t)
		case uint:
			val = int64(t)
		case uint16:
			val = int64(t)
		case int8:
			val = int64(t)
		case uint8:
			val = int64(t)
		default:
			break
		}
		return NewInt64(val.(int64), status...), nil
	case Float64Type:
		if len(status) > 1 {
			if status[1] {
				val = 0.0
			}
		}
		if _, ok := val.(float64); ok {
			return NewFloat64(val.(float64), status...), nil
		}
	case BoolType:
		if len(status) > 1 {
			if status[1] {
				val = false
			}
		}
		if _, ok := val.(bool); ok {
			return NewBool(val.(bool), status...), nil
		}
	case BytesType:
		if len(status) > 1 {
			if status[1] {
				val = []byte(nil)
			}
		}
		if _, ok := val.(string); ok {
			val = []byte(val.(string))
		}
		if _, ok := val.([]byte); ok {
			return NewBytes(val.([]byte), status...), nil
		}
	case DateType:
		if len(status) > 1 {
			if status[1] {
				val = civil.Date{}
			}
		}
		if _, ok := val.(civil.Date); ok {
			return NewDate(val.(civil.Date), status...), nil
		}
	case TimeType:
		if len(status) > 1 {
			if status[1] {
				val = time.Time{}
			}
		}
		if _, ok := val.(time.Time); ok {
			return NewTime(val.(time.Time), status...), nil
		}
	default:
		return nil, errors.New("unsupported wrapped type, type name: " + typ.Name())
	}
	return nil, InvalidTypeValue{typ: typ, val: val}
}

func IsZeroValue(t WrappedType) bool {
	return t.Val() == reflect.New(reflect.TypeOf(t.Val())).Elem().Interface()
}

func NewWrappedTypeReflection(typ reflect.Type, val interface{}, status ...bool) (reflect.Value, error) {
	value, err := NewWrappedType(typ, val, status...)
	if err != nil {
		return reflect.Value{}, err
	}
	switch typ {
	case StringType:
		return reflect.ValueOf(value.(String)), nil
	case Int64Type:
		return reflect.ValueOf(value.(Int64)), nil
	case Float64Type:
		return reflect.ValueOf(value.(Float64)), nil
	case BytesType:
		return reflect.ValueOf(value.(Bytes)), nil
	case BoolType:
		return reflect.ValueOf(value.(Bool)), nil
	case TimeType:
		return reflect.ValueOf(value.(Time)), nil
	case DateType:
		return reflect.ValueOf(value.(Date)), nil
	}
	return reflect.Value{}, ErrUnsupportedFieldType
}
