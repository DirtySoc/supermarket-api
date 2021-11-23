package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

func (h *produceHandlers) produce(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
		return
	case http.MethodPost:
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func (h *produceHandlers) get(w http.ResponseWriter, r *http.Request) {
	// Convert map to array of produce items.
	produce := make([]Produce, len(h.store))

	h.mu.Lock()
	i := 0
	for _, item := range h.store {
		produce[i] = item
		i++
	}
	h.mu.Unlock()

	jsonBytes, err := json.Marshal(produce)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // TODO
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *produceHandlers) post(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
}

func main() {
	produceHandlers := newProduceHandlers()

	router := http.NewServeMux()
	router.HandleFunc("/produce", produceHandlers.produce)
	log.Fatal(http.ListenAndServe(":6620", router))
}

type Produce struct {
	Name        string  `json:"name"`
	ProduceCode string  `json:"produceCode"`
	UnitPrice   float64 `json:"unitPrice"`
}

type produceHandlers struct {
	mu    sync.Mutex
	store map[string]Produce
}

func newProduceHandlers() *produceHandlers {
	return &produceHandlers{
		store: map[string]Produce{
			"A12T-4GH7-QPL9-3N4M": {Name: "Lettuce", ProduceCode: "A12T-4GH7-QPL9-3N4M", UnitPrice: 3.46},
			"E5T6-9UI3-TH15-QR88": {Name: "Peach", ProduceCode: "E5T6-9UI3-TH15-QR88", UnitPrice: 2.99},
			"YRT6-72AS-K736-L4AR": {Name: "Green Pepper", ProduceCode: "YRT6-72AS-K736-L4AR", UnitPrice: 0.79},
			"TQ4C-VV6T-75ZX-1RMR": {Name: "Gala Apple", ProduceCode: "TQ4C-VV6T-75ZX-1RMR", UnitPrice: 3.59},
		},
	}
}
