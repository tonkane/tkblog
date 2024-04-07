package log

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tkane/tkblog/internal/pkg/known"
)

// 定义日志接口
type Logger interface {
	Debugw(msg string, kv ...interface{})
	Infow(msg string, kv ...interface{})
	Warnw(msg string, kv ...interface{})
	Errorw(msg string, kv ...interface{})
	Panicw(msg string, kv ...interface{})
	Fatalw(msg string, kv ...interface{})
	Sync()
}


// logger 
type zapLogger struct {
	z *zap.Logger
}

// 确保接口实现
var _ Logger = &zapLogger{}

var (
	mu sync.Mutex
	std = NewLogger(NewOptions())
)

// 初始化
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = NewLogger(opts)
}

func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 将 opts 里面文本类型的日志级别转换为 zapcore.Level 类型
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 如果日志级别非法则指定 info 为默认
		zapLevel = zapcore.InfoLevel
	}

	// 创建 encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义key名称
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	// 时间序列化指定格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// duration 序列化函数, 精确到毫秒
	encoderConfig.EncodeDuration = func (d time.Duration, enc zapcore.PrimitiveArrayEncoder)  {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	// 设置 zap logger 配置
	cfg := &zap.Config{
		DisableCaller: opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Level: zap.NewAtomicLevelAt(zapLevel),
		Encoding: opts.Format,
		EncoderConfig: encoderConfig,
		OutputPaths: opts.OutputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	// 创建 zap logger 对象
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger := &zapLogger{z: z}

	// 将标准库 log.Logger 的info级别输出重定向到 zap.Logger 中
	zap.RedirectStdLog(z)
	return logger

}

// 调用zap sunc 将缓存中的日志刷新到磁盘上
func Sync() {std.z.Sync()}

func (l *zapLogger) Sync() {
	l.z.Sync()
}

// debug
func Debugw(msg string, kv ...interface{}) {
	std.z.Sugar().Debugw(msg, kv...)
}

func (l *zapLogger) Debugw(msg string, kv ...interface{}) {
	l.z.Sugar().Debugw(msg, kv...)
}

// info
func Infow(msg string, kv ...interface{}) {
	std.z.Sugar().Infow(msg, kv...)
}

func (l *zapLogger) Infow(msg string, kv ...interface{}) {
	l.z.Sugar().Infow(msg, kv...)
}

// warn
func Warnw(msg string, kv ...interface{}) {
	std.z.Sugar().Warnw(msg, kv...)
}

func (l *zapLogger) Warnw(msg string, kv ...interface{}) {
	l.z.Sugar().Warnw(msg, kv...)
}

// error
func Errorw(msg string, kv ...interface{}) {
	std.z.Sugar().Errorw(msg, kv...)
}

func (l *zapLogger) Errorw(msg string, kv ...interface{}) {
	l.z.Sugar().Errorw(msg, kv...)
}

// panic
func Panicw(msg string, kv ...interface{}) {
	std.z.Sugar().Panicw(msg, kv...)
}

func (l *zapLogger) Panicw(msg string, kv ...interface{}) {
	l.z.Sugar().Panicw(msg, kv...)
}

// fatal
func Fatalw(msg string, kv ...interface{}) {
	std.z.Sugar().Fatalw(msg, kv...)
}

func (l *zapLogger) Fatalw(msg string, kv ...interface{}) {
	l.z.Sugar().Fatalw(msg, kv...)
}

// 解析 context
func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}

func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()

	if requestID := ctx.Value(known.XRequestIDKey); requestID != nil {
		lc.z = lc.z.With(zap.Any(known.XRequestIDKey, requestID))
	}

	return lc
}

// clone 防止并发时候 id 串
func (l *zapLogger) clone() *zapLogger {
	lc := *l
	return &lc
}