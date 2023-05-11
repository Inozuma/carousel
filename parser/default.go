package parser

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	Parsers["default"] = &DefaultParser{}
}

var (
	metaExpr = regexp.MustCompile(`[\(\[][^\]\)]+[\)\]]`)
)

type DefaultParser struct{}

func (p *DefaultParser) Parse(path string) (TitleMetadata, error) {
	_, name := filepath.Split(path)

	m := TitleMetadata{
		Path:  path,
		Title: strings.TrimSuffix(name, filepath.Ext(name)),
	}
	m.Title = strings.ReplaceAll(m.Title, "_", " ")
	matchIndexes := metaExpr.FindAllStringSubmatchIndex(name, -1)
	var offset int
	for _, indexes := range matchIndexes {
		begin, end := indexes[0], indexes[1]
		length := end - begin
		m.Title = string(append([]byte(m.Title[:begin-offset]), []byte(m.Title[end-offset:])...))
		offset += length
		m.Extra = append(m.Extra, strings.Trim(name[begin:end], "[()]"))
	}

	// parse episode
	m.Title, m.Episode = parseEpisode(m.Title)

	return m, nil
}

var numSepExprs = []*regexp.Regexp{
	regexp.MustCompile(`((#|＃|-|V[oO][lL]\.|ep\.|- Episode|- Chapter|Ep|Chapter|第)\s?)[0]*(?P<EP>\d+)\s?(END|_|巻)?`),
	regexp.MustCompile(`\s[0]*(?P<EP>\d+)$`),
}

func parseEpisode(s string) (string, int) {
	var ep int
	s = strings.TrimSpace(s)
	for _, re := range numSepExprs {
		matches := re.FindAllStringSubmatchIndex(s, -1)
		epIdx := re.SubexpIndex("EP")
		if len(matches) > 0 {
			epBegin, epEnd := matches[0][epIdx*2], matches[0][epIdx*2+1]
			ep, _ = strconv.Atoi(s[epBegin:epEnd])
			s = strings.TrimSpace(s[:matches[0][0]])
			break
		}
	}
	return s, ep
}
