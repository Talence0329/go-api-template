package dbconn

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var cfg = Config{}

var isInited = false

var dbList = make(map[DBName]DBConfig)

func Init(initConfig Config) {
	cfg = initConfig
	isInited = true
}

func New(cList map[DBName]DBConfig) {
	if !isInited {
		log.Fatalf("[func->New] dbConn 尚未進行過Init")
	}

	for k, c := range cList {
		dsn := ""
		db := &sql.DB{}
		c.Once.Do(func() {
			// 取得DSN
			if _dsn, _err := getDSN(c.DSNSource); _err != nil {
				log.Fatalf("%v", _err)
			} else {
				dsn = _dsn
			}
			if _db, _err := sql.Open(string(c.DBDriver), dsn); _err != nil {
				log.Fatalf("[func->New] %s 連線失敗 %s", dsn, _err)
			} else {
				if _err := _db.Ping(); _err != nil {
					log.Fatal(_err.Error())
					return
				}
				db = _db
			}

			// 設定該連線的參數
			if c.ConnConfig.MaxOpenConns == nil {
				db.SetMaxOpenConns(cfg.DefaultConnConfig.MaxOpenConns)
			} else {
				db.SetMaxOpenConns(*c.ConnConfig.MaxOpenConns)
			}
			if c.ConnConfig.ConnMaxLifetime == nil {
				db.SetConnMaxLifetime(cfg.DefaultConnConfig.ConnMaxLifetime)
			} else {
				db.SetConnMaxLifetime(*c.ConnConfig.ConnMaxLifetime)
			}
			if c.ConnConfig.MaxIdleConns == nil {
				db.SetMaxIdleConns(cfg.DefaultConnConfig.MaxIdleConns)
			} else {
				db.SetMaxIdleConns(*c.ConnConfig.MaxIdleConns)
			}
			if c.ConnConfig.ConnMaxIdleTime == nil {
				db.SetConnMaxIdleTime(cfg.DefaultConnConfig.ConnMaxIdleTime)
			} else {
				db.SetConnMaxIdleTime(*c.ConnConfig.ConnMaxIdleTime)
			}

			c.db = db

			// 裝載至清單中供後續取用
			dbList[k] = c
		})
	}
}

// DB : 取得DB
func (dbName DBName) DB() (*sql.DB, error) {
	return dbList[dbName].db, dbList[dbName].db.Ping()
}

// TX : 取得TX
func (dbName DBName) TX() (*sql.Tx, error) {
	return dbList[dbName].db.Begin()
}

// getDSN : 取得連線字串
func getDSN(dsnSource any) (dsn string, err error) {
	switch d := dsnSource.(type) {
	case Envkey:
		if _dsn, _err := d.getDSN(); _err != nil {
			return "", _err
		} else {
			dsn = _dsn
		}

	default:
		return "", errors.New("[func->getDSN] unkown DSNSource")
	}

	return dsn, nil
}

// getDSN : 從env取得連線字串
func (ek Envkey) getDSN() (string, error) {
	dsn := os.Getenv(string(ek))
	if dsn == "" {
		return "", errors.New("[func->getDSN] missing Env")
	}

	return dsn, nil
}
