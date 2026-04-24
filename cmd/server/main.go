package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greg901896/go-shopflow/internal/handler"
	"github.com/greg901896/go-shopflow/internal/repository"
	"github.com/greg901896/go-shopflow/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgres://shopflow:shopflow_dev@localhost:5433/shopflow"

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

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db_down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/products", productHandler.Create)
	r.GET("/products", productHandler.List)
	r.GET("/products/:id", productHandler.Get)
	r.PUT("/products/:id", productHandler.Update)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
