// Package database 提供了数据库相关的封装，目前实现的是主流的 gorm + mysql
// SaaS 开发者可根据需要替换为其他 orm（如 SQLBoiler，Ent）或者其他数据库（如 mongodb）
// 如果对性能要有很高的话，也可以考虑 sqlx，这是一个高性能的标准 sql 库增强 & 扩展包，
// 其缺点是没有提供完整的 ORM 功能（如自动迁移，关系处理等等），开发者用起来不太方便（需要写不少的 SQL）
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	log "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

var (
	db         *gorm.DB
	dbInitOnce sync.Once
)

const (
	// string 类型字段的默认长度
	defaultStringSize = 256
	// 默认批量创建数量
	defaultBatchSize = 100
	// 默认最大空闲连接
	defaultMaxIdleConns = 20
	// 默认最大连接数
	defaultMaxOpenConns = 100
)

// Client 获取数据库客户端
func Client(ctx context.Context) *gorm.DB {
	if db == nil {
		log.Fatal("database client not init")
	}
	// 设置上下文目的：让 slogGorm 记录日志时带上 Request ID
	return db.WithContext(ctx)
}

// InitDBClient 初始化数据库客户端
func InitDBClient(ctx context.Context, cfg *config.MysqlConfig, slogger *slog.Logger) {
	if db != nil {
		return
	}
	if cfg == nil {
		log.Fatal("mysql config is required when init database client")
	}
	dbInitOnce.Do(func() {
		dbInfo := fmt.Sprintf("mysql %s:%d/%s", cfg.Host, cfg.Port, cfg.Name)

		var err error
		if db, err = newClient(ctx, cfg, slogger); err != nil {
			log.Fatalf("failed to connect database %s: %s", dbInfo, err)
		} else {
			log.Infof(ctx, "database: %s connected", dbInfo)
		}
	})
}

// 初始化 DB Client
func newClient(ctx context.Context, cfg *config.MysqlConfig, slogger *slog.Logger) (*gorm.DB, error) {
	sqlDB, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(defaultMaxIdleConns)
	sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	cCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查 DB 是否可用
	if err = sqlDB.PingContext(cCtx); err != nil {
		return nil, err
	}

	mysqlCfg := mysql.Config{
		DSN:                       cfg.DSN(),
		DefaultStringSize:         defaultStringSize,
		SkipInitializeWithVersion: false,
	}

	gormCfg := &gorm.Config{
		ConnPool: sqlDB,
		// 禁用默认事务（需要手动管理）
		SkipDefaultTransaction: true,
		// 缓存预编译语句
		PrepareStmt: true,
		// Mysql 本身即不支持嵌套事务
		DisableNestedTransaction: true,
		// 批量操作数量
		CreateBatchSize: defaultBatchSize,
		// 数据库迁移时，忽略外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 日志相关
		Logger: slogGorm.New(
			slogGorm.WithTraceAll(),
			slogGorm.WithHandler(slogger.Handler()),
			slogGorm.WithSlowThreshold(200*time.Millisecond),
			slogGorm.WithContextValue(common.RequestIDLogKey, common.RequestIDCtxKey),
			slogGorm.WithContextFunc(common.TraceIDLogKey, log.ExtractTraceID),
			slogGorm.WithContextFunc(common.SpanIDLogKey, log.ExtractSpanID),
		),
	}

	client, err := gorm.Open(mysql.New(mysqlCfg), gormCfg)
	if err != nil {
		return nil, err
	}
	// OpenTelemetry Tracing
	if err = client.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return nil, err
	}

	return client, nil
}

var (
	migSet         *migrationSet
	migSetInitOnce sync.Once
)

// 初始化数据库迁移集
func getMigrationSet() *migrationSet {
	migSetInitOnce.Do(func() {
		migSet = &migrationSet{
			mapping: map[string]*gormigrate.Migration{},
		}
	})
	return migSet
}

// RegisterMigration 注册迁移文件
func RegisterMigration(m *gormigrate.Migration) {
	if err := getMigrationSet().register(m); err != nil {
		log.Fatalf("failed to register migration: %s", err)
	}
}
