package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetProduce(t *testing.T) {

	produceHandler := newProduceHandlers()

	initProduceDB := []Produce{}
	for _, item := range produceHandler.store {
		initProduceDB = append(initProduceDB, item)
	}

	t.Run("GET all produce", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/produce", nil)
		response := httptest.NewRecorder()
		produceHandler.get(response, request)

		var got []Produce
		json.Unmarshal(response.Body.Bytes(), &got)

		if !isEqual(got, initProduceDB) {
			t.Errorf("got %+v, want %+v", got, initProduceDB)
		}
	})

	t.Run("GET produce by ID", func(t *testing.T) {
		// TODO check for 200
	})

	t.Run("GET produce with fake ID", func(t *testing.T) {
		// TODO check for 404
	})
}

func TestAddProduce(t *testing.T) {
	produceHandler := newProduceHandlers()

	// TODO: fix this with fancy custom unmarshaler probably
	// t.Run("Add single new produce", func(t *testing.T) {
	// 	newProduce := Produce{
	// 		Name:        "Red Apple",
	// 		ProduceCode: "RRRR-VV6T-75ZX-1RMR",
	// 		UnitPrice:   3.44,
	// 	}

	// 	reqBody, _ := json.Marshal(newProduce)
	// 	request, _ := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(reqBody))
	// 	request.Header.Add("content-type", "application/json")
	// 	response := httptest.NewRecorder()
	// 	produceHandler.post(response, request)

	// 	got := response.Code
	// 	if got != http.StatusCreated {
	// 		t.Errorf("got %d, want %d", got, http.StatusOK)
	// 	}
	// })

	t.Run("Add produce that already exists", func(t *testing.T) {
		newProduce := []Produce{
			{Name: "Lettuce", ProduceCode: "A12T-4GH7-QPL9-3N4M", UnitPrice: 3.46},
		}

		reqBody, _ := json.Marshal(newProduce)
		request, _ := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(reqBody))
		request.Header.Add("content-type", "application/json")
		response := httptest.NewRecorder()
		produceHandler.post(response, request)

		got := response.Code
		if got != http.StatusCreated { // TODO: is this really what we want here?
			t.Errorf("got %d, want %d", got, http.StatusCreated)
		}
	})

	t.Run("Add multiple new produce", func(t *testing.T) {
		newProduce := []Produce{
			{Name: "Red Apple", ProduceCode: "RRRR-VV6T-75ZX-1RMR", UnitPrice: 3.44},
			{Name: "Blue Apple", ProduceCode: "BBBB-VV6T-75ZX-1RMR", UnitPrice: 40.12},
			{Name: "Purple Apple", ProduceCode: "PPPP-VV6T-75ZX-1RMR", UnitPrice: 43.99},
		}

		reqBody, _ := json.Marshal(newProduce)
		request, _ := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(reqBody))
		request.Header.Add("content-type", "application/json")
		response := httptest.NewRecorder()
		produceHandler.post(response, request)

		got := response.Code
		if got != http.StatusCreated {
			t.Errorf("got %d, want %d", got, http.StatusOK)
		}
	})

	t.Run("Add items concurrently", func(t *testing.T) {
		// TODO
	})

	t.Run("Add item with crazy JSON body", func(t *testing.T) {
		// TODO
	})
}

func TestDeleteProduce(t *testing.T) {
	produceHandler := newProduceHandlers()

	t.Run("Delete existing produce", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/produce/A12T-4GH7-QPL9-3N4M", nil)
		response := httptest.NewRecorder()

		produceHandler.delete(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("got %d, want %d", response.Code, http.StatusOK)
		}
	})

	t.Run("Delete non-existant produce", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/produce/non-existant-produce", nil)
		response := httptest.NewRecorder()

		produceHandler.delete(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("got %d, want %d", response.Code, http.StatusNotFound)
		}
	})
}

// Check and ensure that all elements exists in another slice and
// check if the length of the slices are equal.
func isEqual(aa, bb []Produce) bool {
	eqCtr := 0
	for _, a := range aa {
		for _, b := range bb {
			if reflect.DeepEqual(a, b) {
				eqCtr++
			}
		}
	}
	if eqCtr != len(bb) || len(aa) != len(bb) {
		return false
	}
	return true
}
