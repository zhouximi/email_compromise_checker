package main

import (
	"github.com/zhouximi/email_compromise_checker/data_model"
	"github.com/zhouximi/email_compromise_checker/middleware/cache"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const HIBP_API = "https://haveibeenpwned.com/api/v3/breachedaccount/"

func main() {
	localCache, err := cache.NewLocalCache()
	if err != nil {
		log.Fatalf("Failed to initialize local cache: %v", err)
	}

	remoteCache, err := cache.NewRemoteCache()
	if err != nil {
		log.Fatalf("Failed to initialize remote cache: %v", err)
	}

	multilayerCache := cache.NewMultiLayerCache(localCache, remoteCache)

	r := gin.Default()

	// POST /check
	r.POST("/check", func(c *gin.Context) {
		var req *data_model.EmailRequest
		if err := c.BindJSON(&req); err != nil || req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}

		if !isValidEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		if resp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			c.JSON(http.StatusOK, gin.H{"compromised": true, "breaches": string(body)})
		} else if resp.StatusCode == 404 {
			c.JSON(http.StatusOK, gin.H{"compromised": false})
		} else {
			c.JSON(resp.StatusCode, gin.H{"error": "API error", "status": resp.Status})
		}
	})

	r.Run(":8080")
}
