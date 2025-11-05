package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GibGyb/todo-project/auth"
	"github.com/GibGyb/todo-project/todo"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	buildCommit = "dev"
	buildTime = time.Now().String()
)
  
func main() {
	// Liveness Probe 
	_,err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Println("please consider environment variables: %s", err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	db.AutoMigrate(&todo.Todo{})

	router := gin.Default()
	// Allow Cors
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowHeaders = []string{
		"Origin", 
		"Authorization",
		"Transaction-Id",
	}
	router.Use(cors.New(config))
	
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Readiness Probe
	router.GET("/healthz", func (c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/limitz", limitHandler)

	router.GET("/x", func (c *gin.Context) {
		c.JSON(200, gin.H{
			"buildCommit": buildCommit, 
			"buildTime": buildTime,
		})
	})

	router.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := router.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	handler := todo.NewHandler(db)
	protected.POST("/todos", handler.NewTask)
	protected.GET("/todos", handler.List)
	protected.DELETE("/todos/:id", handler.Remove)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20, 
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}


// Rate Limiter
var limiter = rate.NewLimiter(5, 5)

func limitHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		return
	}
	c.JSON(200, gin.H{"message": "pong"})
	} 