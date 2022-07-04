package zaplogger

import (
	"os"
	"runtime"
	"time"

	"github.com/bluele/zapslack"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	StdFormatLog      = `%s %s %s %s "{"request":%s, "response":%s}"`
	StdFormatErrorLog = `%s %s %s %s "{"request":%s, "response":%s}" %s`

	Topic     = "topic"
	Partition = "partition"
	Message   = "message"
	WorkerID  = "workerID"
	Offset    = "offset"
	Time      = "time"

	GRPC     = "GRPC"
	SIZE     = "SIZE"
	URI      = "URI"
	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"
)

type (
	ListErrors struct {
		Error    string
		File     string
		Function string
		Line     int
		Extra    interface{}
	}
	Fields map[string]interface{}
)

//Logger is our contract for the logger
type Logger interface {
	SetMessageLog(err error, depthList ...int) *ListErrors

	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	WarnMsg(msg string, err error)

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Fatal(args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger

	WithName(name string)

	Sync() error

	Desugar() *zap.Logger

	KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time)

	KafkaLogCommittedMessage(topic string, partition int, offset int64)

	GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error)

	GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error)
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewZapLogger(logPath, slackWebHookUrl string) Logger {

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     syslogTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, fileSyncer, highPriority),
	)

	logger := zap.New(
		core,
		zap.AddCallerSkip(1),
		zap.AddCaller()).Sugar()

	if slackWebHookUrl != "" {
		logger = zap.New(
			core,
			zap.AddCallerSkip(1),
			zap.AddCaller(),
			zap.Hooks(zapslack.NewSlackHook(slackWebHookUrl, zap.ErrorLevel).GetHook())).Sugar()
	}

	return &zapLogger{sugaredLogger: logger}
}

func (s zapLogger) SetMessageLog(err error, depthList ...int) *ListErrors {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	le := new(ListErrors)
	if function, file, line, ok := runtime.Caller(depth); ok {
		le.Error = err.Error()
		le.File = file
		le.Function = runtime.FuncForPC(function).Name()
		le.Line = line
	} else {
		le = nil
	}
	return le
}

func (l *zapLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	l.sugaredLogger.Info(
		GRPC,
		zap.String(METHOD, method),
		zap.Duration(TIME, time),
		zap.Any(METADATA, metaData),
		zap.Error(err),
	)
}

func (l *zapLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.sugaredLogger.Info(
		GRPC,
		zap.String(METHOD, method),
		zap.Any(REQUEST, req),
		zap.Any(REPLY, reply),
		zap.Duration(TIME, time),
		zap.Any(METADATA, metaData),
		zap.Error(err),
	)
}

func (l *zapLogger) KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time) {
	l.sugaredLogger.Debug(
		"Processing Kafka message",
		zap.String(Topic, topic),
		zap.Int(Partition, partition),
		zap.String(Message, message),
		zap.Int(WorkerID, workerID),
		zap.Int64(Offset, offset),
		zap.Time(Time, time),
	)
}

func (l *zapLogger) KafkaLogCommittedMessage(topic string, partition int, offset int64) {
	l.sugaredLogger.Info(
		"Committed Kafka message",
		zap.String(Topic, topic),
		zap.Int(Partition, partition),
		zap.Int64(Offset, offset),
	)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

// WithName add logger microservice name
func (l *zapLogger) WithName(name string) {
	l.sugaredLogger = l.sugaredLogger.Named(name)
}

// WarnMsg log error message with warn level.
func (l *zapLogger) WarnMsg(msg string, err error) {
	l.sugaredLogger.Warn(msg, zap.String("error", err.Error()))
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *zapLogger) Sync() error {
	return l.sugaredLogger.Sync()
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

func (s zapLogger) Desugar() *zap.Logger {
	return s.sugaredLogger.Desugar()
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(helper.DateTimeFormatDefault))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
