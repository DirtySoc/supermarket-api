package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/health", healthCheck)
	router.HandleFunc("/produce", handleProduce)
	log.Fatal(http.ListenAndServe(":6620", router))
}

func handleProduce(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(ProduceDB)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Health endpoint hit.")
	fmt.Println("Endpoint Hit: healthCheck")
}

type Produce struct {
	Name        string  `json:"name"`
	ProduceCode string  `json:"produceCode"`
	UnitPrice   float64 `json:"unitPrice"`
}

var ProduceDB = []Produce{
	{Name: "Lettuce", ProduceCode: "e0b4-d9f6-4ddf-bc69", UnitPrice: 3.45},
	{Name: "Tomatoe", ProduceCode: "e0A4-d9f6-4dqf-bc6t", UnitPrice: 4.78},
	{Name: "Banana", ProduceCode: "e0b4-dpf6-4dgf-bc77", UnitPrice: 2.50},
}
