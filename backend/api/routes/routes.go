package routes


import (
	"github.com/OderoCeasar/system/api/handlers"
	"github.com/OderoCeasar/system/api/middleware"
	"github.com/OderoCeasar/system/config"
	"github.com/OderoCeasar/system/db"
	"github.com/OderoCeasar/system/db/repositories"
	"github.com/OderoCeasar/system/mpesa"
	"github.com/OderoCeasar/system/services"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.GetDB())
	packageRepo := repositories.NewPackageRepository(db.GetDB())
	paymentRepo := repositories.NewPaymentRepository(db.GetDB())
	sessionRepo := repositories.NewSessionRepository(db.GetDB())

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	sessionService := services.NewSessionService(packageRepo, sessionRepo, paymentRepo)
	radiusService := services.NewRADIUSService(userRepo, sessionRepo, cfg)
	mpesaService := mpesa.NewService(&cfg.Mpesa, paymentRepo, packageRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	packageHandler := handlers.NewPackageHandler(packageRepo)
	paymentHandler := handlers.NewPaymentHandler(mpesaService, sessionService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	adminHandler := handlers.NewAdminHandler(sessionService)
	radiusHandler := handlers.NewRADIUSHandler(radiusService)


	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.Server.FrontendURL))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "WiFi Billing System",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		//public routes
		
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/quick-register", authHandler.QuickRegister)
		}

		packages := v1.Group("/packages")
		{
			packages.GET("", packageHandler.ListPackage)
			packages.GET("/:id", packageHandler.GetPackage)
		}


		payments := v1.Group("/payments")
		{
			payments.POST("/callback", paymentHandler.HandleCallback)
		}


		radius := v1.Group("/radius")
		{
			radius.POST("/auth", radiusHandler.Authenticate)
			radius.POST("/accounting/start", radiusHandler.AccountingStart)
			radius.POST("/accounting/update", radiusHandler.AccountingUpdate)
			radius.POST("/accounting/stop", radiusHandler.AccountingStop)
		}

	// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
		
			protected.GET("/me", authHandler.Me)

			
			protected.POST("/payments/initiate", paymentHandler.InitiatePayment)
			protected.GET("/payments/:id/status", paymentHandler.CheckPaymentStatus)

			sessions := protected.Group("/sessions")
			{
				sessions.POST("", sessionHandler.CreateSession)
				sessions.GET("/active", sessionHandler.GetActiveSession)
				sessions.GET("/history", sessionHandler.ListUserSessions)
				sessions.GET("/:id", sessionHandler.GetSession)
				sessions.GET("/:id/stats", sessionHandler.GetSessionStats)
				sessions.POST("/:id/disconnect", sessionHandler.DisconnectSession)
			}

	// admin routes		
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				
				admin.POST("/packages", packageHandler.CreatePackage)
				admin.PUT("/packages/:id", packageHandler.UpdatePackage)
				admin.DELETE("/packages/:id", packageHandler.DeletePackage)

				adminSessions := admin.Group("/sessions")
				{
					adminSessions.GET("/active", adminHandler.ListActiveSessions)
				}
			}
		}
	}
}
