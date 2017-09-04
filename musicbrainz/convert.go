package musicbrainz

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Converter converts MusicBrainz data stored in a PostgeSQL database to META
// objects.
type Converter struct {
	db    *sql.DB
	store *meta.Store
}

// NewConverter returns a Converter which reads data from the given PostgreSQL
// database connection and stores META object in the given META store.
func NewConverter(db *sql.DB, store *meta.Store) *Converter {
	return &Converter{
		db:    db,
		store: store,
	}
}

// ConvertArtists reads all artists from the database, converts them to META
// objects, stores them in the META store and sends their CIDs to the given
// stream.
func (c *Converter) ConvertArtists(ctx context.Context, outStream chan *cid.Cid) error {
	// get all artists from the db
	rows, err := c.db.Query(artistsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// read the db row into an Artist struct, handling nullable
		// columns
		var (
			a          Artist
			typ        *string
			gender     *string
			area       *string
			beginDate  *string
			endDate    *string
			ipi        []byte
			isni       []byte
			alias      []byte
			annotation []byte
		)
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.SortName,
			&typ,
			&gender,
			&area,
			&beginDate,
			&endDate,
			&ipi,
			&isni,
			&alias,
			&a.MBID,
			&a.DisambiguationComment,
			&annotation,
		)
		if err != nil {
			return err
		}
		if typ != nil {
			a.Type = *typ
		}
		if gender != nil {
			a.Gender = *gender
		}
		if area != nil {
			a.Area = *area
		}
		if beginDate != nil {
			a.BeginDate = *beginDate
		}
		if endDate != nil {
			a.EndDate = *endDate
		}
		if len(ipi) > 2 {
			a.IPI = strings.Split(string(ipi)[1:len(ipi)-1], ",")
		}
		if len(isni) > 2 {
			a.ISNI = strings.Split(string(isni)[1:len(isni)-1], ",")
		}
		if len(alias) > 2 {
			a.Alias = strings.Split(string(alias)[1:len(alias)-1], ",")
		}
		if len(annotation) > 2 {
			a.Annotation = strings.Split(string(annotation)[1:len(annotation)-1], ",")
		}
		a.Context = ArtistContext

		// convert the artist to a META object
		obj, err := meta.Encode(a)
		if err != nil {
			return err
		}
		if err := c.store.Put(obj); err != nil {
			return err
		}

		// send the object's CID to the output stream
		select {
		case outStream <- obj.Cid():
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return rows.Err()
}
