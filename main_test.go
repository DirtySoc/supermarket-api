package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var initProduce = []Produce{
	{Name: "Gala Apple", ProduceCode: "TQ4C-VV6T-75ZX-1RMR", UnitPrice: 3.59},
	{Name: "Green Pepper", ProduceCode: "YRT6-72AS-K736-L4AR", UnitPrice: 0.79},
	{Name: "Lettuce", ProduceCode: "A12T-4GH7-QPL9-3N4M", UnitPrice: 3.46},
	{Name: "Peach", ProduceCode: "E5T6-9UI3-TH15-QR88", UnitPrice: 2.99},
}

func TestGetProduce(t *testing.T) {
	router := setupRouter()

	t.Run("get all produce", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/produce", nil)
		router.ServeHTTP(w, req)

		expectedRes, _ := json.Marshal(initProduce)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedRes), w.Body.String())
	})

	t.Run("get produce by ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/produce/TQ4C-VV6T-75ZX-1RMR", nil)
		router.ServeHTTP(w, req)

		expectedRes, _ := json.Marshal(initProduce[0])

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedRes), w.Body.String())
	})

	t.Run("get produce by ID 404", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/produce/kalhefid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPostProduce(t *testing.T) {
	router := setupRouter()

	var newProduce = []Produce{
		{Name: "Red Apple", ProduceCode: "RRRR-VV6T-75ZX-1RMR", UnitPrice: 3.44},
		{Name: "Blue Apple", ProduceCode: "BBBB-VV6T-75ZX-1RMR", UnitPrice: 3.44},
		{Name: "Green Apple", ProduceCode: "GGGG-VV6T-75ZX-1RMR", UnitPrice: 3.44},
	}

	t.Run("add single new produce", func(t *testing.T) {
		newSingle := make([]Produce, 1)
		newSingle[0] = newProduce[0]
		newProduce, _ := json.Marshal(newSingle)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(newProduce))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("add multiple new produce", func(t *testing.T) {
		newProduce, _ := json.Marshal(newProduce)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(newProduce))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("add new produce with invalid produceCode", func(t *testing.T) {
		newSingle := make([]Produce, 1)
		newSingle[0] = newProduce[0]
		newSingle[0].ProduceCode = "imNotValid"
		newProduce, _ := json.Marshal(newSingle)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(newProduce))
		router.ServeHTTP(w, req)

		expectedBody := `{ "error": "invalid product code detected" }`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expectedBody, w.Body.String())
	})

	t.Run("add new produce no body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/produce", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("produce can be added concurrently", func(t *testing.T) {
		updates := 1000
		store := newStoreHandlers()

		var wg sync.WaitGroup
		wg.Add(updates)

		for i := 0; i < updates; i++ {
			go func(n int) {
				newProduce := make([]Produce, 1)
				newProduce[0] = Produce{
					ProduceCode: "TTTT-TTTT-TTTT-" + fmt.Sprintf("%04d", n),
					Name:        "TestProduce " + fmt.Sprintf("%04d", n),
					UnitPrice:   float64(n),
				}
				store.updateProduceInv(newProduce)
				wg.Done()
			}(i)
		}
		wg.Wait()

		assert.Equal(t, 1004, len(store.getProduceInvSlice()))
	})
}

func TestDeleteProduce(t *testing.T) {
	router := setupRouter()

	t.Run("delete produce by id", func(t *testing.T) {
		w := httptest.NewRecorder()

		// delete produce
		req, _ := http.NewRequest(http.MethodDelete, "/produce/TQ4C-VV6T-75ZX-1RMR", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// verify produce was deleted
		req2, _ := http.NewRequest(http.MethodGet, "/produce", nil)
		router.ServeHTTP(w, req2)

		expectedBody, _ := json.Marshal(initProduce[1:])
		assert.JSONEq(t, string(expectedBody), w.Body.String())
	})

	t.Run("delete produce that does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/produce/this-does-note-xist", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
