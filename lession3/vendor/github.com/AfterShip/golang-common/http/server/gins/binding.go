package gins

import (
	"sync"

	"github.com/AfterShip/golang-common/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var once = new(sync.Once)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// 1. types.WrappedType本身实现json的Unmarshal接口，所以使用gin.Context本身的ShouldBindJSON是支持types.WrappedType的binding
// 2. 但是这里的ShouldBindJSON等ShouldBind系列方法，用于替代gin.Context的ShouldBind，内部使用jsoniter来做json binding，提供了binding的一些扩展点支持
// 3. 提供了jsoniter的autoSetNullExtension扩展
//    用于满足API request body 里多级嵌套的时候，某个父级字段为null的时候代表所有子字段应该为null
//    如果父字段用pointer表示null，那么业务层代码想要表达所有嵌套子字段为null的话会非常难写难受
//    所以合适的方式应该是用struct value，然后用WrappedType来表达null(基本类型没有nil的支持)
//    典型场景 REST PATCH API
// 4. 提供了before 和 after decoding的hook，hook返回err的时候则decode失败
//    TODO hook的使用场景
// 5. types.WrappedType本身实现json的Unmarshal接口，所以使用gin.Context本身的ShouldBindJSON是支持types.WrappedType的binding
func ShouldBindJSON(ctx *gin.Context, ptr interface{}) error {
	once.Do(func() {
		jsoniter.RegisterExtension(&autoSetNullExtension{})
	})
	if ctx.Request == nil || ctx.Request.Body == nil {
		return errors.ErrUnprocessableEntity
	}
	decoder := json.NewDecoder(ctx.Request.Body)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if beforeDecodingHook, ok := ptr.(BeforeDecodingHook); ok {
		if err := beforeDecodingHook.BeforeDecode(ctx); err != nil {
			return err
		}
	}
	if err := decoder.Decode(ptr); err != nil {
		return err
	}
	if afterDecodingHook, ok := ptr.(AfterDecodingHook); ok {
		if err := afterDecodingHook.AfterDecode(ctx); err != nil {
			return err
		}
	}

	if binding.Validator == nil {
		return nil
	}
	if beforeValidatingHook, ok := ptr.(BeforeValidatingHook); ok {
		if err := beforeValidatingHook.BeforeValidate(ctx, binding.Validator); err != nil {
			return err
		}
	}
	if err := binding.Validator.ValidateStruct(ptr); err != nil {
		return err
	}
	if afterValidatingHook, ok := ptr.(AfterValidatingHook); ok {
		if err := afterValidatingHook.AfterValidate(ctx, binding.Validator); err != nil {
			return err
		}
	}
	return nil
}

// TODO: 这里缺了对WrappedType的支持
func ShouldBindQuery(ctx *gin.Context, ptr interface{}) error {
	return ctx.ShouldBindQuery(ptr)
}

func ShouldBindXml(ctx *gin.Context, ptr interface{}) error {
	return ctx.ShouldBindXML(ptr)
}

func ShouldBind(ctx *gin.Context, ptr interface{}) error {
	return ctx.ShouldBind(ptr)
}

func ShouldBindWith(ctx *gin.Context, ptr interface{}, b binding.Binding) error {
	if b == binding.JSON {
		return ShouldBindJSON(ctx, ptr)
	}
	return b.Bind(ctx.Request, ptr)
}
