package errors

import (
	"context"
	"encoding/json"
	"golang.org/x/xerrors"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	vv9 "gopkg.in/go-playground/validator.v9"
)

var (
	//统一的0 子码
	SubCodeZero = &Code{code: 0, messageKey: ""}
)

var (
	//Bad Request
	CodeBadRequest = &Code{code: 400, messageKey: "errors.bad_request"}
	ErrBadRequest  = &APIError{mainCode: CodeBadRequest}

	//Unauthorized，没有传递认证信息或认证信息无效
	CodeUnauthorized    = &Code{code: 401, messageKey: "errors.unauthorized"}
	ErrUnauthorized     = &APIError{mainCode: CodeUnauthorized}
	ErrInvalidPrivilege = &APIError{mainCode: CodeUnauthorized, subCode: &Code{code: 1, messageKey: "errors.invalid_privileges"}}
	ErrInvalidToken     = &APIError{mainCode: CodeUnauthorized, subCode: &Code{code: 2, messageKey: "errors.invalid_token"}}
	ErrTokenExpired     = &APIError{mainCode: CodeUnauthorized, subCode: &Code{code: 3, messageKey: "errors.token_expired"}}
	// subcode = 4 means verify signature failed, but messageKey won`t tell client this.
	// Because we want response a fuzzy message against baleful actions.
	ErrTokenVerifyingFailed = &APIError{mainCode: CodeUnauthorized, subCode: &Code{code: 4, messageKey: "errors.invalid_token"}}
	ErrTokenNotContainsUser = &APIError{mainCode: CodeUnauthorized, subCode: &Code{code: 5, messageKey: "errors.invalid_token"}}

	//Payment Required
	CodePaymentRequired = &Code{code: 402, messageKey: "errors.payment.required"}

	//Forbidden
	CodeForbidden = &Code{code: 403, messageKey: "errors.forbidden"}

	//API或业务记录不存在
	CodeNotFound = &Code{code: 404, messageKey: "errors.notfound"}
	ErrNotFound  = &APIError{mainCode: CodeNotFound}

	//API method不允许调用
	CodeMethodNotAllowed = &Code{code: 405, messageKey: "errors.method_not_allowed"}
	ErrMethodNotAllowed  = &APIError{mainCode: CodeMethodNotAllowed}

	//Conflict
	CodeConflict = &Code{code: 409, messageKey: "errors.conflict"}
	ErrConflict  = &APIError{mainCode: CodeConflict}

	//Precondition Failed
	CodePreconditionFailed = &Code{code: 412, messageKey: "errors.precondition_failed"}
	ErrPreconditionFailed  = &APIError{mainCode: CodePreconditionFailed}

	//UnprocessableEntity
	CodeUnprocessableEntity = &Code{code: 422, messageKey: "errors.unprocessable_entity"}
	ErrUnprocessableEntity  = &APIError{mainCode: CodeUnprocessableEntity}
	CodeInvalidArgument     = CodeUnprocessableEntity
	ErrInvalidArgument      = &APIError{mainCode: CodeInvalidArgument}

	//TooManyRequests
	CodeTooManyRequests = &Code{code: 429, messageKey: "errors.too_many_requests"}
	ErrTooManyRequests  = &APIError{mainCode: CodeTooManyRequests}

	//系统内部错误
	CodeInternalError    = &Code{code: 500, messageKey: "errors.internal"}
	ErrInternalError     = &APIError{mainCode: CodeInternalError}
	ErrS3PostPolicyError = &APIError{mainCode: CodeInternalError, subCode: &Code{code: 90, messageKey: "errors.s3.invalid_post_policy"}}

	//系统不可用
	CodeUnavailable = &Code{code: 503, messageKey: "errors.unavailable"}
	ErrUnavailable  = &APIError{mainCode: CodeUnavailable}

	//处理超时
	CodeDeadlineExceeded = &Code{code: 504, messageKey: "errors.deadline_exceeded"}
	ErrDeadlineExceeded  = &APIError{mainCode: CodeDeadlineExceeded}
)

//带error code、scene的自定义error类型，方便对error进行分类处理，常规用于API层
//比如方法调用参数非法、依赖资源访问失败，是需要能识别进行不同的处理
//约定code为2级结构,main code、sub code，方便分类统计，比如code 为 500内部错误、sub code为 001 访问数据库失败
type APIError struct {
	mainCode *Code
	subCode  *Code
	scene    *Scene
}

//错误码
type Code struct {
	code       int
	messageKey string
}

var emptyCode = &Code{}

func (e *APIError) Error() string {
	return e.JsonString()
}

func (e *APIError) JsonString() string {
	data, err := json.Marshal(e)
	if err == nil {
		return string(data)
	} else {
		return ""
	}
}

func (e *APIError) String() string {
	return e.JsonString()
}

func (e *APIError) MainCode() *Code {
	if e.mainCode != nil {
		return e.mainCode
	} else {
		return emptyCode
	}
}

func (e *APIError) SubCode() *Code {
	if e.subCode != nil {
		return e.subCode
	} else {
		return emptyCode
	}
}

func (e *APIError) Scene() *Scene {
	if e.scene != nil {
		return e.scene
	} else {
		return emptyScene
	}
}

func (c *Code) Code() int {
	return c.code
}

func (c *Code) MessageKey() string {
	return c.messageKey
}

func (c *Code) JsonString() string {
	data, err := json.Marshal(c)
	if err == nil {
		return string(data)
	} else {
		return ""
	}
}

func (c *Code) String() string {
	return c.JsonString()
}

// =========================
// APIError builder functions
// =========================

func NewAPIError(mainCode *Code, subCode *Code) *APIError {
	return &APIError{
		mainCode: mainCode,
		subCode:  subCode,
	}
}

// Deprecated
func NewError(mainCode *Code, subCode *Code) *APIError {
	return NewAPIError(mainCode, subCode)
}

func NewCode(code int, messageKey string) *Code {
	return &Code{
		code:       code,
		messageKey: messageKey,
	}
}

func WithSubCode(basic *APIError, subCode *Code) *APIError {
	err := *basic
	err.subCode = subCode
	return &err
}

// Deprecated
func WithScene(basic *APIError, sceneFields ...SceneField) *APIError {
	err := *basic
	if err.scene != nil {
		scene := *err.scene
		err.scene = &scene
	} else {
		err.scene = &Scene{}
	}

	for _, field := range sceneFields {
		field(err.scene)
	}
	if err.scene.stack == "" {
		err.scene.stack = SmallerStacktrace(1, 1)
	}
	return &err
}

// 作用：基于自定义好的APIError，包装好对应的错误。
// 注意：用于在api层调用，不要在业务层和其它底层的调用。
//	可以建议自定义错误 --> APIError的映射，以便明确发生了那些错误应该返回给到API调用方具体什么APIError
// 参考：https://docs.google.com/document/d/1jU2TlEaf0-m7Z9i5h5L780bChKlaNjgOLmFZkkbJcog/edit#
func APIErrorWithScene(basic *APIError, sceneFields ...SceneField) *APIError {
	err := *basic
	if err.scene != nil {
		scene := *err.scene
		err.scene = &scene
	} else {
		err.scene = &Scene{}
	}

	for _, field := range sceneFields {
		field(err.scene)
	}
	if err.scene.stack == "" {
		err.scene.stack = SmallerStacktrace(1, 1)
	}
	return &err
}

// =========================
// APIError parser
// =========================
func ConvertToAPIError(rawErr error) *APIError {
	if xerrors.Is(rawErr, context.DeadlineExceeded) {
		return APIErrorWithScene(ErrDeadlineExceeded, Cause(rawErr), Stack(SmallerStacktrace(1, 1)))
	}

	items := ParseUnprocessableEntityErrorItems(rawErr)
	if items != nil {
		// 保留链路信息
		return APIErrorWithScene(ErrUnprocessableEntity, Cause(rawErr), Items(items), Stack(SmallerStacktrace(1, 1)))
	}

	switch err := rawErr.(type) {
	case *APIError:
		return err
	case *ErrorWithScene:
		return toAPIErrorWithScene(err)
	default:
		return APIErrorWithScene(ErrInternalError, Cause(rawErr), Stack(SmallerStacktrace(1, 1)))
	}
}

func ConvertBindingErrorToAPIError(rawErr error) *APIError {
	switch err := rawErr.(type) {
	case *APIError:
		return err
	case *ErrorWithScene:
		return toAPIErrorWithScene(err)
	default:
		items := ParseUnprocessableEntityErrorItems(rawErr)
		if items != nil {
			return APIErrorWithScene(ErrUnprocessableEntity, Items(items), Stack(SmallerStacktrace(1, 1)))
		}
		return APIErrorWithScene(ErrBadRequest, Cause(rawErr), Stack(SmallerStacktrace(1, 1)))
	}
}

func ParseUnprocessableEntityErrorItems(rawErr error) []interface{} {
	if rawErr == nil {
		return nil
	}

	// 这里使用As的方式判断，避免因为wrap过后，引起的类型丢失，进而导致提取校验错误信息失败
	var targetErrV vv9.ValidationErrors
	if xerrors.As(rawErr, &targetErrV) {
		return buildValidationErrorsV9Items(targetErrV)
	}

	var targetErrN *strconv.NumError
	if xerrors.As(rawErr, &targetErrN) {
		return []interface{}{
			validationErrorItem{
				Path: "",
				Info: targetErrN.Error(),
			},
		}
	}
	return nil
}

type validationErrorItem struct {
	Path string `json:"path,omitempty"`
	Info string `json:"info"`
}

func buildValidationErrorsV9Items(vv9Errors vv9.ValidationErrors) []interface{} {
	var items []interface{}
	for _, e := range vv9Errors {
		if e != nil {
			path := strToSnake(e.Namespace())
			info := path + " validation failed on the rule: " + e.Tag()
			if e.Param() != "" {
				info = info + "=" + e.Param()
			}
			items = append(items, validationErrorItem{Path: strToSnake(path), Info: info})
		}
	}
	return items
}

func strToSnake(str string) string {
	return strings.Replace(strcase.ToSnake(str), "._", ".", -1)
}

// 作用
// 对于链路是WithScene1 --> Wrap2 --> WithScene3 --> Wrap4 --> Wrap5
// WithScene1 --> Wrap2 --> WithScene3 --> Wrap4 --> Wrap5 --> 基于WithScene3 mainCode和subCode的
// 避免链路丢失，避免底层WithScene后, 上层又不断的wrap导致判断有问题。
// ！！！！由于历史原因这一块目前是无法避免的，而且为了记录错误调用链路也是需要wrap的
func toAPIErrorWithScene(rawErr error) *APIError {
	var ae *APIError
	if xerrors.As(rawErr, &ae) {
		return APIErrorWithScene(NewAPIError(ae.MainCode(), ae.SubCode()), Cause(rawErr), Stack(SmallerStacktrace(2, 1)))
	} else {
		return APIErrorWithScene(ErrInternalError, Cause(rawErr), Stack(SmallerStacktrace(2, 1)))
	}
}
