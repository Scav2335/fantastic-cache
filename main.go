package main

import (
	"flag"
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

func createGroup() *fancache.Group {
	return fancache.NewGroup("scores", 2<<10, fancache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addresses []string, group *fancache.Group) {
	peers := fancache.NewHTTPPool(addr)
	peers.Set(addresses...)
	group.RegisterPeers(peers)
	log.Println("fan-cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *fancache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(view.ByteSlice())
		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "fan-cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addresses []string
	for _, addr := range addrMap {
		addresses = append(addresses, addr)
	}

	group := createGroup()
	if api {
		go startAPIServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], addresses, group)
}
