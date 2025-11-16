package main

import (
	"context"
	"log"
	"os"

	httpServer "github.com/f4ke-n0name/avito/internal/app/http"
	"github.com/f4ke-n0name/avito/internal/domain/services"
	"github.com/f4ke-n0name/avito/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	userRepo := db.NewUserRepositoryPG(database)
	teamRepo := db.NewTeamRepositoryPG(database)
	prRepo := db.NewPRRepositoryPG(database)

	withTx := func(ctx context.Context, fn func(ctx context.Context) error) error {
		return database.WithTx(ctx, fn)
	}

	userSvc := services.NewUserService(userRepo)
	teamSvc := services.NewTeamService(teamRepo, userRepo)
	prSvc := services.NewPRService(userRepo, teamRepo, prRepo, withTx)

	server := httpServer.NewServer(prSvc, userSvc, teamSvc)
	r := gin.Default()
	server.RegisterRoutes(r)

	log.Println("Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
