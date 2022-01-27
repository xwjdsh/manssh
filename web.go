package manssh

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
)

//go:embed web/dist
var assetData embed.FS

type Resp struct {
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func WebServe(path string, addr string, cors bool) error {
	router := mux.NewRouter().StrictSlash(true)
	h := &webHandler{
		path: path,
		cors: cors,
	}

	var staticFS = fs.FS(assetData)
	htmlContent, err := fs.Sub(staticFS, "web/dist")
	if err != nil {
		return err
	}
	router.HandleFunc("/api/records", h.listRecords).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/record/{key}", h.deleteRecord).Methods(http.MethodDelete, http.MethodOptions)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(htmlContent))))
	return http.ListenAndServe(addr, router)
}

type webHandler struct {
	path string
	cors bool
}

func (h *webHandler) listRecords(w http.ResponseWriter, req *http.Request) {
	records, err := List(h.path, ListOption{})
	if err != nil {
		h.resp(w, &Resp{Err: err.Error()})
		return
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Alias > records[j].Alias
	})

	h.resp(w, &Resp{Data: records})
	return
}

func (h *webHandler) deleteRecord(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	if key == "" {
		h.resp(w, &Resp{Err: "key is empty"})
		return
	}
	_, err := Delete(h.path, key)
	if err != nil {
		h.resp(w, &Resp{Err: err.Error()})
		return
	}

	h.resp(w, nil)
	return
}

func (h *webHandler) resp(w http.ResponseWriter, r *Resp) {
	if h.cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	if r == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(r)
	w.Write(data)
}
