package helper

import (
	"fmt"
	"github.com/sz47/shortenIt/structs"
	"log"
	"math/rand"
	"time"
)

func GenerateRandomString() string {
	rand.Seed(time.Now().UTC().UnixNano())
	chars := "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	for true {
		temp := fmt.Sprintf("%c%c%c", chars[rand.Intn(57)], chars[rand.Intn(57)], chars[rand.Intn(57)])

		structs.Mutex.RLock()
		_, ok := structs.LinkSet[temp]
		structs.Mutex.RUnlock()

		if !ok {
			structs.Mutex.Lock()
			defer structs.Mutex.Unlock()

			structs.LinkSet[temp] = true
			return temp
		}
	}
	return ""
}

func SyncAliasesWithTimer() {
	for range time.Tick(time.Hour * 24) {
		SyncAliases()
	}
}

func SyncAliases() {

	log.Println("Syncing DB")

	structs.Mutex.Lock()
	defer structs.Mutex.Unlock()

	rows, err := structs.DB.Query(fmt.Sprintf("SELECT key_url FROM %s;", structs.Config.DB.TableName))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Sync end")
}
