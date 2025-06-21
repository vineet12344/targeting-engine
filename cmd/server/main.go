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
	"github.com/vineet12344/targeting-engine/pkg/db"
)

func main() {

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

	// if err := campaign.SeedCampaings(); err != nil {
	// 	log.Fatalf("❌ Failed to seed campaigns: %v", err)
	// }

	// log.Println(" ✅ Database seeded Sucessfully ! ")

	campaignService := campaign.NewCampaignService()

	campaigns, err := campaignService.GetActiveCampaings()

	if err != nil {
		log.Fatalf("❌ Failed to fetch active campaigns: %v", err)
	}

	log.Printf("✅ Loaded %d active campaigns\n", len(campaigns))

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"Message": "Server up and running"}) })

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
		log.Println("✅ DB connection closed.")
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("❌ Error closing DB connection: %v", err)
	}

	log.Println("✅ Server shutdown cleanly")
}
