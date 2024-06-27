package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultImage = "https://www.worldplumbing.org/wp-content/uploads/2016/11/white-background-500x500.jpg"

type routeFindingParams struct {
	Products []int `json:"products"`
}

type storeParams struct {
	Name    string `json:"name"`
	Address string `json:"address"`

	// either CSV or Grid and Start
	CSV   string     `json:"csv"`
	Grid  [][]square `json:"grid"`
	Start point      `json:"start"`
}

type point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type routeFound struct {
	Path []point `json:"path"`
}

func newJSONDecoder(r io.Reader) *json.Decoder {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder
}

func getStoreID(w http.ResponseWriter, storeIDString string, defaultStoreID uint64) uint64 {
	storeID, err := strconv.ParseUint(storeIDString, 10, 64)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return 0
	}
	if storeID == 0 {
		storeID = defaultStoreID
	}
	return storeID
}

func registerMainEndpoints(mux *http.ServeMux, userBox *userBox, productBox *productBox, storeBox *storeBox, defaultStoreID uint64) {
	mux.HandleFunc("GET /categories", func(w http.ResponseWriter, r *http.Request) {
		categories := set[string]{}
		products, err := productBox.GetAll()
		if err != nil {
			log.Printf("Failed to get products: %v", err)
		}

		for _, product := range products {
			categories.insert(product.Category)
		}

		json.NewEncoder(w).Encode(categories.toArray())
	})

	mux.HandleFunc("POST /users/{user}/products", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("user")

		users, err := userBox.Query(user_.username.Equals(username, true)).Find()
		if err != nil {
			log.Printf("Cannot get user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if len(users) == 0 {
			http.Error(w, "User does not exist", http.StatusBadRequest)
			return
		}

		var params product
		err = newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if params.ImageURL == "" {
			params.ImageURL = defaultImage
		}
		params.Owner = users[0]

		_, err = productBox.Insert(&params)
		if err != nil {
			log.Printf("Failed to insert into database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /users/{user}/products", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("user")

		users, err := userBox.Query(user_.username.Equals(username, true)).Find()
		if err != nil {
			log.Printf("Cannot get user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if len(users) == 0 {
			http.Error(w, "User does not exist", http.StatusBadRequest)
			return
		}

		products, err := productBox.Query(product_.Owner.Equals(users[0].id)).Find()
		if err != nil {
			log.Println("Failed to get products:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(products)
	})

	mux.HandleFunc("GET /stores", func(w http.ResponseWriter, r *http.Request) {
		stores, err := storeBox.GetAll()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Failed to get stores: %v", err)
			return
		}
		json.NewEncoder(w).Encode(stores)
	})

	mux.HandleFunc("GET /stores/{store}/layout", func(w http.ResponseWriter, r *http.Request) {
		storeID := getStoreID(w, r.PathValue("store"), defaultStoreID)
		if storeID == 0 {
			return
		}
		store, err := storeBox.Get(storeID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Failed to get store: %v", err)
			return
		}
		json.NewEncoder(w).Encode(decodeGrid(store.grid, store.width))
	})

	mux.HandleFunc("POST /stores", func(w http.ResponseWriter, r *http.Request) {
		user := requireUser(userBox, w, r)
		if user == nil {
			return
		}

		var params storeParams
		err := newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		grid := params.Grid
		start := params.Start
		if params.CSV != "" {
			var err error
			grid, start, err = parseStoreCSV(strings.NewReader(params.CSV))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		id, err := storeBox.Insert(&store{
			Name:    params.Name,
			Address: params.Address,
			width:   getWidth(grid),
			grid:    encodeGrid(grid),
			start:   start,
			owner:   user.id,
		})
		if err != nil {
			log.Printf("Failed to insert store into database: %v", err)
			return
		}
		json.NewEncoder(w).Encode(id)
	})

	mux.HandleFunc("PUT /stores/{store}", func(w http.ResponseWriter, r *http.Request) {
		user := requireUser(userBox, w, r)
		if user == nil {
			return
		}

		storeID := getStoreID(w, r.PathValue("store"), defaultStoreID)
		if storeID == 0 {
			return
		}

		activeStore, err := storeBox.Get(storeID)
		if err != nil {
			log.Printf("Failed to get store with id: %v,: %v", storeID, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.id != activeStore.owner {
			log.Printf("User with id: %v is not the owner of this store", user.id)
			http.Error(w, "Invalid user", http.StatusBadRequest)
			return
		}

		var params storeParams
		err = newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		grid := params.Grid
		start := params.Start
		if params.CSV != "" {
			var err error
			grid, start, err = parseStoreCSV(strings.NewReader(params.CSV))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = storeBox.Update(&store{
			ID:      storeID,
			Name:    params.Name,
			Address: params.Address,
			width:   getWidth(grid),
			grid:    encodeGrid(grid),
			start:   start,
		})
		if err != nil {
			log.Printf("Failed to update store: %v", err)
			return
		}
	})

	mux.HandleFunc("POST /stores/{store}/find-route", func(w http.ResponseWriter, r *http.Request) {
		var params routeFindingParams
		err := newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		storeID := getStoreID(w, r.PathValue("store"), defaultStoreID)
		if storeID == 0 {
			return
		}

		store, err := storeBox.Get(storeID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Failed to get store: %v", err)
			return
		}

		products := set[int]{}
		for _, productID := range params.Products {
			products[productID] = struct{}{}
		}
		begin := time.Now()
		path := theAlgorithm(decodeGrid(store.grid, store.width), store.start, products)
		log.Println("solving time:", time.Since(begin))
		log.Println("length:", len(path))

		json.NewEncoder(w).Encode(routeFound{path})
	})
}
