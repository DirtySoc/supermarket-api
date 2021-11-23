package main

import (
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
}

func TestAddProduce(t *testing.T) {
	produceHandler := newProduceHandlers()

	t.Run("Add new produce", func(t *testing.T) {
		// TODO: add body of new produce & check it exists after

		request, _ := http.NewRequest(http.MethodPost, "/produce", nil)
		response := httptest.NewRecorder()

		produceHandler.post(response, request)

		got := response.Code
		if got != http.StatusOK {
			t.Errorf("got %d, want %d", got, http.StatusOK)
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
