package gins

import (
	"context"
	"fmt"
	"strings"

	"github.com/AfterShip/golang-common/errors"
	"github.com/AfterShip/golang-common/logger"
	"github.com/AfterShip/golang-common/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type APIErrorLoggerFunc func(ginCtx *gin.Context, ctx context.Context, apiError *errors.APIError)

var GlobalAPIErrorLoggerFunc APIErrorLoggerFunc = APIErrorLogging

func APIErrorLogging(ginCtx *gin.Context, ctx context.Context, apiError *errors.APIError) {
	//logging
	loggingFields := []zapcore.Field{
		zap.String("category", "http_response_error"),
		zap.String("remote_addr", ginCtx.Request.RemoteAddr),
		zap.String("request_query_string", ginCtx.Request.URL.RawQuery),
	}

	loggingFields = append(loggingFields, logger.APIErrorLoggingFields(apiError)...)

	if ginCtx.Request.Header != nil {
		headerMap := make(map[string]string)
		for key, values := range ginCtx.Request.Header {
			if !security.IsSensitiveHeaderKey(key) {
				headerMap[key] = strings.Join(values, ",")
			}
		}
		loggingFields = append(loggingFields, zap.Any("request_header", headerMap))
	}

	loggingMsg := fmt.Sprintf("%s %s %d %s",
		ginCtx.Request.Method,
		ginCtx.Request.URL.EscapedPath(),
		apiError.MainCode().Code(),
		apiError.MainCode().MessageKey())
	if apiError.MainCode().Code() >= errors.CodeInternalError.Code() {
		logger.Error(
			ctx,
			fmt.Sprintf("[ERROR] %s", loggingMsg),
			loggingFields...,
		)
	} else if apiError.MainCode().Code() != errors.CodeNotFound.Code() {
		//log expect 404 not found
		logger.Warn(
			ctx,
			fmt.Sprintf("[WARNING] %s", loggingMsg),
			loggingFields...,
		)
	}
}
