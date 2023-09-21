package main

import "database/sql"

type Config struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
}

type dataSet struct {
	path  string
	value sql.NullString
}
