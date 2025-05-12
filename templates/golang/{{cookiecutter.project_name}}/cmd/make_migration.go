package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/database"
	log "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

var migrationTmpl = `
// Package migration stores all database migrations
package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
)


func init() {
	// Do Not Edit Migration ID!
	migrationID := "[[ .id ]]"

	database.RegisterMigration(&gormigrate.Migration{
		ID: migrationID,
		Migrate: func(tx *gorm.DB) error {
			logApplying(migrationID)

			// TODO implement migrate code
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			logRollingBack(migrationID)

			// TODO implement rollback code
			return nil
		},
	})
}
`

var makeMigrationCmd = &cobra.Command{
	Use:   "make-migration",
	Short: "Generate an empty migration file.",
	Run: func(cmd *cobra.Command, args []string) {
		migrationID := database.GenMigrationID()

		// 文件
		fileName := fmt.Sprintf("%s.go", migrationID)
		filePath := path.Join(config.BaseDir, "pkg/migration", fileName)
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed to create migration file with path: %s, err: %s", filePath, err)
		}
		defer file.Close()

		// 模板
		// Q：为什么需要修改 Delims 为方括号
		// A：避免在开发者中心渲染模板时候渲染此处
		tmpl, err := template.New("migration").
			Delims("[[", "]]").
			Parse(strings.TrimLeft(migrationTmpl, "\n"))
		if err != nil {
			log.Fatal("failed to initialize migration template")
		}
		if err = tmpl.Execute(file, map[string]string{"id": migrationID}); err != nil {
			log.Fatal("failed to render migration file from template")
		}

		log.Infof(
			context.Background(),
			"migration file %s generated, you must edit it and "+
					"implement the migration logic and then run `migrate` to apply",
			fileName,
		)
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
}
