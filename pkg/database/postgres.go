package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/rohanchauhan02/internal-transfer/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres database connection and initialization
type Postgres interface {
	InitClient(ctx context.Context) (*gorm.DB, error)
}

type database struct {
	conf config.ImmutableConfigs
}

var (
	once sync.Once
	db   *gorm.DB
	err  error
)

// NewPostgres creates a new Postgres instance with the provided configuration.
func NewPostgres(conf config.ImmutableConfigs) Postgres {
	return &database{
		conf: conf,
	}
}

func (d *database) InitClient(ctx context.Context) (*gorm.DB, error) {
	once.Do(func() {
		log.Info("Initializing PostgreSQL connection...")

		dbConfig := d.conf.GetDBConf()

		connectionString := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=%s",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Host,
			dbConfig.SSLMode,
		)

		// Retry mechanism for transient failures
		maxRetries := 3
		for i := range maxRetries {
			db, err = gorm.Open(postgres.New(postgres.Config{
				DSN:                  connectionString,
				PreferSimpleProtocol: true,
			}), &gorm.Config{
				DisableAutomaticPing: false,
				PrepareStmt:          true,
			})

			if err == nil {
				break
			}

			log.Errorf("Postgres connection attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Errorf("PostgreSQL connection failed after retries: %v", err)
			return
		}

		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			log.Errorf("Failed to get underlying sql.DB: %v", sqlErr)
			err = sqlErr
			return
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		log.Info("Successfully connected to PostgreSQL!")
	})

	return db, err
}
