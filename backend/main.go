package main

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
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

// JSON automagically encodes []byte in base64, but not [N]byte

const sessionIDLength = 15

type sessionID [sessionIDLength]byte

const sessionDuration = 7 * 24 * time.Hour

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

type session struct {
	user *user
}

var sessions = map[sessionID]session{}

func checkSesssion(w http.ResponseWriter, r *http.Request) *session {
	cookie, _ := r.Cookie("session")
	if cookie == nil {
		http.Error(w, "no session cookie", http.StatusBadRequest)
		return nil
	}
	sessionIDSlice, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if len(sessionIDSlice) != sessionIDLength {
		http.Error(w, fmt.Sprint("Session ID must be ", sessionIDLength, " bytes"), http.StatusBadRequest)
		return nil
	}
	session, ok := sessions[*(*sessionID)(sessionIDSlice)]
	if !ok {
		http.Error(w, "Session does not exist (might have expired)", http.StatusBadRequest)
		return nil
	}
	return &session
}

func authenticate(userBox *userBox, w http.ResponseWriter, username string, password string) *user {
	users, err := userBox.Query(user_.username.Equals(username, true)).Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	if len(users) == 0 {
		http.Error(w, "Invalid name or password", http.StatusBadRequest)
		return nil
	}
	if len(users) > 1 {
		log.Println("Multiple users with the same name")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil
	}
	user := users[0]

	if [sha512.Size]byte(user.passwordHash) != sha512.Sum512([]byte(password)) {
		http.Error(w, "Invalid name or password", http.StatusBadRequest)
		return nil
	}
	return user
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

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var params userCredentials
		err := newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if params.Username == "" || params.Password == "" {
			http.Error(w, "Username and password must be non-empty", http.StatusBadRequest)
			return
		}

		existing_users, err := userBox.Query(user_.username.Contains(params.Username, true)).Find()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if len(existing_users) != 0 {
			http.Error(w, "Such user already exists", http.StatusBadRequest)
			return
		}

		passwordHash := sha512.Sum512([]byte(params.Password))
		user := &user{username: params.Username, passwordHash: passwordHash[:]}
		id, err := userBox.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Created user", id, params.Username, params.Password)
	})

	mux.HandleFunc("POST /users/login", func(w http.ResponseWriter, r *http.Request) {
		var params userCredentials
		err := newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := authenticate(userBox, w, params.Username, params.Password)
		if user == nil {
			return
		}

		var sessionID sessionID
		rand.Read(sessionID[:])
		sessions[sessionID] = session{user}
		go func() {
			time.Sleep(sessionDuration)
			delete(sessions, sessionID)
		}()

		encodedSessionID := base64.StdEncoding.EncodeToString(sessionID[:])
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    encodedSessionID,
			Path:     "/",
			Expires:  time.Now().Add(sessionDuration),
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
	})

	mux.HandleFunc("GET /check-session", func(w http.ResponseWriter, r *http.Request) {
		_ = checkSesssion(w, r)
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

	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		products, err := productBox.GetAll()
		if err != nil {
			log.Println("Failed to get products:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)
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
