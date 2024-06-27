package handlers

import (
	"backend2/cmd/api/database"
	"backend2/cmd/api/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "secret"

func Register(c *fiber.Ctx) error {
	log.Println("Received a registration request")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	log.Println("Username: ", data["username"], "Password: ", data["password"])

	// Check if the email already exists
	var existingUser models.User
	if err := database.DB.Where("username = ?", data["username"]).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username taken",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	log.Println("Creating User...")
	user := &models.User{
		Username: data["username"],
		Password: string(hashedPassword),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	log.Println("User registered successfully")

	log.Println("Generating JWT token")
	// Generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  strconv.Itoa(int(user.ID)),
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
}

func Login(c *fiber.Ctx) error {
	log.Println("Received a Login request")

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if user exists
	var user models.User
	database.DB.Where("username = ?", data["username"]).First(&user)
	if user.ID == 0 {
		log.Println("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// Compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		log.Println("Invalid Password:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	log.Println("Generating JWT token")
	// Generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  strconv.Itoa(int(user.ID)),
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
}

func GetUser(c *fiber.Ctx) error {
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

	var id = (*claims)["sub"].(string)
	var role = (*claims)["role"].(string)

	// Query user from database using ID
	var user models.User
	database.DB.Where("id = ?", id).First(&user)

	// Go only exports types with a capital letter
	// holy **** thats stupid
	// and then you also have to rename it
	data := struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}{Username: user.Username, Role: role}

	return c.JSON(data)
}

func Logout(c *fiber.Ctx) error {
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
}
