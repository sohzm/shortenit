package handlers

import (
	"fmt"
	"github.com/sz47/shortenIt/helper"
	"github.com/sz47/shortenIt/structs"
	"log"
	"time"

	"encoding/json"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

func NewHandler(w http.ResponseWriter, r *http.Request) {

	inputBody := structs.NewLinkRequest{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &inputBody)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	newLink := inputBody.NewLink

	// handling some user side errors
	newLink = strings.ReplaceAll(newLink, " ", "")
	if len(newLink) < 8 {
		response := structs.NewLinkResponse{
			Error: "Your link is already small enough.",
		}
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, string(b))
		return
	}

	if newLink[0:7] != "http://" && newLink[0:8] != "https://" {
		newLink = "http://" + newLink
	}

	// parsing URL
	parsedURL, err := url.Parse(newLink)
	if err != nil {
		response := structs.NewLinkResponse{
			Error: err.Error(),
		}
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, string(b))
		return
	}
	if !parsedURL.IsAbs() {
		response := structs.NewLinkResponse{
			Error: "Unable to shorten that link. It is not a valid url.",
		}
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, string(b))
		return
	}

	// valid links, we can continue
	// check rows in db

	var n int
	numberOfRows := structs.DB.QueryRow("SELECT COUNT(*) FROM links;")
	if err := numberOfRows.Scan(&n); err != nil {
		log.Fatal(err)
	}

	if n >= structs.Config.URLLimit {
		response := structs.NewLinkResponse{
			Error: "Maximum Number of shortened links reached, try again in few hours.",
		}
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, string(b))
		return
	}

	// add to db
	newAlias := helper.GenerateRandomString()

	_, err = structs.DB.Exec("INSERT INTO links VALUES (?, ?, now())", newAlias, newLink)
	if err != nil {
		log.Fatal(err)
	}

	response := structs.NewLinkResponse{
		ShortenedAlias: newAlias,
		WillExpire:     fmt.Sprintf("%v", time.Now().Add(time.Hour*24)),
	}
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
	return
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Home Page Visit
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./index.html")
		return
	}

	p := strings.Split(r.URL.Path, "/")
	row := structs.DB.QueryRow("SELECT * FROM links where key_url = ?", p[1])
	var temp structs.Link
	if err := row.Scan(&temp.KeyURL, &temp.ValueURL, &temp.LastUpdate); err != nil {
		log.Println(err.Error())
		return
	}

	val := temp.ValueURL.String

	http.Redirect(w, r, val, http.StatusPermanentRedirect)
	return
}
