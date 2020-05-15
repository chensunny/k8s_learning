package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AfterShip/golang-common/errors"
	"github.com/AfterShip/golang-common/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger
var globalConf LoggerConf
var _beforeLogHook BeforeLogHook

const (
	LogLevelDebug    = "DEBUG"
	LogLevelInfo     = "INFO"
	LogLevelWarn     = "WARNING"
	LogLevelError    = "ERROR"
	LogLevelCritical = "CRITICAL"
)

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  LogLevelDebug,
	zapcore.InfoLevel:   LogLevelInfo,
	zapcore.WarnLevel:   LogLevelWarn,
	zapcore.ErrorLevel:  LogLevelError,
	zapcore.DPanicLevel: LogLevelCritical,
	zapcore.PanicLevel:  LogLevelCritical,
	zapcore.FatalLevel:  LogLevelCritical,
}

func init() {
	zapLogger, _ := BuildZapDriverLogger(LoggerConf{
		Level:    "info",
		Encoding: "json",
	})
	globalLogger = zapLogger
	_beforeLogHook = DefaultBeforeLogHookImpl{}
}

type LoggerConf struct {
	Level    string
	Encoding string
}

// 打印log之前的hook操作。
// 典型场景：用于从ctx中提取值给每条logger添加对应的值
type BeforeLogHook interface {
	BeforeLog(ctx context.Context, msg string, fields []zap.Field) (context.Context, string, []zap.Field)
}

type DefaultBeforeLogHookImpl struct{}

// 默认实现，保持和原来一致
func (d DefaultBeforeLogHookImpl) BeforeLog(ctx context.Context, msg string, fields []zap.Field) (context.Context, string, []zap.Field) {
	if ctx == nil {
		return ctx, msg, fields
	}

	if ctx.Err() != nil {
		fields = append(fields, zap.NamedError("context_error", ctx.Err()))
	}
	cloudflareRay := tracing.GetCloudflareRayFromContext(ctx)
	if cloudflareRay != "" {
		fields = append(fields, zap.String("context_cloudflare_ray", cloudflareRay))
	}

	traceID := tracing.GetTraceIDFromContext(ctx)
	if traceID != "" {
		fields = append(fields, zap.String("context_trace_id", traceID))
	}
	if deadline, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.Time("context_deadline", deadline))
	}
	requestMethod := tracing.GetRequestMethodFromContext(ctx)
	if requestMethod != "" {
		fields = append(fields, zap.String("context_request_method", requestMethod))
	}
	requestPath := tracing.GetRequestPathFromContext(ctx)
	if requestPath != "" {
		fields = append(fields, zap.String("context_request_path", requestPath))
	}
	return ctx, msg, fields
}

func EncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[lv])
}

func BuildZapLogger(conf LoggerConf) (*zap.Logger, error) {
	return BuildZapDriverLogger(conf)
}

//adapt stackdriver format
func BuildZapDriverLogger(conf LoggerConf) (*zap.Logger, error) {
	globalConf = conf
	//zap log level
	zapLevel := zap.NewAtomicLevel()
	err := zapLevel.UnmarshalText([]byte(conf.Level))
	if err != nil {
		fmt.Println("unmarshal zap log level failed:", err)
	}
	if conf.Encoding == "" {
		conf.Encoding = "json"
	}

	//encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "eventTime",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    EncodeLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConf := zap.Config{
		Level:            zapLevel,
		Encoding:         conf.Encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	//our errors will always be wrapped that contains stackstrace
	//so add stacktrace only when level >= DPanic
	return zapConf.Build(zap.AddStacktrace(zapcore.DPanicLevel), zap.AddCallerSkip(1))
}

func SetZapLogger(logger *zap.Logger) {
	globalLogger = logger
}

func GetZapLogger() *zap.Logger {
	if globalLogger == nil {
		return zap.L()
	}
	return globalLogger
}

// 设置包内变量：_beforeLogHook，提供调用方自定义BeforeLogHook接口的能力
func SetBeforeLogHook(beforeLogHook BeforeLogHook) {
	_beforeLogHook = beforeLogHook
}

// 统一规范 error level从INFO开始
// 如果有需要输出debug的地方需要判断IsDebugEnabled，节省不必要的内存资源
func IsDebugEnabled() bool {
	return strings.TrimSpace(strings.ToLower(globalConf.Level)) == "debug"
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
	GetZapLogger().Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
	GetZapLogger().Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
	GetZapLogger().Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
	GetZapLogger().Error(msg, fields...)
}

func Critical(ctx context.Context, msg string, fields ...zap.Field) {
	ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
	GetZapLogger().DPanic(msg, fields...)
}

// 可自定义打印的log等级
// 场景：例如底层公共库，允许使用方自定义打印的log等级
func Print(ctx context.Context, logLevel string, msg string, fields ...zap.Field) {
	switch logLevel {
	case LogLevelDebug:
		ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
		GetZapLogger().Debug(msg, fields...)
	case LogLevelInfo:
		ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
		GetZapLogger().Info(msg, fields...)
	case LogLevelWarn:
		ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
		GetZapLogger().Warn(msg, fields...)
	case LogLevelError:
		ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
		GetZapLogger().Error(msg, fields...)
	case LogLevelCritical:
		ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
		GetZapLogger().DPanic(msg, fields...)
	default:
		ctx, msg, fields = _beforeLogHook.BeforeLog(ctx, msg, fields)
		GetZapLogger().Debug(msg, fields...)
	}
}

func ErrorField(err error) zap.Field {
	if wrappedErr, ok := err.(*errors.APIError); ok {
		return APIErrorLoggingFields(wrappedErr)[0]
	}
	if wrappedErr, ok := err.(*errors.ErrorWithScene); ok {
		return ErrorWithSceneLoggingFields(wrappedErr)[0]
	}
	return zap.Error(err)
}

func APIErrorLoggingFields(err *errors.APIError) []zap.Field {
	return []zap.Field{zap.Any("error", json.RawMessage(err.JsonString()))}
}

func ErrorWithSceneLoggingFields(err *errors.ErrorWithScene) []zap.Field {
	return []zap.Field{zap.Any("error", json.RawMessage(err.JsonString()))}
}
