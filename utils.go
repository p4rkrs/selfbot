package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Global Variables
var (
	codename = "Blueberry"
	version  = "3.1"
	vColor   = 0x3498db
	api      = "https://api.coinbase.com/v2/prices/spot?currency="
)

//GetPrice returns price
func GetPrice(currency string) (string, error) {
	resp, err := http.Get(api + currency)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Invalid response for " + currency + " | <" + resp.Status + ">")
		return "", errors.New("Invalid response for " + currency + ", " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	data := map[string]map[string]string{}
	json.Unmarshal(body, &data)
	return data["data"]["amount"], nil
}

//GetTime returns unix timestamp in milliseconds
func GetTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
