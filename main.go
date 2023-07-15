package main

import (
	"fmt"
	"geeCache/geecache"
	"log"
	"net/http"
)

var luckyNumsDB = map[string]string{
	"chubby": "7777",
	"rshao":  "8888",
	"Jack":   "6666",
	"Others": "23132",
}

func main() {
	geecache.NewGroup("luckNums", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Printf("[%s] is searching key: %s", "luckNums", key)
			if v, ok := luckyNumsDB[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key:%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at:", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
