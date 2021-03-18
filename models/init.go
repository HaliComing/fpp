package models

import (
	"time"

	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB 数据库链接单例
var DB *gorm.DB

// Init 初始化 MySQL 链接
func Init() {
	util.Log().Info("[DB] Init DB.")

	var (
		db  *gorm.DB
		err error
	)

	switch conf.DatabaseConfig.Type {
	case "memory":
		db, err = gorm.Open("sqlite3", ":memory:")
	case "sqlite", "sqlite3":
		db, err = gorm.Open("sqlite3", conf.DatabaseConfig.DBFile)
	default:
		util.Log().Panic("[DB] Database type = %s, Error = Database type is not supported.", conf.DatabaseConfig.Type)
	}

	//db.SetLogger(util.Log())
	if err != nil {
		util.Log().Panic("[DB] Failed to connect to database, Error = %s", err)
	}

	// Debug模式下，输出所有 SQL 日志
	if conf.SystemConfig.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(50)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)

	DB = db

	// 清空数据
	DB.DropTable(&ProxyIP{})
	// 执行迁移
	DB.AutoMigrate(&ProxyIP{})
}
