package main

import (
	"github.com/OderoCeasar/system/api/routes"
	"github.com/OderoCeasar/system/config"
	"github.com/OderoCeasar/system/db"
	"github.com/OderoCeasar/system/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()
	logger := utils.GetLogger()

	cfg := config.Load()
	logger.Info("Configuration loaded successfully")

	if err := db.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}
	logger.Info("Database connected successfully")

	if err := db.AutoMigrate(); err != nil {
		logger.Fatal("Failed to migrate the database: " + err.Error())
	}
	logger.Info("Database migration completed")

	// setup GIN
	if cfg.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// routes
	routes.SetupRoutes(r, cfg)
	logger.Info("Routes configured successfully")

	addr := ":" + cfg.Server.Port
	logger.Info("Server starting on " + addr)

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}
