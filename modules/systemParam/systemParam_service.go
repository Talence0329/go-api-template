package systemparam

import (
	"backend/basic/database"
	"backend/modules/tools"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

// getOne
func getOne(db *sql.DB, key string) (SystemParamData, error) {
	pssq := database.PSSQ()
	data := SystemParamData{LastTime: time.Now()}

	query := pssq.Select(`name, value, category`).From(`systemparam`).Where(sq.Eq{"name": key})
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return SystemParamData{}, sqErr
	}

	switch err := db.QueryRow(sqll, args...).Scan(&data.Key, &data.Value, &data.Category); err {
	case sql.ErrNoRows:
		return SystemParamData{}, err
	case nil:
		return data, nil
	default:
		return SystemParamData{}, err
	}
}

// getAll
func getAll(db *sql.DB) (result []SystemParamData, err error) {
	pssq := database.PSSQ()
	query := pssq.Select(`name, value, category`).From(`systemparam`)

	_sql, _args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	_rows, err := db.Query(_sql, _args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := _rows.Close(); err != nil {
			tools.Log(tools.LogLevelError, "[systemParamManage/getListData/rows.Close] %v", err)
		}
	}()

	result = make([]SystemParamData, 0)
	for _rows.Next() {
		data := SystemParamData{LastTime: time.Now()}
		if err = _rows.Scan(&data.Key, &data.Value, &data.Category); err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	if err = _rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
