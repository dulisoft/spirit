package database

import (
	"gorm.io/gorm/logger"
	"time"
)

// Options defines optsions for mysql database.
type Options struct {
	DBType                string
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
	IsDebug               bool
	TablePrefix           string
}
