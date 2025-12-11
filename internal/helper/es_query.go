package helper

type ESQuery struct {
	Query BoolQuery `json:"query"`
}

type BoolQuery struct {
	Bool BoolShould `json:"bool"`
}

type BoolShould struct {
	Should             []any `json:"should"`
	MinimalShouldMatch int   `json:"minimum_should_match"`
}

type MatchPhrase struct {
	MatchPhrase map[string]MatchPhraseField `json:"match_phrase"`
}

type MatchPhraseField struct {
	Query string  `json:"query"`
	Slop  int     `json:"slop"`
	Boost float64 `json:"boost,omitempty"`
}

type MultiMatch struct {
	MultiMatch MultiMatchField `json:"multi_match"`
}

type MultiMatchField struct {
	Query    string   `json:"query"`
	Fields   []string `json:"fields"`
	Fuzzines string   `json:"fuzziness,omitempty"`
}
