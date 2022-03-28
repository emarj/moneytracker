package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const apiKey = "keyAuR8F3wLAUXZAL"

func main() {

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://api.airtable.com/v0/appfbyqOSj86V5Wku/Comuni", nil)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
	}
	fmt.Println(buf.String())

}
