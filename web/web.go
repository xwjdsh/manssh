package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/xwjdsh/manssh"
)

type Resp struct {
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

var path string

func Serve(p string, addr string) error {
	path = p
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/records", listRecords).Methods(http.MethodGet)

	return http.ListenAndServe(addr, router)
}

func listRecords(w http.ResponseWriter, req *http.Request) {
	records, err := manssh.List(path, manssh.ListOption{})
	if err != nil {
		resp(w, &Resp{Err: err.Error()})
		return
	}
	resp(w, &Resp{Data: records})
	return
}

func resp(w http.ResponseWriter, r *Resp) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(r)
	w.Write(data)
}
