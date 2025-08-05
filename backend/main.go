package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zhouximi/email_compromise_checker/data_model"
	"github.com/zhouximi/email_compromise_checker/handler"
	"github.com/zhouximi/email_compromise_checker/middleware/cache"
	"github.com/zhouximi/email_compromise_checker/middleware/db"
)

func main() {
	initService()

	r := gin.Default()

	// POST /check
	r.POST("/check", func(c *gin.Context) {
		var req data_model.EmailCheckAPIRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}

		emailInfo, err := handler.GlobalHandler.IsEmailCompromised(req.Email)
		if err != nil {
			log.Printf("[/check] Error checking email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if emailInfo == nil {
			c.JSON(http.StatusOK, gin.H{
				"compromised": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"compromised": emailInfo.Compromised,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}

func initService() {
	localCache, err := cache.NewLocalCache()
	if err != nil {
		log.Fatalf("Failed to initialize local cache: %v", err)
	}
	remoteCache, err := cache.NewRemoteCache()
	if err != nil {
		log.Fatalf("Failed to initialize remote cache: %v", err)
	}
	multilayerCache := cache.NewMultiLayerCache(localCache, remoteCache)

	mysqlDB, err := db.NewMySQLQuerier()
	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}

	handler.GlobalHandler = handler.NewEmailCheckHandler(multilayerCache, mysqlDB)
}
