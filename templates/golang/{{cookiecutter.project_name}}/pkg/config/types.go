package config

import "fmt"

// MysqlConfig Mysql 增强服务配置
type MysqlConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Charset  string
}

// DSN ...
func (cfg *MysqlConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
	)
}

// RabbitMQConfig RabbitMQ 增强服务配置
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Vhost    string
	Password string
}

// DSN ...
func (cfg *RabbitMQConfig) DSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Vhost,
	)
}

// RedisConfig Redis 增强服务配置
type RedisConfig struct {
	Username string
	Host     string
	Port     int
	Password string
	DB       int
}

// DSN ...
func (cfg *RedisConfig) DSN() string {
	return fmt.Sprintf("redis://%s:%s@%s:%d/%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
}

// BkRepoConfig BkRepo 增强服务配置
type BkRepoConfig struct {
	EndpointUrl   string
	Project       string
	Username      string
	Password      string
	Bucket        string
	PublicBucket  string
	PrivateBucket string
}

// BkOtelConfig BkOtel 增强服务配置
type BkOtelConfig struct {
	Trace       bool
	GrpcUrl     string
	BkDataToken string
	Sampler     string
}

// AddonsConfig 增强服务配置
type AddonsConfig struct {
	Mysql    *MysqlConfig
	RabbitMQ *RabbitMQConfig
	Redis    *RedisConfig
	BkRepo   *BkRepoConfig
	BkOtel   *BkOtelConfig
}

// BkPlatUrlConfig 蓝鲸各平台服务地址
type BkPlatUrlConfig struct {
	// 蓝鲸开发者中心地址
	BkPaaS string
	// 统一登录地址
	BkLogin string
	// 组件 API 地址
	BkCompApi string
	// NOTE: SaaS 开发者可按需添加诸如 BkIAM，BkLog 等服务配置
}

// PlatformConfig 平台配置
type PlatformConfig struct {
	// 蓝鲸应用 ID
	AppID string
	// 蓝鲸应用密钥
	AppSecret string
	// 模块名称
	ModuleName string
	// 运行环境：stag 预发布环境，prod 生产环境
	RunEnv string

	// 应用引擎版本
	Region string
	// 推荐的 DB 加密算法有：SHANGMI（对应 SM4CTR 算法）和 CLASSIC（对应 Fernet 算法）
	CryptoType string

	// 蓝鲸根域名，用于获取登录票据，国际化语言等 cookie 信息
	BkDomain string
	// 网关 API 访问地址模板
	ApiUrlTmpl string

	// 蓝鲸平台服务地址配置
	BkPlatUrl BkPlatUrlConfig
	// 增强服务配置
	Addons AddonsConfig
}

// LogConfig 日志配置
type LogConfig struct {
	// 日志级别，可选值为：debug、info、warn、error
	Level string
	// 日志目录，部署于 PaaS 平台上时，该值必须为 /app/v3logs，否则无法采集日志
	Dir string
	// 是否强制标准输出，不输出到文件（一般用于本地开发，标准输出日志查看比较方便）
	ForceToStdout bool
}

// ServerConfig Gin Web Server 配置
type ServerConfig struct {
	// 服务端口
	Port int
	// 优雅退出等待时间
	GraceTimeout int
	// Gin 运行模式
	GinRunMode string
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	// Web Server 配置
	Server ServerConfig
	// 日志配置
	Log LogConfig

	// CORS 允许来源列表
	AllowedOrigins []string
	// AllowedUsers 允许访问的用户列表（UserID）
	AllowedUsers []string
	// CSRF Cookie 域名
	CSRFCookieDomain string
	// 健康探针 Token
	HealthzToken string
	// 指标 API Token
	MetricToken string

	// 是否启用 swagger docs
	EnableSwagger bool
	// 文档文件存放目录
	DocFileBaseDir string
	// 静态文件存放目录
	StaticFileBaseDir string
	// 模板文件存放目录
	TmplFileBaseDir string

	// 缓存内存大小（单位为 MB）
	MemoryCacheSize int

	// DB 加密密钥，若未使用加密功能可不配置
	// 生成方式参见 Readme 文档 - 数据库字段加密
	EncryptSecret string
}

// BizConfig 业务相关配置
type BizConfig struct {
	// NOTE: SaaS 开发者可在此处添加业务相关配置项
}

// Config SaaS 配置
type Config struct {
	// 平台内置配置
	Platform PlatformConfig
	// 服务配置
	Service ServiceConfig
	// 业务配置
	Biz BizConfig
}
