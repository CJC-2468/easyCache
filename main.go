package main

import (
	"easyCache/easycache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *easycache.Group {
	//把匿名函数转换为easycache.GetterFunc类型，传入的时候其实是作为一个getter，而getter是个接口类型
	return easycache.NewGroup("scores", 2, easycache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))
}

//启动缓存服务器：创建 HTTPPool，添加节点信息到哈希环，注册到 easy 中作为PeerPicker，
//启动 HTTP 服务（共3个端口，8001/8002/8003），用户不感知。
func startCacheServer(addr string, addrs []string, easy *easycache.Group) {
	peers := easycache.NewHTTPPool(addr)
	peers.Set(addrs...)
	easy.RegisterPeers(peers)
	log.Println("easycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

//启动一个 API 服务（端口 9999），与用户进行交互，用户感知。
func startAPIServer(apiAddr string, easy *easycache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := easy.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		},
	))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}
func main() {
	//指定端口启动 HTTP 服务。
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "easycache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	easy := createGroup()
	if api {
		go startAPIServer(apiAddr, easy)
	}
	startCacheServer(addrMap[port], []string(addrs), easy)
}
