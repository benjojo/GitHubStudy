package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*
To make the table you want:

CREATE TABLE `githubkeys` (
	`keyid` INT NOT NULL,
	`username` VARCHAR(96) NULL,
	`key` TEXT NULL,
	PRIMARY KEY (`keyid`),
	INDEX `username` (`username`)
)
COLLATE='utf8_bin'
ENGINE=InnoDB;

*/

func main() {
	con, err := sql.Open("mysql", "root:@/random")
	defer con.Close()
	check(err)
	ff, _ := os.OpenFile(os.Args[1], 0, 0)
	f := bufio.NewReader(ff)
	for {
		read_line, _ := f.ReadString('\n')
		if read_line == "" {
			break
		}
		fmt.Print(read_line + "\n")
		go processUser(read_line, con)
		time.Sleep(time.Second)
	}
	ff.Close()
	Wait := time.NewTimer(time.Second * 4)
	<-Wait.C
}

func processUser(username string, con *sql.DB) {
	rawkeys := getURL("https://github.com/" + username + ".keys")
	readKeys(rawkeys, username, con)
	fmt.Print(rawkeys)
	// First we will need to read the keys

}

func getURL(url string) string {
	conf := &tls.Config{InsecureSkipVerify: true}
	trans := &http.Transport{TLSClientConfig: conf}
	client := &http.Client{Transport: trans}
	resp, err := client.Get(url)
	check(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func readKeys(keychain string, username string, con *sql.DB) {
	for _, keystr := range strings.Split(keychain, "\n") {
		if keystr != "" {
			keysiz := getKeySize(keystr)
			storeKey(username, keystr, keysiz, con)
		}
	}
}

func getKeySize(key string) int {
	WriteFile("./"+SHAHash(key)+".key", key)
	app := "ssh-keygen"
	arg0 := "-l"
	arg1 := "-f"
	arg2 := "./" + SHAHash(key) + ".key"

	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()
	check(err)
	parts := strings.Split(string(out), " ")
	i, err := strconv.Atoi(parts[0])
	check(err)
	os.Remove(arg2)
	return int(i)
}

func WriteFile(path string, contents string) {
	d1 := []byte(contents)
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)
}

/* I have gathered that this is fairly standard. */
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func SHAHash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return string(sha)
}

func storeKey(username string, key string, keylen int, con *sql.DB) {
	_, e := con.Exec("INSERT INTO `random`.`githubkeys` (`username`, `keylen`, `key`) VALUES (?, ?, ?);", username, keylen, key)
	check(e)
}
