package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

// JSON automagically encodes []byte in base64, but not [N]byte

const sessionIDLength = 15

type sessionID [sessionIDLength]byte

const sessionDuration = 7 * 24 * time.Hour

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func handleAuth(mux *http.ServeMux, userBox *userBox) {
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
}
