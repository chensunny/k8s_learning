package errors

import (
	"errors"
	"testing"

	"golang.org/x/xerrors"
)

var (
	err = errors.New("hello error")
)

func ExampleWrapError(t *testing.T) {
	// can used in any layer
	wrappedErr := Wrap(err, Field("field1", "value1"))
	if xerrors.Is(wrappedErr, err) {
		//do something
	}
}

func ExampleAPIError(t *testing.T) {
	err := WithScene(ErrInternalError, Field("field1", "value1"))
	if xerrors.Is(err, ErrInternalError) {
		//will be true in this case , do some thing
	}
}

func ExampleAPIError2(t *testing.T) {
	//usually for api layer
	err := Wrap(ErrInternalError, Field("field1", "value1"))
	if xerrors.Is(err, ErrInternalError) {
		//will be true in this case , do some thing
	}
}
