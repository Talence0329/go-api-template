package database

import (
	dbconn "backend/modules/dbConn"
	"time"

	sq "github.com/Masterminds/squirrel"
)

// 控管所有連線，其他模組僅組裝連線字串與進行連線等作業

type Config struct {
	DefaultPostgresConnConfig dbconn.ConnConfig
}

// todo : 改成可以透過init修改DefaultConnConfig

func Init() {
	connMaxIdleTime := time.Hour * 1
	connMaxLifetime := time.Hour * 1
	maxOpenConns := 10
	maxIdleConns := 5
	dbconn.Init(dbconn.Config{
		DefaultConnConfig: dbconn.DefaultConnConfig{
			ConnMaxIdleTime: connMaxIdleTime,
			ConnMaxLifetime: connMaxLifetime,
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
		},
	})

	dbconn.New(map[dbconn.DBName]dbconn.DBConfig{
		MEMBER: {
			DBDriver:  dbconn.DBDriverPostgres,
			DSNSource: dbconn.Envkey("MEMBER_POSTGRES_URL"),
		},
		BACKSTAGE: {
			DBDriver:  dbconn.DBDriverPostgres,
			DSNSource: EnvkeyBackstage,
		},
	})
}

func PSSQ() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
