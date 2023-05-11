package main

type Scraper interface {
	ScrapMetadata(q string) map[string]interface{}
}
