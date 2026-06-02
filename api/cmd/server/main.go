package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nenrinyear/small-library-system/api/internal/config"
	appdb "github.com/nenrinyear/small-library-system/api/internal/db"
	"github.com/nenrinyear/small-library-system/api/internal/server"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	if err := migrateDatabase(cfg); err != nil {
		panic(err)
	}

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

func migrateDatabase(cfg config.Config) error {
	db, err := sql.Open("mysql", cfg.MySQLMigrationDSN())
	if err != nil {
		return fmt.Errorf("failed to open migration db: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect migration db: %w", err)
	}

	if err := appdb.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	return nil
}
