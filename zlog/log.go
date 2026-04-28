package zlog

// import "zinx/ziface"

type Logger struct {
}

type NewLoggerOption func(*Logger) *Logger

func WithLoggerLevel(level int) NewLoggerOption {
	return func(logger *Logger) *Logger {
		// 设置日志级别
		return logger
	}
}

// WithLoggerOutput 设置日志输出目标
func WithLoggerOutput(output string) NewLoggerOption {
	return func(logger *Logger) *Logger {
		// 设置日志输出目标
		return logger
	}
}

// WithLoggerFormat 设置日志格式
func WithLoggerFormat(format string) NewLoggerOption {
	return func(logger *Logger) *Logger {
		// 设置日志格式
		return logger
	}
}

// NewLogger 创建一个新的 Logger 实例，并应用提供的选项
func NewLogger(opts ...NewLoggerOption) *Logger {
	logger := &Logger{}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}
