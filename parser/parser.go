package parser

type Parser interface {
	Parse(name string) (TitleMetadata, error)
}

type TitleMetadata struct {
	Path    string            `json:"path"`
	Title   string            `json:"title"`
	Type    string            `json:"type,omitempty"`
	Episode int               `json:"episode,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
	Extra   []string          `json:"extra,omitempty"`
}

var Parsers = make(map[string]Parser)
