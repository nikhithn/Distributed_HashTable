package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	numChars := 20
	var currKey string = ""
	var currValue string = ""
	randi := 0
	var keys []string
	for j := 0; j < 100; j++ {
		for i := 0; i < numChars; i++ {
			randi = rand.Intn(26)
			currKey = currKey + string(randi+97)
			randi = rand.Intn(26)
			currValue = currValue + string(randi+97)
		}
		http.Get("http://localhost:8080/input.html?key=" + currKey + "&value=" + currValue)
		keys = append(keys, currKey)
		currKey = ""
		currValue = ""
	}

	for _, k := range keys {
		resp, _ := http.Get("http://localhost:8080/search.html?key=" + k)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}

}
