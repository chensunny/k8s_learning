package errors

import (
	"encoding/json"

	"golang.org/x/xerrors"
)

// APIError
func (e *APIError) Unwrap() error {
	if e.scene != nil {
		return e.scene.cause
	}
	return nil
}

func (e *APIError) Is(targetErr error) bool {
	if targetErr == nil {
		return false
	}
	if targetStructErr, ok := targetErr.(*APIError); ok {
		return apiErrorIs(e, targetStructErr)
	}
	return false
}

func apiErrorIs(err *APIError, targetErr *APIError) bool {
	if targetErr.subCode == nil || targetErr.SubCode().code == 0 {
		//target err没有定义sub code的时候，说明是target error 是一个大类，符合大类的都算
		return err.MainCode().Code() == targetErr.MainCode().code
	}
	return err.MainCode().code == targetErr.MainCode().code && err.SubCode().code == targetErr.SubCode().code
}

///// error with scene
func Wrap(err error, sceneFields ...SceneField) *ErrorWithScene {
	scene := &Scene{}
	for _, field := range sceneFields {
		field(scene)
	}
	if len(scene.stack) == 0 {
		scene.stack = SmallerStacktrace(1, 1)
	}
	return &ErrorWithScene{
		cause: err,
		scene: scene,
	}
}

type ErrorWithScene struct {
	cause error
	scene *Scene
}

func (e *ErrorWithScene) Error() string {
	return e.JsonString()
}

func (e *ErrorWithScene) Err() error {
	return e.cause
}

func (e *ErrorWithScene) Cause() error {
	return e.cause
}

func (e *ErrorWithScene) Scene() *Scene {
	if e.scene != nil {
		return e.scene
	} else {
		return emptyScene
	}
}

func (e *ErrorWithScene) JsonString() string {
	data, err := json.Marshal(e)
	if err == nil {
		return string(data)
	} else {
		return ""
	}
}

func (e *ErrorWithScene) Unwrap() error {
	return e.cause
}

func (e *ErrorWithScene) Is(targetErr error) bool {
	if targetErr == nil {
		return false
	}
	if e.cause == targetErr {
		return true
	}
	if e.scene != nil && e.scene.cause != nil {
		return xerrors.Is(e.scene.cause, targetErr)
	}
	return false
}
