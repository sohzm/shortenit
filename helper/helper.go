package helper

import (
	"fmt"
	"github.com/sz47/shortenIt/structs"
	"log"
	"math/rand"
	"time"
)

// GenerateRandomString will generate random string of length 3 which isnt 
// already present in the structs.LinkSet
func GenerateRandomString() string {
	rand.Seed(time.Now().UTC().UnixNano())

	// notice we didn't include l, I, 0, O in the chars as these often cause 
	// confusion due to the fact that they look kind of same under most fonts.
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

// SyncAliasesWithTimer will call SyncAliases() every 24 hour, maintaining a 
// local store allows us to provide a new link pretty quickly without doing a DB 
// call struct.LinkSet will have the links we had when last sync was done + at 
// max 20_000 more links, which makes it have the maximum size of 40_000. We 
// have ample of variations of chars so having a local store double the size of
// actual db wont cause any problems but will result in a speed boost.
func SyncAliasesWithTimer() {
	for range time.Tick(time.Hour * 24) {
		SyncAliases()
	}
}

// SyncAliases ... syncs the aliases.
// The reason Mutex is used is that we dont want to add values to the LinkSet
// when we're syncing with the db.
func SyncAliases() {

	log.Println("Syncing DB.")

	structs.Mutex.Lock()
	defer structs.Mutex.Unlock()

	// empty set
	structs.LinkSet = make(map[string]bool)

	rows, err := structs.DB.Query(fmt.Sprintf("SELECT key_url FROM %s;", structs.Config.DB.TableName))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			log.Fatal(err)
		}
		structs.LinkSet[key] = true
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Syncing Ended Successfully.")
}
