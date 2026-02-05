package main

import (
	"log"

	"github.com/OderoCeasar/system/api/handlers"
	"github.com/OderoCeasar/system/config"
	"github.com/OderoCeasar/system/db/models"
	"github.com/OderoCeasar/system/db/repositories"
	"github.com/OderoCeasar/system/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4/database"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db := database.GetDB()
	if err := db.AutoMigrate(
		&models.User{},
		&models.Package{},
		&models.Payment{},
		&models.Session{},
		&models.RADIUSAccount{},
		&models.RADIUSAccounting{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// initialize repositories
	userRepo := repositories.NewUserRepository(db)
	packageRepo := repositories.NewPackageRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)


	// initialize service
	authService := services.NewAuthService(userRepo, cfg)
	paymentService := services.NewPaymentService(paymentRepo, packageRepo, userRepo, cfg)
	sessionService := services.NewSessionService(sessionRepo, packageRepo, paymentRepo)


	// initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	packageHandler := handlers.NewPackageHandler(packageRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService, sessionService)
	sessionHandler := handlers.NewSessionHandler(sessionService)


	gin.SetMode(cfg.Server.GinMode)
	r := gin.Default()

	routes.SetUpRoutes(r, authService, authHandler, packageHandler, paymentHandler, sessionHandler, cfg.Server.FrontendURL)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
