package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"net/http"

	"github.com/sz47/shortenIt/databaseHelper"
	"github.com/sz47/shortenIt/handlers"
	"github.com/sz47/shortenIt/helper"
	"github.com/sz47/shortenIt/structs"

	_ "github.com/go-sql-driver/mysql"
)

func initConfigurationStruct() {
	structs.Mutex = sync.RWMutex{}
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &structs.Config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// initialize
	initConfigurationStruct()
	databaseHelper.InitDB()

	// sync
	helper.SyncAliases()
	go helper.SyncAliasesWithTimer()

	// add handlers
	http.HandleFunc("/new", handlers.NewHandler)
	http.HandleFunc("/", handlers.RedirectHandler)
	log.Println("Server Ready at Port:", structs.Config.AppPort)
	log.Fatal(http.ListenAndServe(structs.Config.AppPort, nil))
}
