package logger

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLevel      = zapcore.InfoLevel
	DefaultTimeLayout = time.RFC3339
)

type Option func(*option)

type option struct {
	level          zapcore.Level
	timeLayout     string
	file           io.Writer
	fields         *sync.Map
	disableConsole bool
}

func WithLevel(level zapcore.Level) Option {
	return func(o *option) {
		o.level = level
	}
}

func WithField(key string, value any) Option {
	return func(o *option) {
		if o.fields == nil {
			o.fields = &sync.Map{}
		}
		o.fields.Store(key, value)
	}
}

func WithTimeLayout(layout string) Option {
	return func(o *option) {
		o.timeLayout = layout
	}
}

func WithFilePath(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return func(o *option) {
		o.file = zapcore.Lock(f)
	}
}

func WithDriver(driver io.Writer) Option {
	return func(o *option) {
		o.file = driver
	}
}

func WithFileRotation(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	return func(o *option) {
		o.file = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    128,
			MaxBackups: 300,
			MaxAge:     30,
			Compress:   true,
			LocalTime:  true,
		}
	}
}

func WithDisableConsole() Option {
	return func(o *option) {
		o.disableConsole = false
	}
}

func NewJsonLogger(opts ...Option) (*zap.Logger, error) {
	opt := &option{
		level:      DefaultLevel,
		timeLayout: DefaultTimeLayout,
		fields:     &sync.Map{},
	}

	for _, f := range opts {
		f(opt)
	}

	// similar to zap.NewProductionEncoderConfig()
	encoderConfigure := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(opt.timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfigure)

	// LowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfigure),
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfigure),
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	logger := zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	opt.fields.Range(func(key, value any) bool {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{Key: key.(string), Type: zapcore.StringType, String: value.(string)}))
		return true
	})

	return logger, nil
}
