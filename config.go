package main

// Config
type Config struct {
	// Libraries list of library configuration.
	Libraries []LibraryConfig `json:"libraries"`
	// Database path.
	Database string `json:"database"`
	// PreviewDirectory where thumbnails are stored.
	PreviewDirectory string `json:"preview_dir"`
}

// LibraryConfig
type LibraryConfig struct {
	// Name of the library
	Name string `json:"name"`
	// Path to the library directory where media files are located.
	Path string `json:"path"`
	// Scraper to use for metadata. Optional.
	Scraper *ScraperConfig `json:"scraper,omitempty"`
	// Parsers to use to find metadata from file name and info.
	Parser string `json:"parser,omitempty"`
}

// ScraperConfig
type ScraperConfig struct {
	Name    string            `json:"name"`
	Options map[string]string `json:"options,omitempty"`
}
