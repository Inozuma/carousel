package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Title     string `json:"title"`
	Path      string `json:"path" gorm:"primaryKey"`
	Library   string `json:"library"`
	Type      string `json:"type,omitempty"`
	Episode   int    `json:"episode,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MediaSerie struct {
	gorm.Model
}

func (db *Database) SaveItems(items []MediaItem) error {
	return db.conn.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(items).Error
}

func (db *Database) LoadItems() ([]MediaItem, error) {
	var items []MediaItem

	return items, db.conn.Find(&items).Error
}
