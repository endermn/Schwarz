package main

// import (
// 	"crypto/rand"
// 	"crypto/sha512"
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"
// )

// // JSON automagically encodes []byte in base64, but not [N]byte

// const sessionIDLength = 15

// type sessionID [sessionIDLength]byte

// const sessionDuration = 7 * 24 * time.Hour

// type userCredentials struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// type session struct {
// 	user *user
// }

// var sessions = map[sessionID]session{}

// func checkSesssion(w http.ResponseWriter, r *http.Request) *session {
// 	cookie, _ := r.Cookie("session")
// 	if cookie == nil {
// 		http.Error(w, "no session cookie", http.StatusBadRequest)
// 		return nil
// 	}
// 	sessionIDSlice, err := base64.StdEncoding.DecodeString(cookie.Value)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return nil
// 	}
// 	if len(sessionIDSlice) != sessionIDLength {
// 		http.Error(w, fmt.Sprint("Session ID must be ", sessionIDLength, " bytes"), http.StatusBadRequest)
// 		return nil
// 	}
// 	session, ok := sessions[*(*sessionID)(sessionIDSlice)]
// 	if !ok {
// 		http.Error(w, "Session does not exist (might have expired)", http.StatusBadRequest)
// 		return nil
// 	}
// 	return &session
// }

// func authenticate(userBox *userBox, w http.ResponseWriter, username string, password string) *user {
// 	users, err := userBox.Query(user_.username.Equals(username, true)).Find()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return nil
// 	}

// 	if len(users) == 0 {
// 		http.Error(w, "Invalid name or password", http.StatusBadRequest)
// 		return nil
// 	}
// 	if len(users) > 1 {
// 		log.Println("Multiple users with the same name")
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return nil
// 	}
// 	user := users[0]

// 	if [sha512.Size]byte(user.passwordHash) != sha512.Sum512([]byte(password)) {
// 		http.Error(w, "Invalid name or password", http.StatusBadRequest)
// 		return nil
// 	}
// 	return user
// }

// func logIn(w http.ResponseWriter, user *user) {
// 	var sessionID sessionID
// 	rand.Read(sessionID[:])
// 	sessions[sessionID] = session{user}
// 	go func() {
// 		time.Sleep(sessionDuration)
// 		delete(sessions, sessionID)
// 	}()

// 	encodedSessionID := base64.StdEncoding.EncodeToString(sessionID[:])
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "session",
// 		Value:    encodedSessionID,
// 		Path:     "/",
// 		Expires:  time.Now().Add(sessionDuration),
// 		SameSite: http.SameSiteNoneMode,
// 	})
// }

// func handleAuth(mux *http.ServeMux, userBox *userBox) {
// 	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
// 		var params userCredentials
// 		err := newJSONDecoder(r.Body).Decode(&params)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		if params.Username == "" || params.Password == "" {
// 			http.Error(w, "Username and password must be non-empty", http.StatusBadRequest)
// 			return
// 		}

// 		existing_users, err := userBox.Query(user_.username.Contains(params.Username, true)).Find()
// 		if err != nil {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
// 		if len(existing_users) != 0 {
// 			http.Error(w, "Such user already exists", http.StatusBadRequest)
// 			return
// 		}

// 		passwordHash := sha512.Sum512([]byte(params.Password))
// 		user := &user{username: params.Username, passwordHash: passwordHash[:]}
// 		id, err := userBox.Insert(user)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		log.Println("Created user", id, params.Username, params.Password)

// 		logIn(w, user)
// 	})

// 	mux.HandleFunc("POST /users/login", func(w http.ResponseWriter, r *http.Request) {
// 		var params userCredentials
// 		err := newJSONDecoder(r.Body).Decode(&params)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		user := authenticate(userBox, w, params.Username, params.Password)
// 		if user == nil {
// 			return
// 		}

// 		logIn(w, user)
// 	})

// 	mux.HandleFunc("GET /check-session", func(w http.ResponseWriter, r *http.Request) {
// 		_ = checkSesssion(w, r)
// 	})
// }

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "secret"

func registerAuthHandlers(app *fiber.App, userBox *userBox) {
	app.Post("/api/register", func(c *fiber.Ctx) error {
		log.Println("Received a registration request")
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse request body",
			})
		}

		log.Println("Username: ", data["username"], "Password: ", data["password"])

		// Check if the email already exists
		existingUserCount, err := userBox.Query(user_.username.Equals(data["username"], true)).Count()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
		}
		if existingUserCount > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username taken",
			})
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}

		log.Println("Creating User...")
		user := &user{
			username:     data["username"],
			passwordHash: passwordHash,
		}
		if _, err := userBox.Insert(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		log.Println("User registered successfully")

		log.Println("Generating JWT token")
		// Generate JWT token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  strconv.Itoa(int(user.id)),
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
			"role": "user", // Expires in 24 hours
		})
		token, err := claims.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Println("Error generating token:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		log.Println("Setting cookie")

		// Set JWT token in cookie
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24), // Expires in 24 hours
			HTTPOnly: true,
			Secure:   true,
		}
		c.Cookie(&cookie)

		log.Println("Authentication successful, returning")
		// Authentication successful, return success response
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Login successful",
		})
	})

	app.Post("/api/login", func(c *fiber.Ctx) error {
		log.Println("Received a Login request")

		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse request body",
			})
		}

		// Check if user exists
		users, err := userBox.Query(user_.username.Equals(data["username"], true)).Find()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get user",
			})
		}
		if len(users) == 0 {
			log.Println("User not found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword(users[0].passwordHash, []byte(data["password"]))
		if err != nil {
			log.Println("Invalid Password:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		log.Println("Generating JWT token")
		// Generate JWT token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  strconv.Itoa(int(users[0].id)),
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
			"role": "user", // Expires in 24 hours
		})
		token, err := claims.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Println("Error generating token:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		log.Println("Setting cookie")

		// Set JWT token in cookie
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24), // Expires in 24 hours
			HTTPOnly: true,
			Secure:   true,
		}
		c.Cookie(&cookie)

		log.Println("Authentication successful, returning")
		// Authentication successful, return success response
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Login successful",
		})
	})

	app.Post("/api/logout", func(c *fiber.Ctx) error {
		log.Println("Received a logout request")

		// Clear JWT token by setting an empty value and expired time in the cookie
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour), // Expired 1 hour ago
			HTTPOnly: true,
			Secure:   true,
		}
		c.Cookie(&cookie)

		// Return success response indicating logout was successful
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Logout successful",
		})
	})

	app.Get("/api/user", func(c *fiber.Ctx) error {
		log.Println("Request to get user...")

		// Retrieve JWT token from cookie
		cookie := c.Cookies("jwt")

		log.Println(cookie)

		// Parse JWT token with claims
		token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		// Handle token parsing errors
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Extract claims from token
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to parse claims",
			})
		}

		var idString = (*claims)["sub"].(string)
		var role = (*claims)["role"].(string)

		// Query user from database using ID
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		user, err := userBox.Get(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get user",
			})
		}

		// Go only exports types with a capital letter
		// holy **** thats stupid
		// and then you also have to rename it
		data := struct {
			Username string `json:"username"`
			Role     string `json:"role"`
		}{Username: user.username, Role: role}

		return c.JSON(data)
	})
}

func requireUser(userBox *userBox, w http.ResponseWriter, r *http.Request) *user {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "unauth", http.StatusUnauthorized)
		return nil
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		http.Error(w, "unauth", http.StatusUnauthorized)
		return nil
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		http.Error(w, "failed to parse claims", http.StatusUnauthorized)
		return nil
	}

	var idString = (*claims)["sub"].(string)

	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	user, err := userBox.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return user
}
