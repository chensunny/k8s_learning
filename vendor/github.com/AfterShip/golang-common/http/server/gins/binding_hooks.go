package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BeforeDecodingHook interface {
	BeforeDecode(ctx *gin.Context) error
}

type AfterDecodingHook interface {
	AfterDecode(ctx *gin.Context) error
}

type BeforeValidatingHook interface {
	BeforeValidate(ctx *gin.Context, validator binding.StructValidator) error
}

type AfterValidatingHook interface {
	AfterValidate(ctx *gin.Context, validator binding.StructValidator) error
}
