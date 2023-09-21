package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	*sql.DB
}

func newConnection(credentials Config) (*Mysql, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		credentials.Username,
		credentials.Password,
		credentials.Host,
		credentials.DatabaseName,
	)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, errors.New("failed to connect to the database")
	}
	return &Mysql{db}, nil
}

func (db *Mysql) getPathAndValues() ([]dataSet, error) {
	set := dataSet{}
	data := []dataSet{}
	rows, err := db.Query("SELECT path,value FROM core_config_data")
	if err != nil {
		rows.Close()
		return []dataSet{}, err
	}

	for rows.Next() {
		var path string

		var value sql.NullString

		if err := rows.Scan(&path, &value); err != nil {
			panic(err)
		}
		set.path = path
		set.value = value

		data = append(data, set)
	}

	return data, nil

}

func (db *Mysql) updatePathValue(data dataSet) {
	query := "UPDATE core_config_data SET value = ? WHERE path = ?"
	result, err := db.Exec(query, data.value, data.path)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
}

func (db *Mysql) pathExist(data dataSet) (dataSet, error) {
	result := dataSet{}
	query := "SELECT path,value FROM core_config_data where path = ?"
	rows, err := db.Query(query, data.path)
	if err != nil {
		rows.Close()
		return dataSet{}, err
	}
	for rows.Next() {
		var path string
		var value sql.NullString
		if err := rows.Scan(&path, &value); err != nil {
			panic(err)
		}
		if (path) != "" {
			result.path = path
			result.value = value
		}

	}
	return result, nil
}
