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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/objectbox/objectbox-go/objectbox"
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

func main() {
	box, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		log.Fatal(err)
	}
	defer box.Close()

	userBox := BoxForuser(box)
	productBox := BoxForproduct(box)
	storeBox := BoxForstore(box)

	var defaultStoreID uint64

	_, err = userBox.Query(user_.username.Equals("admin", true)).Remove()
	if err != nil {
		panic(err)
	}
	passwordHash := sha512.Sum512([]byte("admin"))
	defaultUser := &user{username: "admin", passwordHash: passwordHash[:]}
	_, err = userBox.Insert(defaultUser)
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 4 && os.Args[1] == "create-admin" {
		passwordHash := sha512.Sum512([]byte(os.Args[3]))
		_, err = userBox.Put(&user{username: os.Args[2], passwordHash: passwordHash[:]})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	} else if len(os.Args) == 3 {
		records := readRecordsFromCSV(os.Args[1])
		loadProducts(records, productBox, defaultUser)

		storeText, err := os.ReadFile(os.Args[2])
		if err != nil {
			log.Printf("Failed to open data: %v", err)
			os.Exit(1)
		}
		grid, start, err := parseCSV(string(storeText))
		if err != nil {
			log.Printf("Failed to parse layout CSV: %v", err)
			os.Exit(1)
		}
		_, err = storeBox.Query(store_.Name.Equals("default", true)).Remove()
		if err != nil {
			panic(err)
		}
		defaultStoreID, err = storeBox.Insert(&store{
			Name:  "default",
			Width: getWidth(grid),
			Grid:  encodeGrid(grid),
			Start: start,
			Owner: 0,
		})
		if err != nil {
			log.Printf("Failed to insert store into database: %v", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("usage:", os.Args[0], "(create-admin <username> <password> | <products.csv> <store.csv>)")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	handler := enableCORS(mux)
	server := &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}

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
		json.NewEncoder(w).Encode(decodeGrid(store.Grid, store.Width))
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
			Owner:   user.id,
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

		if user.id != activeStore.Owner {
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
			grid, start, err = parseCSV(params.CSV)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = storeBox.Update(&store{
			ID:      storeID,
			Name:    params.Name,
			Address: params.Address,
			Width:   getWidth(grid),
			Grid:    encodeGrid(grid),
			Start:   start,
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
		path := theAlgorithm(decodeGrid(store.Grid, store.Width), store.Start, products)
		log.Println("solving time:", time.Since(begin))
		log.Println("length:", len(path))

		json.NewEncoder(w).Encode(routeFound{path})
	})

	terminationChan := make(chan os.Signal, 1)
	signal.Notify(terminationChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-terminationChan
		server.Shutdown(context.Background())
	}()

	go func() {
		app := fiber.New()

		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:5173",
			AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders:     "Content-Type,Authorization,Accept,Origin,Access-Control-Request-Method,Access-Control-Request-Headers,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods,Access-Control-Expose-Headers,Access-Control-Max-Age,Access-Control-Allow-Credentials",
			AllowCredentials: true,
		}))

		registerAuthHandlers(app, userBox)

		app.Listen(":3000")
	}()

	log.Println("Server started")
	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
