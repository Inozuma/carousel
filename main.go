package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"carousel/parser"
	"carousel/scraper"

	"github.com/h2non/filetype"
)

func main() {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}
	cfg := &Config{}
	if err = json.Unmarshal(configData, cfg); err != nil {
		log.Fatalf("failed to decode config: %s", err)
	}

	db, err := NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("cannot open database at %q: %s", cfg.Database, err)
	}

	// File crawler
	var count int
	var items []MediaItem
	for _, lb := range cfg.Libraries {
		log.Println("Library", lb.Name)
		var lparser parser.Parser
		var exist bool
		if lb.Parser != "" {
			lparser, exist = parser.Parsers[lb.Parser]
			if !exist {
				log.Fatalf("parser %s does not exist", lb.Parser)
			}
			log.Println("Parser", lb.Parser)
		}

		var scr scraper.Scraper
		if lb.Scraper != nil {
			fn, exist := scraper.Scrapers[lb.Scraper.Name]
			if !exist {
				log.Fatalf("scraper %s does not exist", lb.Scraper)
			}
			scr = fn(lb.Scraper.Options)
			log.Println("Scraper", lb.Scraper)
		}

		walkMedia(lb.Path, func(filename, path string) error {
			count++

			var item MediaItem
			if lparser != nil {
				tm, err := lparser.Parse(path)
				if err == nil {
					item = MediaItem{
						Title:   tm.Title,
						Path:    tm.Path,
						Type:    tm.Type,
						Episode: tm.Episode,
						Library: lb.Name,
					}
				}
			} else {
				item = MediaItem{
					Path:    path,
					Library: lb.Name,
				}
			}
			items = append(items, item)

			return nil
			if scr != nil && item.Title != "" {
				res, err := scr.Scrape(item.Title)
				if err != nil {
					log.Printf("failed to scrape result for %q: %s", item.Title, err)
				} else {
					fmt.Println(item.Title, "=>", res.Title)
				}
			}
			return nil
		})
	}

	// TODO: metadata

	if err := db.SaveItems(items); err != nil {
		log.Fatalf("failed to save media items: %s", err)
	}

	log.Printf("found %d media items", count)

	// DRAW
	app := newApplication(db)
	if err = app.run(); err != nil {
		log.Fatalf("failed to run carousel: %s", err)
	}
}

func walkMedia(path string, process func(filename, dir string) error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("failed to read directory: %s", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			walkMedia(filepath.Join(path, entry.Name()), process)
			continue
		}

		fileext := filepath.Ext(entry.Name())
		if fileext == "" {
			continue
		}

		loc := filepath.Join(path, entry.Name())

		fd, err := os.Open(loc)
		if err != nil {
			log.Printf("failed to open file: %s", err)
			continue
		}

		header := make([]byte, 261)
		fd.Read(header)
		fd.Close()

		if filetype.IsVideo(header) {
			process(entry.Name(), loc)
		}
	}
}
