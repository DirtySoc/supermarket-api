// Package main provides the main implementation of the supermarket api
package main

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
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
func (h *storeHandlers) getProduce(c *gin.Context) {
	c.JSON(http.StatusOK, h.getProduceInvSlice())
}

// handles GET /produce/:id requests
func (h *storeHandlers) getProduceByID(c *gin.Context) {
	id := c.Param("id")

	h.produceInvMu.Lock()
	item, found := h.produceInv[id]
	h.produceInvMu.Unlock()

	if !found {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, item)
}

// handles POST /produce requests with body containing
// JSON array of new/updated produce entries
func (h *storeHandlers) postProduce(c *gin.Context) {
	var newProduce []Produce

	if err := c.BindJSON(&newProduce); err != nil {
		return
	}

	if err := h.validateProduceIDs(newProduce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err.Error())
		return
	}

	h.updateProduceInv(newProduce)
	c.Status(http.StatusCreated)
}

// handles DELETE /produce/:id requests
func (h *storeHandlers) deleteProduceByID(c *gin.Context) {
	id := c.Param("id")

	h.produceInvMu.Lock()
	defer h.produceInvMu.Unlock()

	delete(h.produceInv, id)
}

// confiugres and initializes the GIN router
func setupRouter() *gin.Engine {
	store := newStoreHandlers()
	router := gin.Default()
	router.GET("/produce", store.getProduce)
	router.GET("/produce/:id", store.getProduceByID)
	router.POST("/produce", store.postProduce)
	router.DELETE("/produce/:id", store.deleteProduceByID)
	return router
}

// starts http server
func main() {
	router := setupRouter()
	router.Run(":6620")
}
