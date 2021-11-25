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
		req, err := http.NewRequest(http.MethodGet, "/produce", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)

		expectedRes, err := json.Marshal(initProduce)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedRes), w.Body.String())
	})

	t.Run("get produce by ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/produce/TQ4C-VV6T-75ZX-1RMR", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		expectedRes, err := json.Marshal(initProduce[0])
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, string(expectedRes), rr.Body.String())
	})

	t.Run("get produce by ID 404", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/produce/kalhefid", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
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
		newProduce, err := json.Marshal(newSingle)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(newProduce))
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("add multiple new produce", func(t *testing.T) {
		newProduce, err := json.Marshal(newProduce)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(newProduce))
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("add new produce with invalid produceCode", func(t *testing.T) {
		newSingle := make([]Produce, 1)
		newSingle[0] = newProduce[0]
		newSingle[0].ProduceCode = "imNotValid"
		newProduce, _ := json.Marshal(newSingle)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader(newProduce))
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		expectedBody := `{ "error": "invalid product code detected" }`

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("add new produce nil body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/produce", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("content-type", "application/json")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("add new produce empty body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewReader([]byte("")))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("content-type", "application/json")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("produce can be added concurrently", func(t *testing.T) {
		updates := 9999
		store, err := newStoreHandlers()
		if err != nil {
			t.Fatal(err)
		}

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

		assert.Equal(t, updates+4, len(store.getProduceInvSlice()))
	})
}

func TestDeleteProduce(t *testing.T) {
	router := setupRouter()

	t.Run("delete produce by id", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// delete produce
		req, err := http.NewRequest(http.MethodDelete, "/produce/TQ4C-VV6T-75ZX-1RMR", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// verify produce was deleted
		req2, err := http.NewRequest(http.MethodGet, "/produce", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req2)

		expectedBody, _ := json.Marshal(initProduce[1:])
		assert.JSONEq(t, string(expectedBody), rr.Body.String())
	})

	t.Run("delete produce that does not exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/produce/this-does-note-xist", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
