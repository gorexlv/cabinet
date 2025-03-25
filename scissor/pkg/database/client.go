package database

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorexlv/cabinet/scissor/internal/config"
	"github.com/gorexlv/cabinet/scissor/pkg/ent"
)

// NewClient 创建一个新的数据库客户端
func NewClient(cfg config.DatabaseConfig) (*ent.Client, error) {
	drv, err := sql.Open(dialect.MySQL, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	client := ent.NewClient(ent.Driver(drv))

	// 运行数据库迁移
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return client, nil
}
