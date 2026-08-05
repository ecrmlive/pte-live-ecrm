package db

import (
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "time/tzdata" // alpine 无系统 zoneinfo 时仍可 LoadLocation Asia/Shanghai
)

const shanghaiLocationName = "Asia/Shanghai"

// EnsureShanghaiTimezone 将进程 time.Local 固定为 Asia/Shanghai。
// 管理后台展示与业务写入统一按上海时区理解「本地时间」。
func EnsureShanghaiTimezone() {
	loc, err := time.LoadLocation(shanghaiLocationName)
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	time.Local = loc
}

// NormalizeShanghaiDSN 强制 parseTime、loc=Asia/Shanghai、会话 time_zone=+08:00。
func NormalizeShanghaiDSN(dsn string) (string, error) {
	if dsn == "" {
		return "", fmt.Errorf("mysql dsn empty")
	}
	cfg, err := mysqldriver.ParseDSN(dsn)
	if err != nil {
		return "", fmt.Errorf("parse mysql dsn: %w", err)
	}
	loc, err := time.LoadLocation(shanghaiLocationName)
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	cfg.ParseTime = true
	cfg.Loc = loc
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	cfg.Params["time_zone"] = "'+08:00'"
	return cfg.FormatDSN(), nil
}

func OpenMySQL(dsn string, debug bool) (*gorm.DB, error) {
	EnsureShanghaiTimezone()
	normalized, err := NormalizeShanghaiDSN(dsn)
	if err != nil {
		return nil, err
	}
	logLevel := logger.Warn
	if debug {
		logLevel = logger.Info
	}
	gdb, err := gorm.Open(mysql.Open(normalized), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return gdb, nil
}
