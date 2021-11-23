// Package main provides the main implementation of the supermarket api
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

// produce handles different http methods on the /produce endpoint.
func (h *produceHandlers) produce(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	switch r.Method {
	case http.MethodGet:
		switch len(parts) {
		case 2:
			h.get(w, r)
		case 3:
			h.getById(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	case http.MethodPost:
		h.post(w, r)
		return
	case http.MethodDelete:
		h.delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

// get handles get request on the /produce endpoint
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
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// getById handles get request for a produce by produceCode
func (h *produceHandlers) getById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")

	h.mu.Lock()
	item, found := h.store[parts[2]]
	h.mu.Unlock()

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// post handles post requests on the /produce endpoint.
// post accepts valid json that conforms to the Produce spec. Arrays of
// Produce or single Produce JSON is accepted.
func (h *produceHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "need content-type 'application/json', but got '%s'", ct)
		return
	}

	var newProduce []Produce // TODO handle single json object and json arrays
	err = json.Unmarshal(bodyBytes, &newProduce)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.mu.Lock()
	for _, item := range newProduce {
		_, exists := h.store[item.ProduceCode]
		if exists {
			fmt.Fprintf(w, "Warning: Duplicate ProduceCodes detected: Item %s already exists.", item.ProduceCode)
		}
		h.store[item.ProduceCode] = item
	}
	defer h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func (h *produceHandlers) delete(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func main() {
	produceHandlers := newProduceHandlers()

	router := http.NewServeMux()
	router.HandleFunc("/produce", produceHandlers.produce)
	router.HandleFunc("/produce/", produceHandlers.produce)
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
