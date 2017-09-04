package musicbrainz

// Context represents a JSON-LD context.
type Context map[string]string

// ArtistContext is the JSON-LD context to use for META objects representing
// MusicBrainz artists.
var ArtistContext = Context{
	"name":       "https://musicbrainz.org/doc/Artist#Name",
	"sortName":   "https://musicbrainz.org/doc/Artist#Sort_name",
	"type":       "https://musicbrainz.org/doc/Artist#Type",
	"gender":     "https://musicbrainz.org/doc/Artist#Gender",
	"area":       "https://musicbrainz.org/doc/Artist#Area",
	"begin_date": "https://musicbrainz.org/doc/Artist#Begin_and_end_dates",
	"end_date":   "https://musicbrainz.org/doc/Artist#Begin_and_end_dates",
	"ipi":        "https://musicbrainz.org/doc/Artist#IPI_code",
	"isni":       "https://musicbrainz.org/doc/Artist#ISNI_code",
	"alias":      "https://musicbrainz.org/doc/Artist#Alias",
	"mbid":       "https://musicbrainz.org/doc/Artist#MBID",
	"disambiguation_comment": "https://musicbrainz.org/doc/Artist#Disambiguation_comment",
	"annotation":             "https://musicbrainz.org/doc/Artist#Annotation",
}

// Artist represents a MusicBrainz artist, see
// https://musicbrainz.org/doc/Artist
type Artist struct {
	Context               Context  `json:"@context"`
	ID                    int64    `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	SortName              string   `json:"sortName,omitempty"`
	Type                  string   `json:"type,omitempty"`
	Gender                string   `json:"gender,omitempty"`
	Area                  string   `json:"area,omitempty"`
	BeginDate             string   `json:"begin_date,omitempty"`
	EndDate               string   `json:"end_date,omitempty"`
	IPI                   []string `json:"ipi,omitempty"`
	ISNI                  []string `json:"isni,omitempty"`
	Alias                 []string `json:"alias,omitempty"`
	MBID                  string   `json:"mbid,omitempty"`
	DisambiguationComment string   `json:"disambiguation_comment,omitempty"`
	Annotation            []string `json:"annotation,omitempty"`
}
