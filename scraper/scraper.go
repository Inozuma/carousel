package scraper

var Scrapers = make(map[string]func(opts map[string]string) Scraper)

type Scraper interface {
	Scrape(title string) (*ScrapeResult, error)
}

type ScrapeResult struct {
	ID                string
	Title             string
	AlternativeTitles []string
	Synopsis          string
	Background        string
	Genres            []string
	Rating            string
}
