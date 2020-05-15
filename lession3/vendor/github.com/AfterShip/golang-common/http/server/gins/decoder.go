package gins

import (
	json2 "encoding/json"
	"reflect"
	"unsafe"

	"github.com/AfterShip/golang-common/lang/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type autoSetNullExtension struct {
	jsoniter.DummyExtension
}

type autoSetNullDecoder struct {
	typ        reflect2.Type
	valDecoder jsoniter.ValDecoder
}

func (ext *autoSetNullExtension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	if !typ.Type1().Implements(types.WrappedTypeInterface) && typeIn(typ, []reflect.Kind{reflect.Struct, reflect.Array, reflect.Slice, reflect.Map}) {
		return &autoSetNullDecoder{valDecoder: decoder, typ: typ}
	}
	return decoder
}

func (decoder *autoSetNullDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.WhatIsNext() == jsoniter.NilValue {
		// Process json null on any struct/slice/array/map which not implement `lang.types`.
		decoder.traversalNull(decoder.typ, ptr, iter, false)
	} else {
		decoder.valDecoder.Decode(ptr, iter)
	}
}

func (decoder *autoSetNullDecoder) traversalNull(typ reflect2.Type, ptr unsafe.Pointer, iter *jsoniter.Iterator, recursive bool) {
	var typ1 reflect.Type
	var typeKind reflect.Kind
	defer func() {
		if !recursive {
			iter.Skip()
		}
	}()
	typ1 = typ.Type1()
kind:
	typeKind = typ1.Kind()
	switch typeKind {
	case reflect.Array, reflect.Slice:
		newIter := iter.Pool().BorrowIterator([]byte("[]"))
		defer iter.Pool().ReturnIterator(newIter)
		decoder.valDecoder.Decode(ptr, newIter)
	case reflect.Map:
		newIter := iter.Pool().BorrowIterator([]byte("{}"))
		defer iter.Pool().ReturnIterator(newIter)
		decoder.valDecoder.Decode(ptr, newIter)
	case reflect.Struct:
		structType := typ.(reflect2.StructType)
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			if field.Type().Type1().Implements(types.WrappedTypeInterface) {
				val := field.UnsafeGet(ptr)
				tmp := reflect.NewAt(field.Type().Type1(), val).Interface()
				_ = tmp.(json2.Unmarshaler).UnmarshalJSON([]byte("null"))
				p := reflect.ValueOf(tmp).Pointer()
				field.UnsafeSet(ptr, unsafe.Pointer(p))
				continue
			}
			val := field.UnsafeGet(ptr)
			decoder.traversalNull(field.Type(), val, iter, true)
		}
	case reflect.Ptr:
		typ1 = typ1.Elem()
		typ = reflect2.Type2(typ1)
		goto kind
	case reflect.String:
		*((*string)(ptr)) = ""
	case reflect.Int:
		*((*int)(ptr)) = 0
	case reflect.Int8:
		*((*int8)(ptr)) = int8(0)
	case reflect.Int16:
		*((*int16)(ptr)) = int16(0)
	case reflect.Int32:
		*((*int32)(ptr)) = int32(0)
	case reflect.Int64:
		*((*int64)(ptr)) = int64(0)
	case reflect.Uint:
		*((*uint)(ptr)) = uint(0)
	case reflect.Uint8:
		*((*uint8)(ptr)) = uint8(0)
	case reflect.Uint16:
		*((*uint16)(ptr)) = uint16(0)
	case reflect.Uint32:
		*((*uint32)(ptr)) = uint32(0)
	case reflect.Uint64:
		*((*uint64)(ptr)) = uint64(0)
	case reflect.Bool:
		*((*bool)(ptr)) = false
	case reflect.Float64:
		*((*float64)(ptr)) = float64(0)
	case reflect.Float32:
		*((*float32)(ptr)) = float32(0)
	default:
		decoder.valDecoder.Decode(ptr, iter)
	}
}

func typeIn(typ reflect2.Type, types []reflect.Kind) bool {
	for _, t := range types {
		if typ.Kind() == t {
			return true
		}
	}
	return false
}
