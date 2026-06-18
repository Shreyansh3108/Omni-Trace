package main

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	db "omnitrace/db/sqlc"
)

// Notice the new `validate` tags! This is what the recruiter asked for.
type CreateUserRequest struct {
	Email        string `json:"email" validate:"required,email"`
	PasswordHash string `json:"password_hash" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
}

// Initialize the validator globally
var validate = validator.New()

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting OmniTrace Initialization...")

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, checking system variables")
	}

	dbURL := os.Getenv("DB_SOURCE")
	if dbURL == "" {
		logger.Fatal("CRITICAL: DB_SOURCE environment variable is not set")
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to the database", zap.Error(err))
	}
	defer dbPool.Close()
	logger.Info("Successfully connected to the Neon Cloud Database")

	queries := db.New(dbPool)

	app := fiber.New(fiber.Config{
		AppName: "OmniTrace API v1.0",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		logger.Info("Health check endpoint hit")
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "OmniTrace system is fully operational",
		})
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body format"})
		}

		// THE RECRUITER REQUIREMENT: Run the validation check
		if err := validate.Struct(&req); err != nil {
			logger.Warn("Validation failed", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
		}

		user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
			Email:        req.Email,
			PasswordHash: req.PasswordHash,
			FullName:     req.FullName,
		})
		
		if err != nil {
			logger.Error("Database insertion failed", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user. Email might already exist."})
		}

		logger.Info("User created successfully", zap.String("email", user.Email))
		return c.Status(fiber.StatusCreated).JSON(user)
	})

	app.Get("/users/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")

		user, err := queries.GetUserByEmail(context.Background(), email)
		if err != nil {
			logger.Error("User lookup failed", zap.String("email", email), zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		logger.Info("User retrieved successfully", zap.String("email", user.Email))
		return c.Status(fiber.StatusOK).JSON(user)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logger.Info("Starting server", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}