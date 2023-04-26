package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/env"
)

var (
	logger *zap.Logger
)

func Init(conf *config.LoggerConfig) {
	file := os.Stdout
	if conf.Filename != "/dev/stdout" {
		var err error
		file, err = os.Open(conf.Filename)
		if err != nil {
			panic(fmt.Sprintf("error open file %v", err))
		}
	}
	ecf := zap.NewProductionEncoderConfig()
	ecf.FunctionKey = "func"
	ecf.EncodeTime = zapcore.ISO8601TimeEncoder
	ecf.ConsoleSeparator = " "
	ecf.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		EncodeWrapper{zapcore.NewConsoleEncoder(ecf)},
		&zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(file),
			Size:          0,
			FlushInterval: time.Second * 1,
			Clock:         nil,
		},
		zap.InfoLevel,
	)
	logger = zap.New(core, zap.AddCallerSkip(1), zap.AddCaller())
}

func NewLogger() *zap.Logger {
	return logger.WithOptions(zap.AddCallerSkip(-1))
}

func getLogger() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field)  { getLogger().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)   { getLogger().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)   { getLogger().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field)  { getLogger().Error(msg, fields...) }
func DPanic(msg string, fields ...zap.Field) { getLogger().DPanic(msg, fields...) }
func Panic(msg string, fields ...zap.Field)  { getLogger().Panic(msg, fields...) }
func Fatal(msg string, fields ...zap.Field)  { getLogger().Fatal(msg, fields...) }
func Sync()                                  { getLogger().Sync(); logger.Sync() }

func S() *zap.SugaredLogger { return getLogger().Sugar() }

func Debugw(msg string, keysAndValues ...interface{}) { S().Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...interface{})  { S().Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...interface{})  { S().Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...interface{}) { S().Errorw(msg, keysAndValues...) }
func Panicw(msg string, keysAndValues ...interface{}) { S().Panicw(msg, keysAndValues...) }
func Fatalw(msg string, keysAndValues ...interface{}) { S().Fatalw(msg, keysAndValues...) }

func Debugf(msg string, keysAndValues ...interface{}) { S().Debugf(msg, keysAndValues...) }
func Infof(msg string, keysAndValues ...interface{})  { S().Infof(msg, keysAndValues...) }
func Warnf(msg string, keysAndValues ...interface{})  { S().Warnf(msg, keysAndValues...) }
func Errorf(msg string, keysAndValues ...interface{}) { S().Errorf(msg, keysAndValues...) }
func Panicf(msg string, keysAndValues ...interface{}) {
	if env.IsLocal() {
		S().Errorf(msg, keysAndValues...)
	} else {
		S().Panicf(msg, keysAndValues...)
	}
}
func Fatalf(msg string, keysAndValues ...interface{}) { S().Fatalf(msg, keysAndValues...) }

type EncodeWrapper struct {
	zapcore.Encoder
}

func (e EncodeWrapper) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	function := ent.Caller.Function
	count := 0
	index := strings.IndexFunc(function, func(r rune) bool {
		if r == '/' {
			count++
		}
		if count >= 3 {
			return true
		}
		return false
	})
	function = function[index+1:]

	ent.Caller.Function = function
	return e.Encoder.EncodeEntry(ent, fields)
}
