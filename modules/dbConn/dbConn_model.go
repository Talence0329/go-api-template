package dbconn

import (
	"database/sql"
	"sync"
	"time"
)

type DBDriver string

const DBDriverMySQL DBDriver = "mysql"
const DBDriverPostgres DBDriver = "postgres"

type Config struct {
	DefaultConnConfig DefaultConnConfig
}

type DBName string
type Envkey string

type DBConfig struct {
	db         *sql.DB
	DBDriver   DBDriver
	DSNSource  any
	ConnConfig ConnConfig
	Once       sync.Once
}

type ConnConfig struct {
	ConnMaxIdleTime *time.Duration
	ConnMaxLifetime *time.Duration
	MaxOpenConns    *int
	MaxIdleConns    *int
}

type DefaultConnConfig struct {
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
