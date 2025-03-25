package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorexlv/cabinet/scissor/internal/config"
	"github.com/gorexlv/cabinet/scissor/internal/handler"
	"github.com/gorexlv/cabinet/scissor/internal/repository"
	"github.com/gorexlv/cabinet/scissor/internal/service"
	"github.com/gorexlv/cabinet/scissor/pkg/database"
	"github.com/gorexlv/cabinet/scissor/pkg/wechat"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewClient(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 初始化微信客户端
	wxClient := wechat.NewClient(cfg.Wechat.AppID, cfg.Wechat.AppSecret)

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo, cfg.JWT.Secret, wxClient)
	articleService := service.NewArticleService(articleRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	articleHandler := handler.NewArticleHandler(articleService)

	// 创建 WebService
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	// 注册路由
	userHandler.Register(ws)
	articleHandler.Register(ws)

	// 创建 WebService 容器
	container := restful.NewContainer()
	container.Add(ws)

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: container,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
