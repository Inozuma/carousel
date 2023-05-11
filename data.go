package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	conn *gorm.DB
}

func NewDatabase(dbpath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(MediaItem{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{
		conn: db,
	}, err
}

type MediaItem struct {
	gorm.Model

	Title   string `json:"title"`
	Path    string `json:"path" gorm:"primaryKey"`
	Library string `json:"library"`
	Type    string `json:"type,omitempty"`
	Episode int    `json:"episode,omitempty"`
}

func (db *Database) SaveItems(items []MediaItem) error {
	return db.conn.Create(items).Error
}

func (db *Database) LoadItems() ([]MediaItem, error) {
	var items []MediaItem

	return items, db.conn.Find(&items).Error
}
