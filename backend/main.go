package main

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/objectbox/objectbox-go/objectbox"
)

type routeFindingParams struct {
	Products []int `json:"products"`
}

type storeParams struct {
	Name    string     `json:"name"`
	Address string     `json:"address"`
	CSV     string     `json:"csv"`
	Grid    [][]square `json:"grid"`
	Start   point      `json:"start"`
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

func main() {
	box, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		log.Fatal(err)
	}
	defer box.Close()

	userBox := BoxForuser(box)
	productBox := BoxForproduct(box)
	storeBox := BoxForstore(box)

	if len(os.Args) >= 2 {
		if os.Args[1] == "create-admin" && len(os.Args) == 4 {
			passwordHash := sha512.Sum512([]byte(os.Args[3]))
			_, err = userBox.Put(&user{username: os.Args[2], passwordHash: passwordHash[:]})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		} else if len(os.Args) == 2 {
			readCSV(os.Args[1], productBox)
		} else {
			fmt.Println("usage:", os.Args[0], "[create-admin <username> <password>]")
			os.Exit(1)
		}
	}

	_, err = userBox.Query(user_.username.Equals("admin", true)).Remove()
	if err != nil {
		panic(err)
	}
	passwordHash := sha512.Sum512([]byte("admin"))
	_, err = userBox.Insert(&user{username: "admin", passwordHash: passwordHash[:]})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	handler := enableCORS(mux)
	server := &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}

	handleAuth(mux, userBox)

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

	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		products, err := productBox.GetAll()
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
		storeIDString := r.PathValue("store")
		storeID, err := strconv.ParseUint(storeIDString, 10, 64)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}
		store, err := storeBox.Get(storeID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Failed to get store: %v", err)
			return
		}
		json.NewEncoder(w).Encode(decodeGrid(store.Grid, store.Width))
	})

	mux.HandleFunc("POST /stores", func(w http.ResponseWriter, r *http.Request) {
		session := checkSesssion(w, r)
		if session == nil {
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
			grid, start, err = parseCSV(params.CSV)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		id, err := storeBox.Insert(&store{
			Name:    params.Name,
			Address: params.Address,
			Width:   getWidth(grid),
			Grid:    encodeGrid(grid),
			Start:   start,
			Owner:   session.user.id,
		})
		if err != nil {
			log.Printf("Failed to insert store into database: %v", err)
			return
		}
		json.NewEncoder(w).Encode(id)
	})

	mux.HandleFunc("POST /stores/{store}/find-route", func(w http.ResponseWriter, r *http.Request) {
		var params routeFindingParams
		err := newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		storeIDString := r.PathValue("store")
		storeID, err := strconv.ParseUint(storeIDString, 10, 64)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
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
		path := theAlgorithm(decodeGrid(store.Grid, store.Width), store.Start, products)
		end := time.Now()
		log.Println("length:", len(path))
		log.Println("solving time:", end.Sub(begin))

		json.NewEncoder(w).Encode(routeFound{path})
	})

	terminationChan := make(chan os.Signal, 1)
	signal.Notify(terminationChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-terminationChan
		server.Shutdown(context.Background())
	}()

	log.Println("Server started")
	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
