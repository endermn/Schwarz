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
	mux := http.NewServeMux()
	handler := enableCORS(mux)
	server := &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}

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

	mux.HandleFunc("POST /find-route", func(w http.ResponseWriter, r *http.Request) {
		var params routeFindingParams
		err := newJSONDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		products := set[int]{}
		for _, productID := range params.Products {
			products[productID] = struct{}{}
		}
		path := solve(products)

		json.NewEncoder(w).Encode(routeFound{path})
	})

	mux.HandleFunc("GET /store-layout", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(grid)
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
