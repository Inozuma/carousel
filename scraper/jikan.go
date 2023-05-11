package scraper

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/darenliang/jikan-go"
)

func init() {
	Scrapers["jikan"] = func(opts map[string]string) Scraper {
		return NewJikanScraperWithOpts(opts)
	}
}

type JikanScraper struct {
	ticker *time.Ticker
	opts   map[string]string
}

func NewJikanScraper() *JikanScraper {
	return &JikanScraper{
		ticker: time.NewTicker(time.Second),
		opts:   make(map[string]string),
	}
}

func NewJikanScraperWithOpts(opts map[string]string) *JikanScraper {
	return &JikanScraper{
		ticker: time.NewTicker(time.Second),
		opts:   opts,
	}
}

func (scraper *JikanScraper) Scrape(title string) (*ScrapeResult, error) {
	<-scraper.ticker.C

	q := make(url.Values)
	q.Set("q", title)

	for k, v := range scraper.opts {
		q.Set(k, v)
	}
	result, err := jikan.GetAnimeSearch(q)
	if err != nil {
		return nil, fmt.Errorf("failed to search for title: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("title not found")
	}

	anime := result.Data[0]
	var genres []string
	for _, g := range anime.Genres {
		genres = append(genres, g.Name)
	}
	for _, g := range anime.ExplicitGenres {
		genres = append(genres, g.Name)
	}

	return &ScrapeResult{
		ID:                strconv.Itoa(anime.MalId),
		Title:             anime.Title,
		AlternativeTitles: append([]string{anime.TitleJapanese}, anime.TitleSynonyms...),
		Synopsis:          anime.Synopsis,
		Background:        anime.Background,
		Genres:            genres,
		Rating:            anime.Rating,
	}, nil
}
