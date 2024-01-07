package dbconn

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

const TESTDB DBName = "TEST_DB"

var testMock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	if db, mock, err := sqlmock.New(); err != nil {
		fmt.Printf("TestMain error = %v", err)
		return
	} else {
		testMock = mock
		dbList[TESTDB] = DBConfig{
			DSNSource: Envkey(TESTDB),
			db:        db,
		}
	}

	m.Run()
}
func TestDBName_DB(t *testing.T) {
	tests := []struct {
		name    string
		dbName  DBName
		wantErr bool
	}{
		{
			name:    "測試取得Tx",
			dbName:  TESTDB,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dbName.DB()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBName.DB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDBName_TX(t *testing.T) {
	tests := []struct {
		name    string
		dbName  DBName
		wantErr bool
	}{
		{
			name:    "測試取得Tx",
			dbName:  TESTDB,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMock.ExpectBegin()
			_, err := tt.dbName.TX()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBName.TX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
