package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	ff, _ := os.OpenFile(os.Args[1], 0, 0)
	f := bufio.NewReader(ff)
	for {
		read_line, _ := f.ReadString('\n')
		if read_line == "" {
			break
		}
		fmt.Print(read_line + "\n")
		go processUser(read_line)
	}
	ff.Close()
	Wait := time.NewTimer(time.Second * 4)
	<-Wait.C
}

func processUser(username string) {
	fmt.Print(getURL("https://github.com/" + username + ".keys"))
}

func getURL(url string) string {
	conf := &tls.Config{InsecureSkipVerify: true}
	trans := &http.Transport{TLSClientConfig: conf}
	client := &http.Client{Transport: trans}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func readKeys(keychain string) {

}

func getKeySize(key string) int {
	return 1
}
