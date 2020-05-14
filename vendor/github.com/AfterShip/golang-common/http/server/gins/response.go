package gins

import (
	"context"
	"net/http"

	"github.com/AfterShip/golang-common/errors"
	"github.com/AfterShip/golang-common/http/model"
	"github.com/AfterShip/golang-common/tracing"
	"github.com/gin-gonic/gin"
)

// see https://docs.google.com/spreadsheets/d/16oduZidE9ofdoT6m3I9oCNa9hPjgSCYmSeg3jxjofdo/edit#gid=524712657
var (
	emptyData = struct {
	}{}
)

func ResponseOK(ginCtx *gin.Context, data interface{}) {
	ResponseSuccess(ginCtx, http.StatusOK, data)
}

func ResponseCreated(ginCtx *gin.Context, data interface{}) {
	ResponseSuccess(ginCtx, http.StatusCreated, data)
}

func ResponseAccepted(ginCtx *gin.Context, data interface{}) {
	ResponseSuccess(ginCtx, http.StatusAccepted, data)
}

func ResponseSuccess(ginCtx *gin.Context, httpStatus int, data interface{}) {
	ginCtx.JSON(httpStatus, model.ResponseBody{
		Meta: model.BuildResponseMeta(model.BuildMetaCode(httpStatus, 0), nil),
		Data: data,
	})
}

//400 or 422
func ResponseInputBindingError(ginCtx *gin.Context, ctx context.Context, err error) {
	items := errors.ParseUnprocessableEntityErrorItems(err)
	if len(items) > 0 {
		ResponseAPIErrorWithLogging(ginCtx, ctx, errors.WithScene(errors.ErrUnprocessableEntity, errors.Items(items)))
	} else {
		// body io EOF、ErrUnexpectedEOF、entity 格式语法不正确等
		ResponseAPIErrorWithLogging(ginCtx, ctx, errors.WithScene(errors.ErrBadRequest, errors.Cause(err)))
	}
}

func ResponseError(ginCtx *gin.Context, ctx context.Context, err error) {
	var apiError = errors.ConvertToAPIError(err)
	ResponseAPIErrorWithLogging(ginCtx, ctx, apiError)
}

func ResponseAPIErrorWithLogging(ginCtx *gin.Context, ctx context.Context, err *errors.APIError) {
	ResponseAPIError(ginCtx, ctx, err)
	if GlobalAPIErrorLoggerFunc != nil {
		GlobalAPIErrorLoggerFunc(ginCtx, ctx, err)
	}
}

func ResponseAPIError(ginCtx *gin.Context, ctx context.Context, err *errors.APIError) {
	code := model.BuildMetaCode(err.MainCode().Code(), err.SubCode().Code())

	meta := model.BuildResponseMeta(code, err)
	meta.Errors = err.Scene().Items()

	resp := &model.ResponseBody{
		Meta: meta,
		Data: emptyData,
	}

	ginCtx.AbortWithStatusJSON(err.MainCode().Code(), resp)
	//add process error to request context for unified middleware process, eg. newrelic tracing middleware
	ginCtx.Request = ginCtx.Request.WithContext(tracing.ContextWithTaskProcessError(ginCtx.Request.Context(), err))
}
