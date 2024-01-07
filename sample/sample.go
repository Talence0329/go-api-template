package sample

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type AddStruct struct {
	name string
	age  int
}

func addService(db *sql.DB, addData AddStruct) (int, error) {
	query := sq.Insert("test").Columns("name", "age").Values(addData.name, addData.age)

	_sql, args, sqErr := query.ToSql()
	if sqErr != nil {
		return 0, sqErr
	}

	result, execErr := db.Exec(_sql, args...)
	if execErr != nil {
		return 0, execErr
	}

	count, err := result.RowsAffected()
	if err != nil {
		return int(count), err
	}

	return int(count), nil
}

func updateService(db *sql.DB, id int) (int, error) {
	query := sq.Update("test").Set("Times", 0).Where(sq.Eq{"id": id})
	_sql, args, sqErr := query.ToSql()
	if sqErr != nil {
		return 0, sqErr
	}

	result, execErr := db.Exec(_sql, args...)
	if execErr != nil {
		return 0, execErr
	}

	count, err := result.RowsAffected()
	if err != nil {
		return int(count), err
	}

	return int(count), nil
}

type getTermStruct struct {
	name string
}
type getDataStruct struct {
	name string
	age  int
}

func getService(db *sql.DB, term getTermStruct) (getDataStruct, error) {
	data := getDataStruct{}

	where := sq.Eq{
		"name": term.name,
	}
	query := sq.Select(`
		name,
		age
	`).From(`test
	`).Where(where)
	_sql, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}

	switch err := db.QueryRow(_sql, args...).Scan(
		&data.name,
		&data.age,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		return data, nil
	default:
		return data, err
	}
}
