package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goshinobi/tor"
)

func processListHandle(w http.ResponseWriter, r *http.Request) {
	bin, err := json.Marshal(tor.GetWorkProxyList())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bin)
}

func createProxyHandle(w http.ResponseWriter, r *http.Request) {
	err := tor.StartProxy()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("success"))
}

func killAllProxyHandle(w http.ResponseWriter, r *http.Request) {
	result := ""
	for k, v := range tor.GetWorkProxyList() {
		err := tor.KillTorProxy(k)
		if err == nil {
			result += fmt.Sprintf("%d: %v\n", k, v)
		}
	}

	w.Write([]byte(result))
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/list", processListHandle)
	http.HandleFunc("/create", createProxyHandle)
	http.HandleFunc("/kill", killAllProxyHandle)
	http.ListenAndServe(":8080", nil)
}
