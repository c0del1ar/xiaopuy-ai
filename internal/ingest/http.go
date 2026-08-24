package ingest

import (
	"encoding/json"
	"net/http"
)

type HTTPHandler struct { Service *Service }

type crawlHTTPResponse struct {
	Visited   int `json:"visited"`
	Skipped   int `json:"skipped"`
	Documents int `json:"documents"`
}

func (h *HTTPHandler) Crawl(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Service == nil { http.Error(w, "ingestion service is not configured", http.StatusInternalServerError); return }
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	var req CrawlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid JSON body", http.StatusBadRequest); return }
	result, err := h.Service.Crawl(r.Context(), req)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(crawlHTTPResponse{Visited: result.Visited, Skipped: result.Skipped, Documents: len(result.Documents)})
}
