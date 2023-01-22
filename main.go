package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
)

func getConsulKVClient() *api.KV {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	// Get a handle to the KV API
	kv := client.KV()
	return kv
}

func Greeting(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received.")
	io.WriteString(w, "Hello from kubernetes pod!")
}

func getAppConfig(w http.ResponseWriter, r *http.Request) {
	kv := getConsulKVClient()
	pair, _, err := kv.Get("REDIS_MAXCLIENTS", nil)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, fmt.Sprintf("{\"%s\": \"%s\"}", pair.Key, pair.Value))
}

func loadAppConfig() {
	kv := getConsulKVClient()
	// PUT a new KV pair
	p := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000")}
	_, err := kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
	// Lookup the pair
	pair, _, err := kv.Get("REDIS_MAXCLIENTS", nil)
	if err != nil {
		panic(err)
	}
	log.Printf("KV: %v %s\n", pair.Key, pair.Value)
}

func main() {
	loadAppConfig()
	http.HandleFunc("/hello", Greeting)
	http.HandleFunc("/app-config", getAppConfig)
	log.Println("Server is starting...")
	err := http.ListenAndServe(":2021", nil)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	} else {
		log.Println("Server started.")
	}
}
