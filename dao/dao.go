package dao

import (
	"database/sql"
	"fmt"
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/gapiservice/dao/model"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"

	"github.com/google/wire"

	"gorm.io/gorm/schema"


	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Set = wire.NewSet(
	InzoneUserRepoSet,
) // end

var models = []interface{}{
	new(model.InzoneUser),
}

var db *gorm.DB

func Init() error {
	cfg := config.C.Gorm
	var err error
	db, err = newGormDB()
	if err != nil {
		return err
	}

	if cfg.EnableAutoMigrate {
		err = autoMigrate(db)
		if err != nil {
			return err
		}
	}

	return nil
}

// newGormDB 创建DB实例
func newGormDB() (*gorm.DB, error) {
	cfg := config.C
	err := createDatabaseWithMySQL()
	if err != nil {
		return nil, err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	db, err = gorm.Open(mysql.Open(cfg.MySQL.DSN()), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "",
		SingularTable: true,
	}, Logger: newLogger})
	if err != nil {
		return nil, err
	}

	if cfg.RunEnv == config.Develop {
		//db = db.Debug()
	}

	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Gorm.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Gorm.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Gorm.MaxLifetime) * time.Second)

	return db, nil
}

// autoMigrate Auto migration for given models
func autoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	return db.AutoMigrate(
		models...,
	)
}

func createDatabaseWithMySQL() error {
	cfg := config.C.MySQL
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", cfg.DBName)
	_, err = db.Exec(query)
	return err
}
