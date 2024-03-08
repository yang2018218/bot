package postgres

import (
	"os"
	"path/filepath"
	"time"
	"wechatbot/internal/pkg/code"
	"wechatbot/internal/pkg/utils"
	"wechatbot/pkg/errors"

	"gorm.io/gorm"
)

type FlywayModel struct {
	ID        string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:ctime"`
	Name      string    `gorm:"column:name;unique_index"`
	Sql       string    `gorm:"column:content"`
}

func (FlywayModel) TableName() string {
	return "flyway"
}

func Migrate(db *gorm.DB, flywayDir string) (flyways []FlywayModel, err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = errors.WithCode(code.ErrInternalServerError, "内部服务出错")
		}
		if err != nil {
			flyways = []FlywayModel{}
			tx.Rollback()
		}
	}()

	tx.AutoMigrate(&FlywayModel{})
	files, _ := os.ReadDir(flywayDir)
	fileName := []string{}
	for _, file := range files {
		fileName = append(fileName, file.Name())
	}
	skips := []FlywayModel{}
	tx.Where("name IN ?", fileName).Find(&skips)
	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		if fileExt != ".sql" {
			continue
		}
		var flyway FlywayModel
		for _, skip := range skips {
			if skip.Name == file.Name() {
				flyway = skip
				break
			}
		}
		if flyway.ID != "" {
			continue
		}
		fileByte, _ := os.ReadFile(filepath.Join(flywayDir, file.Name()))
		fileStr := string(fileByte)
		err = tx.Exec(fileStr).Error
		if err != nil {
			return
		}
		// 保存执行结果
		flyway = FlywayModel{
			ID:   utils.NewObjectID(),
			Name: file.Name(),
			Sql:  fileStr,
		}
		err = tx.Create(&flyway).Error
		if err != nil {
			return
		}
		flyways = append(flyways, flyway)
	}
	if err == nil {
		err = tx.Commit().Error
	}
	return
}
