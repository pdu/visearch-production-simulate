package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dropbox/godropbox/errors"
)

// User is the APP's credential
type User struct {
	Username string `json:"username"`
	Access   string `json:"access"`
	Secret   string `json:"secret"`
}

// Users is the slice of User
type Users []User

var users Users

var host = flag.String("host", "http://visearch.visenze.com/", "The search endpoint")
var filenames = flag.String("files", "", "The playback logs, multiple log files separated by comma")
var qps = flag.Int("qps", 20, "The QPS of the load")

func loadCreds() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("ReadFile failed,filename:./creds.json error:%v\n", err)
		os.Exit(1)
	}

	json.Unmarshal(file, &users)
}

func getCred(username string) (string, string, error) {
	for _, user := range users {
		if user.Username == username {
			return user.Access, user.Secret, nil
		}
	}
	return "", "", errors.Newf("Credential not found for %v", username)
}

func main() {
	flag.Parse()

	loadCreds()

	files := strings.Split(*filenames, ",")
	requests, err := parseFiles(files)
	if err != nil {
		log.Printf("parseFiles failed,files:%v err:%v\n", files, err)
		os.Exit(1)
	}

	log.Println("begin playback")

	playback(requests)
}
