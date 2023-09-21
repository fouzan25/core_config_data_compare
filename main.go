package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func isYes(selected string) bool {
	normalized := strings.ToLower(selected)
	return normalized == "y" || normalized == "yes"
}

func main() {
	loadENV()
	tempDBConnection := Config{
		Username:     os.Getenv("TEMP_DB_USER"),
		Password:     os.Getenv("TEMP_DB_PASSWORD"),
		DatabaseName: os.Getenv("TEMP_DB_NAME"),
		Host:         fmt.Sprintf("%s:%s", os.Getenv("TEMP_DB_HOST"), os.Getenv("TEMP_DB_PORT")),
	}
	uatDBConnection := Config{
		Username:     os.Getenv("UAT_DB_USER"),
		Password:     os.Getenv("UAT_DB_PASSWORD"),
		DatabaseName: os.Getenv("UAT_DB_NAME"),
		Host:         fmt.Sprintf("%s:%s", os.Getenv("UAT_DB_HOST"), os.Getenv("UAT_DB_PORT")),
	}

	tempDb, err1 := newConnection(tempDBConnection)
	uatDb, err2 := newConnection(uatDBConnection)

	if err1 != nil {
		panic(err1)
	}
	defer tempDb.Close()
	if err2 != nil {
		panic(err2)
	}
	defer uatDb.Close()
	tempDBqueryResult, err := tempDb.getPathAndValues()

	if err1 != nil {
		panic(err)
	}

	for _, tempItem := range tempDBqueryResult {
		pathExist, err := uatDb.pathExist(tempItem)
		if err != nil {
			panic(err)
		}
		if pathExist.path != "" {
			tempDb.updatePathValue(pathExist)
		} else {
			var selected string
			var insertedValues string
			fmt.Printf("path: %v\t", tempItem.path)
			fmt.Printf("value: %v\n", tempItem.value.String)
			fmt.Printf("do you insert value manually ? y/n\n")
			fmt.Scan(&selected)
			if isYes(selected) {
				fmt.Printf("enter new value: \t")
				fmt.Scan(&insertedValues)
				fmt.Println(insertedValues)
				tempItem.value = sql.NullString{String: insertedValues, Valid: true}
				tempDb.updatePathValue(tempItem)
			}

		}
	}
}
