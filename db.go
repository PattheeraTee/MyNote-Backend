package main

import (
	"fmt"
	// "log"
	"miw/entities"
	// "os"
	// "time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

func setupDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold: time.Second,
	// 		LogLevel:      logger.Info,
	// 		Colorful:      true,
	// 	},
	// )

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&entities.User{},
		&entities.Note{},
		&entities.Reminder{},
		&entities.Tag{},
		&entities.ShareNote{},
		&entities.Event{},
		&entities.NoteTag{},
	)
	if err != nil {
		panic("failed to migrate database schema")
	}

	return db
}

