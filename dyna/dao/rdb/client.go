package rdb

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	db *gorm.DB
}

func newDBConnection(connection string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             1.0 * time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info,       // Log level
				IgnoreRecordNotFoundError: true,              // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,              // Disable color
			},
		),
		PrepareStmt: true,
	}

	db, err := gorm.Open(postgres.Open(connection), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewClient(connection string) (*Client, error) {
	db, err := newDBConnection(connection)
	if err != nil {
		return nil, err
	}

	c := &Client{
		db: db,
	}
	c.InitializeModels()
	return c, nil
}

func (c *Client) InitializeModels() error {
	return c.db.AutoMigrate(
		Message{},
	)
}
