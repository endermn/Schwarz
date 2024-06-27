package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "secret"

func registerAuthEndpoints(app *fiber.App, userBox *userBox) {
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

func runAuth(userBox *userBox) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,Accept,Origin,Access-Control-Request-Method,Access-Control-Request-Headers,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods,Access-Control-Expose-Headers,Access-Control-Max-Age,Access-Control-Allow-Credentials",
		AllowCredentials: true,
	}))

	registerAuthEndpoints(app, userBox)

	app.Listen(":3000")
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
