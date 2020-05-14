package types

import (
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func RegisterEncoder() {
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*String)(nil)).Elem().String(), &StringEncoder{})
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*Int64)(nil)).Elem().String(), &Int64Encoder{})
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*Float64)(nil)).Elem().String(), &Float64Encoder{})
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*Bool)(nil)).Elem().String(), &BoolEncoder{})
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*Bytes)(nil)).Elem().String(), &BytesEncoder{})
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*Time)(nil)).Elem().String(), &TimeEncoder{})
	jsoniter.RegisterTypeEncoder(reflect.TypeOf((*Date)(nil)).Elem().String(), &DateEncoder{})
}

type StringEncoder struct{}

func (encoder *StringEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*String)(ptr))
	if !class.Null() {
		stream.WriteVal(class.String())
	} else {
		stream.WriteNil()
	}
}

func (encoder *StringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*String)(ptr))
	return !class.Assigned()
}

func (encoder *StringEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*String)(ptr))
	return !class.Assigned()
}

type Int64Encoder struct{}

func (encoder *Int64Encoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*Int64)(ptr))
	if !class.Null() {
		stream.WriteVal(class.Int64())
	} else {
		stream.WriteNil()
	}
}

func (encoder *Int64Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*Int64)(ptr))
	return !class.Assigned()
}

func (encoder *Int64Encoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*Int64)(ptr))
	return !class.Assigned()
}

type Float64Encoder struct{}

func (encoder *Float64Encoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*Float64)(ptr))
	if !class.Null() {
		stream.WriteVal(class.Float64())
	} else {
		stream.WriteNil()
	}
}

func (encoder *Float64Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*Float64)(ptr))
	return !class.Assigned()
}

func (encoder *Float64Encoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*Float64)(ptr))
	return !class.Assigned()
}

type BoolEncoder struct{}

func (encoder *BoolEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*Bool)(ptr))
	if !class.Null() {
		stream.WriteVal(class.Bool())
	} else {
		stream.WriteNil()
	}
}

func (encoder *BoolEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*Bool)(ptr))
	return !class.Assigned()
}

func (encoder *BoolEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*Bool)(ptr))
	return !class.Assigned()
}

type BytesEncoder struct{}

func (encoder *BytesEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*Bytes)(ptr))
	if !class.Null() {
		stream.WriteVal(string(class.Bytes()))
	} else {
		stream.WriteNil()
	}
}

func (encoder *BytesEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*Bytes)(ptr))
	return !class.Assigned()
}

func (encoder *BytesEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*Bytes)(ptr))
	return !class.Assigned()
}

type TimeEncoder struct{}

func (encoder *TimeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*Time)(ptr))
	if !class.Null() {
		stream.WriteVal(class.Time().Format(rfc3339Colon))
	} else {
		stream.WriteNil()
	}
}

func (encoder *TimeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*Time)(ptr))
	return !class.Assigned()
}

func (encoder *TimeEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*Time)(ptr))
	return !class.Assigned()
}

type DateEncoder struct{}

func (encoder *DateEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	class := *((*Date)(ptr))
	if !class.Null() {
		stream.WriteVal(class.Date().String())
	} else {
		stream.WriteNil()
	}
}

func (encoder *DateEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	class := *((*Date)(ptr))
	return !class.Assigned()
}

func (encoder *DateEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	class := *((*Date)(ptr))
	return !class.Assigned()
}
