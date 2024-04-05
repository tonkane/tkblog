package log

import(
	"go.uber.org/zap/zapcore"
)

type Options struct {
	// 是否打印日志所在的文件和行号
	DisableCaller bool
	// 是否禁止 panic 及以上级别打印堆栈信息
	DisableStacktrace bool
	// 日志级别 debug, info , warn, error, panic, fatal
	Level string
	// 日志格式 console, json
	Format string
	// 日志输出位置
	OutputPaths []string
}

func NewOptions() *Options {
	return &Options {
		false, false, zapcore.InfoLevel.String(), "console", []string{"stdout"},
	}
}