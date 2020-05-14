package types

import "errors"

var (
	ErrValueAccessingUnexportedStructFields     = errors.New("Invalid field type.")
	ErrValueNotAddressableOrNotExported         = errors.New("Maybe this field is not addressable or is an unexported struct field.")
	ErrNonNestedStructMustImplementWrapperClass = errors.New("Non-nested struct must implement WrappedType.")
	ErrUnsupportedWrapperClassType              = errors.New("Unsupported WrappedType type.")
	ErrUnsupportedFieldType                     = errors.New("Unsupported field type.")
	ErrNestedSliceMustOnlyContainsStruct        = errors.New("Nested slice must only contains struct.")
	ErrTargetMustBePointerOfStruct              = errors.New("must be pointer of struct")
	ErrOnlyStructCanBeConvertToSpannerMutation  = errors.New("must be struct tagged with 'spanner'")
	ErrNestedStructOrArrayNotSupported          = errors.New("Not support any nested struct/array/slice when convert to spanner mutation")
	ErrUnsupportedSpannerTypeCodeMapping        = errors.New("unsupported spanner type code mapping")
)
