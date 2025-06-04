package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("加载: %v", err)
	}

	port := cfg.Server.Port
	fmt.Printf("后端端口:%d \n", port)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 随机数API端点
	r.GET("/api/time", func(c *gin.Context) {
		//返回JSON响应
		c.JSON(http.StatusOK, gin.H{
			"timestamp":    time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/api/port", func(c *gin.Context) {
		// 返回JSON响应
		c.JSON(http.StatusOK, gin.H{
			"port": port,
		})
	})

	// 启动服务器
	r.Run(":8080")
}
