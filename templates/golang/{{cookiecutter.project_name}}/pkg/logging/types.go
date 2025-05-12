package logging

// Gin, Gorm 默认日志等级为 warn，目的是避免记录过多无关日志
// 开发者可根据需求自行调整，可选值：debug、info、warn、error
const (
	// GinLogLevel gin 日志级别
	GinLogLevel = "warn"
	// GormLogLevel gorm 日志级别
	GormLogLevel = "warn"
)

// Options Logger 配置
type Options struct {
	// 日志级别
	Level string
	// 日志内容 Handler，支持 text 和 json
	HandlerName string
	// io.Writer, 支持 stdout、stderr、file
	WriterName string
	// Writer 配置
	WriterConfig map[string]string
}
