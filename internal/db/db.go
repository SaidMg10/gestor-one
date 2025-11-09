package db

import (
	"fmt"
	"log"
	"time"

	"github.com/SaidMg10/gestor-one/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg config.DBConfig) error {
	dsn := cfg.DSN
	if dsn == "" {
		// Fallback: armar el DSN desde los campos
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.Name,
			cfg.Port,
		)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Error connecting to DB: %v", err)
	}

	// Configurar el pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Error connecting to DB: %v", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTime)
	} else {
		sqlDB.SetConnMaxIdleTime(15 * time.Minute) // valor por defecto
	}

	DB = db

	log.Println("✅ Connected to DB")
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("X Error closing DB: %v", err)
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		log.Printf("X Error closing DB: %v", err)
		return err
	}
	log.Println("✅ Closed DB")
	return nil
}
