// Package migration stores all database migrations
package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/database"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/model"
)

func init() {
	// Do Not Edit Migration ID!
	migrationID := "20241022_110326"

	database.RegisterMigration(&gormigrate.Migration{
		ID: migrationID,
		Migrate: func(tx *gorm.DB) error {
			logApplying(migrationID)

			return tx.AutoMigrate(&model.Task{}, &model.PeriodicTask{})
		},
		Rollback: func(tx *gorm.DB) error {
			logRollingBack(migrationID)

			return tx.Migrator().DropTable(&model.Task{}, &model.PeriodicTask{})
		},
	})
}
