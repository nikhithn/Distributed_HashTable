package main

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"os"
	"strconv"
)

type KVPair struct {
	key   string
	value string
}

var table [][]KVPair

func main() {

	// Parse arguments
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments.  Please just enter the size of the hash table")
		os.Exit(1)
	}
	args := os.Args[1]
	tableSize, err := strconv.Atoi(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if tableSize <= 0 {
		fmt.Println("Please enter a size greater than 0")
		os.Exit(1)
	}

	table = make([][]KVPair, tableSize, tableSize)

	// Start http server
	http.HandleFunc("/input.html", input)
	http.HandleFunc("/search.html", search)
	http.ListenAndServe(":8080", nil)
}

func input(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		params := r.URL.Query()
		if params["key"] != nil && params["value"] != nil {
			fmt.Println("Key: ", params["key"][0])
			fmt.Println("Value: ", params["value"][0])
			insert(params["key"][0], params["value"][0])
		}
		http.ServeFile(w, r, "input.html")
	}
}

func insert(key string, value string) {
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32() % uint32(len(table))

	kv := KVPair{key, value}
	if table[hash] == nil {
		table[hash] = []KVPair{kv}
		fmt.Println(hash, len(table[hash]))
	} else {
		for i, e := range table[hash] {
			if e.key == key {
				e.value = value
				table[hash][i] = e
				fmt.Println(hash, len(table[hash]))
				return
			}
		}
		table[hash] = append(table[hash], kv)
		fmt.Println(hash, len(table[hash]))
	}
}

func lookup(key string) (string, string) {
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32() % uint32(len(table))
	if table[hash] == nil {
		fmt.Println("Key not found")
		return "", "Key not found"
	} else {
		for _, e := range table[hash] {
			if e.key == key {
				return e.value, ""
			}
		}
		fmt.Println("Key not found")
		return "", "Key not found"
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		params := r.URL.Query()
		if params["key"] != nil {
			fmt.Println("Key: ", params["key"][0])
			value, err := lookup(params["key"][0])
			if err == "" {
				fmt.Println("Value: ", value)
				w.Write([]byte(value))
			} else {
				fmt.Println(err)
				w.Write([]byte(err))
			}
		} else {
			http.ServeFile(w, r, "search.html")
		}
	}
}
