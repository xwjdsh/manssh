package web

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"

	"github.com/xwjdsh/manssh"
)

type Resp struct {
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

var (
	path string
	cors bool
)

func Serve(p string, addr string, allowCors bool) error {
	path = p
	cors = allowCors
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/records", listRecords).Methods(http.MethodGet, http.MethodOptions)

	return http.ListenAndServe(addr, router)
}

func listRecords(w http.ResponseWriter, req *http.Request) {
	records, err := manssh.List(path, manssh.ListOption{})
	if err != nil {
		resp(w, &Resp{Err: err.Error()})
		return
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Alias < records[j].Alias
	})

	resp(w, &Resp{Data: records})
	return
}

func resp(w http.ResponseWriter, r *Resp) {
	if cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(r)
	w.Write(data)
}
