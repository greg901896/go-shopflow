package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/greg901896/go-shopflow/internal/handler"
	"github.com/greg901896/go-shopflow/internal/middleware"
	"github.com/greg901896/go-shopflow/internal/repository"
	"github.com/greg901896/go-shopflow/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("ping failed: %v", err)
	}
	log.Println("connected to postgres")

	productRepo := repository.NewProductRepository(pool)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(userSvc)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db_down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	r.GET("/products", productHandler.List)
	r.GET("/products/:id", productHandler.Get)

	protected := r.Group("/", middleware.JWTAuth(jwtSecret))
	{
		protected.POST("/products", productHandler.Create)
		protected.PUT("/products/:id", productHandler.Update)
	}

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
