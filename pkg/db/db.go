package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLOptions struct {
	Host string
	Username string
	Password string
	Database string
	MaxIdleConnections int
	MaxOpenConnections int
	MaxConnectionLifeTime time.Duration
	LogLevel int
}

func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
			o.Username,
			o.Password,
			o.Host,
			o.Database,
			true,
			"Local")
}


// 创建 gorm 数据库实例
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	logLevel := logger.Silent
	if opts.LogLevel != 0 {
		logLevel = logger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 最大连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// 最长可重用时间
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// 空闲连接池最大连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}