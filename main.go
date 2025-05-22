package main

import (
	"fmt"
	"log"
	"net/http"
	fancache "scav.abc/fantastic-cache/fan-cache"
)

var db = map[string]string{
	"Tom":   "630",
	"Jack":  "589",
	"Sam":   "567",
	"Kevin": "739",
}

func main() {
	fancache.NewGroup("scores", 2<<10, fancache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := fancache.NewHTTPPool(addr)
	log.Println("fan-cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
