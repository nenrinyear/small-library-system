package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nenrinyear/small-library-system/go-api/internal/config"
	"github.com/nenrinyear/small-library-system/go-api/internal/server"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(gormmysql.Open(cfg.MySQLDSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to open db: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get sql db: %w", err))
	}
	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		panic(fmt.Errorf("failed to connect db: %w", err))
	}

	srv := server.New(cfg, db)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		if err := srv.Engine().Start(addr); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("server error: %w", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Engine().Shutdown(shutdownCtx)
}
