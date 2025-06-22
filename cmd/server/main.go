package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vineet12344/targeting-engine/internal/campaign"
	"github.com/vineet12344/targeting-engine/middleware"
	"github.com/vineet12344/targeting-engine/pkg/db"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// Load ENV Variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// fmt.Println("Hello World")
	// Connectign to Database
	db.ConnectDB()
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get DB Pool %v", err)
	}

	// Automigrate Database
	err = db.DB.AutoMigrate(&campaign.Campaign{}, &campaign.TargetingRule{})
	if err != nil {
		log.Fatalf("❌ Failed to AutoMigrate DB %v", err)
	}

	log.Println(" ✅ Automigration of Database Sucessfull!")

	// seed database
	// if err = campaign.SeedCampaings(); err != nil {
	// 	log.Printf("Cannot seed database %v", err)
	// }
	// log.Println(" ✅ Seeding of Database Sucessfull Successfull !")

	campaignService := campaign.NewCampaignService()

	// Load initial cache from DB
	if err := campaign.LoadToCache(campaignService); err != nil {
		log.Fatalf("❌ Failed to load initial campaign cache: %v", err)
	}

	log.Println("✅ Campaign cache loaded into memory")

	stopChan := make(chan struct{})
	campaign.StartAutoRefresh(campaignService, 1*time.Minute, stopChan)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"Message": "Server up and running"}) })
	router.Use(middleware.PrometheusMiddleware())

	router.GET("/metrics", middleware.MetricsHandler())

	campaign.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Running Server in Goroutine
	go func() {
		log.Printf("🚀 Server running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	} else {
		log.Println("🛑 DB connection closed.")
	}

	// Stop Auto-refresh
	close(stopChan)
	log.Println("🛑 Ticker Stopped successfully")

	if err := sqlDB.Close(); err != nil {
		log.Printf("❌ Error closing DB connection: %v", err)
	}

	log.Println("🛑 Server shutdown cleanly")
}
