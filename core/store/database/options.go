package database

import (
	"gorm.io/gorm/logger"
	"time"
)

// Options defines optsions for mysql database.
type Options struct {
	DBType                string        `json:"dbtype,omitempty"                   mapstructure:"db-type"`
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"password"                           mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connecti"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connecti"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-li"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
	IsDebug               bool          `json:"isdebug"                           mapstructure:"is-debug"`
	TablePrefix           string        `json:"tableprefix"                       mapstructure:"table-prefix"`
	Logger                logger.Interface
}
