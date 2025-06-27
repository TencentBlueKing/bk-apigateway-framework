// Package migration stores all database migrations
package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
	model2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/model"
)

func init() {
	// Do Not Edit Migration ID!
	migrationID := "20241022_105518"

	database.RegisterMigration(&gormigrate.Migration{
		ID: migrationID,
		Migrate: func(tx *gorm.DB) error {
			logApplying(migrationID)

			return tx.AutoMigrate(&model2.Category{}, &model2.Entry{})
		},
		Rollback: func(tx *gorm.DB) error {
			logRollingBack(migrationID)

			return tx.Migrator().DropTable(&model2.Category{}, &model2.Entry{})
		},
	})
}
