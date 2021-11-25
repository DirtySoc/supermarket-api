// Package main provides the main implementation of the supermarket api
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"sync"

	"github.com/gorilla/mux"
)

type Produce struct {
	Name        string  `json:"name"`
	ProduceCode string  `json:"produceCode"`
	UnitPrice   float64 `json:"unitPrice"`
}

type storeHandlers struct {
	produceInvMu sync.Mutex
	produceInv   map[string]Produce
}

func newStoreHandlers() *storeHandlers {
	return &storeHandlers{
		produceInv: map[string]Produce{
			"A12T-4GH7-QPL9-3N4M": {Name: "Lettuce", ProduceCode: "A12T-4GH7-QPL9-3N4M", UnitPrice: 3.46},
			"E5T6-9UI3-TH15-QR88": {Name: "Peach", ProduceCode: "E5T6-9UI3-TH15-QR88", UnitPrice: 2.99},
			"YRT6-72AS-K736-L4AR": {Name: "Green Pepper", ProduceCode: "YRT6-72AS-K736-L4AR", UnitPrice: 0.79},
			"TQ4C-VV6T-75ZX-1RMR": {Name: "Gala Apple", ProduceCode: "TQ4C-VV6T-75ZX-1RMR", UnitPrice: 3.59},
		},
	}
}

// getProduceInvSlice returns a slice containing all of
// the produce entries sorted by name.
func (h *storeHandlers) getProduceInvSlice() []Produce {
	produce := make([]Produce, len(h.produceInv))
	h.produceInvMu.Lock()
	i := 0
	for _, item := range h.produceInv {
		produce[i] = item
		i++
	}
	h.produceInvMu.Unlock()

	// Sorting Produce by their Name
	// Using Slice() function
	sort.Slice(produce, func(p, q int) bool {
		return produce[p].Name < produce[q].Name
	})

	return produce
}

// updateProduceInv adds or updates produce entries in the DB
func (h *storeHandlers) updateProduceInv(newProduce []Produce) {
	h.produceInvMu.Lock()
	defer h.produceInvMu.Unlock()
	for _, item := range newProduce {
		h.produceInv[item.ProduceCode] = item
	}
}

// returns an error if invalid product IDs are detected in []Produce
// todo handle error from regexp.Compile gracefully
func (h *storeHandlers) validateProduceIDs(p []Produce) error {
	prodCodeRegexp, _ := regexp.Compile("^(?:[a-zA-Z0-9]{4}-){3}[a-zA-Z0-9]{4}$")
	for _, item := range p {
		matched := prodCodeRegexp.MatchString(item.ProduceCode)
		if !matched {
			return fmt.Errorf("invalid product code detected")
		}
	}
	return nil
}

// handles GET /produce requests
func (h *storeHandlers) getProduce(w http.ResponseWriter, r *http.Request) {
	resBytes, err := json.Marshal(h.getProduceInvSlice())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Errorf("unable to marshal data into JSON"))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

// handles GET /produce/{id} requests
func (h *storeHandlers) getProduceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Errorf("id not present in request"))
		return
	}

	h.produceInvMu.Lock()
	item, found := h.produceInv[id]
	h.produceInvMu.Unlock()

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resBytes, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

// handles POST /produce requests with body containing
// JSON array of new/updated produce entries
func (h *storeHandlers) postProduce(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Body does not exist.")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	if len(bodyBytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Body is empty.")
		return
	}

	var newProduce []Produce
	if err := json.Unmarshal(bodyBytes, &newProduce); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	if err := h.validateProduceIDs(newProduce); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	h.updateProduceInv(newProduce)
	w.WriteHeader(http.StatusCreated)
}

// handles DELETE /produce/:id requests
func (h *storeHandlers) deleteProduceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Errorf("id not present in request"))
		return
	}

	h.produceInvMu.Lock()
	defer h.produceInvMu.Unlock()

	delete(h.produceInv, id)
}

func setupRouter() *mux.Router {
	storeHandlers := newStoreHandlers()

	r := mux.NewRouter()
	r.HandleFunc("/produce", storeHandlers.getProduce).Methods(http.MethodGet)
	r.HandleFunc("/produce/{id}", storeHandlers.getProduceByID).Methods(http.MethodGet)
	r.HandleFunc("/produce", storeHandlers.postProduce).Methods(http.MethodPost)
	r.HandleFunc("/produce/{id}", storeHandlers.deleteProduceByID).Methods(http.MethodDelete)

	return r
}

func main() {
	r := setupRouter()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":6620", r))
}
