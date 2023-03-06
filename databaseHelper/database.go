package databaseHelper

import (
	"database/sql"
	"fmt"
	"github.com/sz47/shortenIt/structs"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB will initialize database if not already present and create a links
// table if not already present.
func InitDB() {

	var err error

	structs.DB, err = sql.Open(
		"mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/",
			structs.Config.DB.UserName, structs.Config.DB.Password,
			structs.Config.DB.DBURL, structs.Config.DB.DBPort,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = structs.DB.Exec(
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;",
			structs.Config.DB.DBName),
	)
	if err != nil {
		log.Fatal(err)
	}

	structs.DB, err = sql.Open(
		"mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			structs.Config.DB.UserName, structs.Config.DB.Password,
			structs.Config.DB.DBURL, structs.Config.DB.DBPort,
			structs.Config.DB.DBName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = structs.DB.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s.%s "+
				"(key_url varchar(3), value_url text, last_update DATETIME, UNIQUE(key_url));",
			structs.Config.DB.DBName, structs.Config.DB.TableName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = structs.DB.Exec(fmt.Sprintf("use %s;", structs.Config.DB.DBName))
	if err != nil {
		log.Fatal(err)
	}

	_, err = structs.DB.Exec(
		fmt.Sprintf(
			"CREATE EVENT IF NOT EXISTS %s ON SCHEDULE EVERY 1 MINUTE DO "+
				"DELETE FROM %s WHERE DATE_ADD(last_update,INTERVAL %s) < NOW();",
			structs.Config.DB.EventName, structs.Config.DB.TableName,
			structs.Config.DB.MaxAllowedTime,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ShortenIt :: Database Successfully Initialized")

}
