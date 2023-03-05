package structs

import (
	"database/sql"
	"sync"
)

var (
	DB      *sql.DB
	Mutex   sync.RWMutex
	Config  Configuration
	LinkSet map[string]bool = map[string]bool{}
)

type NewLinkRequest struct {
	NewLink string `json:"newLink"`
}

type NewLinkResponse struct {
	ShortenedAlias string `json:"shortenedAlias"`
	WillExpire     string `json:"willExpire"`
	Error          string `json:"error"`
}

type Configuration struct {
	DB       DBConfig `json:"db"`
	AppPort  string   `json:"appPort"`
	URLLimit int      `json:"urlLimit"`
}

type DBConfig struct {
	UserName       string `json:"userName"`
	Password       string `json:"password"`
	DBName         string `json:"dbName"`
	TableName      string `json:"tableName"`
	DBURL          string `json:"dbURL"`
	DBPort         string `json:"dbPort"`
	EventName      string `json:"eventName"`
	MaxAllowedTime string `json:"maxAllowedTime"`
}

type Link struct {
	KeyURL     string
	ValueURL   sql.NullString
	LastUpdate string
}
