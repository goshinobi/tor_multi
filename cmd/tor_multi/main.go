package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func addProxyHandle(w http.ResponseWriter, r *http.Request) {
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

func killProxyHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	n, err := strconv.Atoi(r.PostFormValue("n"))
	if err != nil {
		bin, _ := json.Marshal(struct {
			Status string
			Err    error
		}{
			"failed",
			err,
		})
		w.Write(bin)
		return
	}

	bin, _ := json.Marshal(struct {
		Status string
		N      int
	}{
		"failed",
		n,
	})
	w.Write(bin)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/list", processListHandle)
	http.HandleFunc("/add", addProxyHandle)
	http.HandleFunc("/killAll", killAllProxyHandle)
	http.HandleFunc("/kill", killProxyHandle)
	http.ListenAndServe(":8080", nil)
}
