package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"warehouse-backend/internal/config"
	"warehouse-backend/internal/db"
	"warehouse-backend/internal/handlers"
)

func main() {
	cfg := config.Load()

	gdb, err := db.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("không kết nối được CSDL: %v", err)
	}
	db.Seed(gdb)
	db.SeedDemo(gdb)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	handlers.Register(r, gdb, cfg.JWTSecret)

	log.Printf("server chạy tại :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
