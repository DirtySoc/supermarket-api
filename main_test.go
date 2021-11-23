package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetProduce(t *testing.T) {
	t.Run("GET all produce", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/produce", nil)
		response := httptest.NewRecorder()

		handleProduce(response, request)

		var got []Produce

		json.Unmarshal(response.Body.Bytes(), &got)
		want := ProduceDB

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
