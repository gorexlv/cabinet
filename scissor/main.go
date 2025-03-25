package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorexlv/cabinet/scissor/pkg/kimi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Article struct {
	ID        uint   `gorm:"primarykey"`
	Title     string `gorm:"size:255"`
	Content   string `gorm:"type:text"`
	Summary   string `gorm:"type:text"`
	Tags      string `gorm:"type:text"`
	SourceURL string `gorm:"size:255"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

func main() {
	// 连接数据库
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/scissor?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(&Article{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建Kimi客户端
	kimiClient, err := kimi.NewClient()
	if err != nil {
		log.Fatal("Failed to create Kimi client:", err)
	}

	// 创建Gin路由
	r := gin.Default()

	// 添加文章接口
	r.POST("/articles", func(c *gin.Context) {
		var article Article
		if err := c.ShouldBindJSON(&article); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 生成摘要
		summary, err := kimiClient.GenerateSummary(article.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate summary: " + err.Error()})
			return
		}
		article.Summary = summary

		// 生成标签
		tags, err := kimiClient.GenerateTags(article.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate tags: " + err.Error()})
			return
		}
		article.Tags = tags

		if err := db.Create(&article).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, article)
	})

	// 获取文章列表
	r.GET("/articles", func(c *gin.Context) {
		var articles []Article
		if err := db.Find(&articles).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, articles)
	})

	// 启动服务器
	r.Run(":8080")
}
