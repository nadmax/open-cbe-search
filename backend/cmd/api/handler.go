package main

import (
	"encoding/json"
	"net/http"
	"github.com/nadmax/open-cbe-search/core/engine/postgres"
)

type API struct {
	DB *postgres.Client
}

type SearchRequest struct {
	Query string `json:"q"`
	Table string `json:"table"`
}

func NewAPI(db *postgres.Client) *API {
	return &API{ DB: db}
}

func (api *API) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Query == "" || req.Table == "" {
		http.Error(w, "Missing 'q' or 'table' field in JSON", http.StatusBadRequest)
		return
	}

	results, err := api.DB.SearchTable(req.Table, req.Query)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}