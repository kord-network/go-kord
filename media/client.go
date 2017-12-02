// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	graphql "github.com/neelance/graphql-go"
)

type Client struct {
	url    string
	source *Source
}

func NewClient(url string, source *Source) *Client {
	return &Client{url, source}
}

func (c *Client) Performer(identifier *Identifier) (*Performer, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getPerformerQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Performer struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"performer"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Performer{Name: v.Performer.Name.Value}, nil
}

func (c *Client) Composer(identifier *Identifier) (*Composer, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getComposerQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Composer struct {
			FirstName struct {
				Value string `json:"value"`
			} `json:"first_name"`
			LastName struct {
				Value string `json:"value"`
			} `json:"last_name"`
		} `json:"composer"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Composer{
		FirstName: v.Composer.FirstName.Value,
		LastName:  v.Composer.LastName.Value,
	}, nil
}

func (c *Client) RecordLabel(identifier *Identifier) (*RecordLabel, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getRecordLabelQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		RecordLabel struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"record_label"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &RecordLabel{Name: v.RecordLabel.Name.Value}, nil
}

func (c *Client) Publisher(identifier *Identifier) (*Publisher, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getPublisherQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Publisher struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"publisher"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Publisher{Name: v.Publisher.Name.Value}, nil
}

func (c *Client) Recording(identifier *Identifier) (*Recording, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getRecordingQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Recording struct {
			Title struct {
				Value string `json:"value"`
			} `json:"title"`
			Duration struct {
				Value string `json:"value"`
			} `json:"duration"`
		} `json:"recording"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Recording{
		Title:    v.Recording.Title.Value,
		Duration: v.Recording.Duration.Value,
	}, nil
}

func (c *Client) Work(identifier *Identifier) (*Work, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getWorkQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Work struct {
			Title struct {
				Value string `json:"value"`
			} `json:"title"`
		} `json:"work"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Work{Title: v.Work.Title.Value}, nil
}

func (c *Client) Song(identifier *Identifier) (*Song, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getSongQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Song struct {
			Title struct {
				Value string `json:"value"`
			} `json:"title"`
			Duration struct {
				Value string `json:"value"`
			} `json:"duration"`
		} `json:"song"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Song{
		Title:    v.Song.Title.Value,
		Duration: v.Song.Duration.Value,
	}, nil
}

func (c *Client) Release(identifier *Identifier) (*Release, error) {
	res, err := c.perform(&graphqlQuery{
		Query:     getReleaseQuery,
		Variables: map[string]interface{}{"identifier": identifier},
	})
	if err != nil {
		return nil, err
	}
	var v struct {
		Recording struct {
			Type struct {
				Value string `json:"value"`
			} `json:"type"`
			Title struct {
				Value string `json:"value"`
			} `json:"title"`
			Date struct {
				Value string `json:"value"`
			} `json:"date"`
		} `json:"release"`
	}
	if err := json.Unmarshal(res, &v); err != nil {
		return nil, err
	}
	return &Release{
		Type:  v.Recording.Type.Value,
		Title: v.Recording.Title.Value,
		Date:  v.Recording.Date.Value,
	}, nil
}

func (c *Client) CreatePerformer(performer *Performer, identifier *Identifier) error {
	return c.createResource(
		"performer",
		createPerformerQuery,
		map[string]interface{}{
			"identifier": identifier,
			"name":       performer.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateComposer(composer *Composer, identifier *Identifier) error {
	return c.createResource(
		"composer",
		createComposerQuery,
		map[string]interface{}{
			"identifier": identifier,
			"firstName":  composer.FirstName,
			"lastName":   composer.LastName,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateRecordLabel(recordLabel *RecordLabel, identifier *Identifier) error {
	return c.createResource(
		"record_label",
		createRecordLabelQuery,
		map[string]interface{}{
			"identifier": identifier,
			"name":       recordLabel.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreatePublisher(publisher *Publisher, identifier *Identifier) error {
	return c.createResource(
		"publisher",
		createPublisherQuery,
		map[string]interface{}{
			"identifier": identifier,
			"name":       publisher.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateRecording(recording *Recording, identifier *Identifier) error {
	return c.createResource(
		"recording",
		createRecordingQuery,
		map[string]interface{}{
			"identifier": identifier,
			"title":      recording.Title,
			"duration":   recording.Duration,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateWork(work *Work, identifier *Identifier) error {
	return c.createResource(
		"work",
		createWorkQuery,
		map[string]interface{}{
			"identifier": identifier,
			"title":      work.Title,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateSong(song *Song, identifier *Identifier) error {
	return c.createResource(
		"song",
		createSongQuery,
		map[string]interface{}{
			"identifier": identifier,
			"title":      song.Title,
			"duration":   song.Duration,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateRelease(release *Release, identifier *Identifier) error {
	return c.createResource(
		"release",
		createReleaseQuery,
		map[string]interface{}{
			"identifier": identifier,
			"type":       release.Type,
			"title":      release.Title,
			"date":       release.Date,
			"source":     c.source,
		},
	)
}

func (c *Client) CreatePerformerRecordingLink(link *PerformerRecordingLink) error {
	return c.createResource(
		"link",
		createPerformerRecordingLinkQuery,
		map[string]interface{}{
			"performer_id": link.Performer,
			"recording_id": link.Recording,
			"role":         link.Role,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateComposerWorkLink(link *ComposerWorkLink) error {
	return c.createResource(
		"link",
		createComposerWorkLinkQuery,
		map[string]interface{}{
			"composer_id": link.Composer,
			"work_id":     link.Work,
			"role":        link.Role,
			"source":      c.source,
		},
	)
}

func (c *Client) CreateRecordLabelSongLink(link *RecordLabelSongLink) error {
	return c.createResource(
		"link",
		createRecordLabelSongLinkQuery,
		map[string]interface{}{
			"record_label_id": link.RecordLabel,
			"song_id":         link.Song,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateRecordLabelReleaseLink(link *RecordLabelReleaseLink) error {
	return c.createResource(
		"link",
		createRecordLabelReleaseLinkQuery,
		map[string]interface{}{
			"record_label_id": link.RecordLabel,
			"release_id":      link.Release,
			"source":          c.source,
		},
	)
}

func (c *Client) CreatePublisherWorkLink(link *PublisherWorkLink) error {
	return c.createResource(
		"link",
		createPublisherWorkLinkQuery,
		map[string]interface{}{
			"publisher_id": link.Publisher,
			"work_id":      link.Work,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateSongRecordingLink(link *SongRecordingLink) error {
	return c.createResource(
		"link",
		createSongRecordingLinkQuery,
		map[string]interface{}{
			"song_id":      link.Song,
			"recording_id": link.Recording,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateReleaseRecordingLink(link *ReleaseRecordingLink) error {
	return c.createResource(
		"link",
		createReleaseRecordingLinkQuery,
		map[string]interface{}{
			"release_id":   link.Release,
			"recording_id": link.Recording,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateRecordingWorkLink(link *RecordingWorkLink) error {
	return c.createResource(
		"link",
		createRecordingWorkLinkQuery,
		map[string]interface{}{
			"recording_id": link.Recording,
			"work_id":      link.Work,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateReleaseSongLink(link *ReleaseSongLink) error {
	return c.createResource(
		"link",
		createReleaseSongLinkQuery,
		map[string]interface{}{
			"release_id": link.Release,
			"song_id":    link.Song,
			"source":     c.source,
		},
	)
}

func (c *Client) createResource(name, query string, vars map[string]interface{}) error {
	_, err := c.perform(&graphqlQuery{
		Query:     query,
		Variables: map[string]interface{}{name: vars},
	})
	return err
}

type graphqlQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (c *Client) perform(query *graphqlQuery) (json.RawMessage, error) {
	data, _ := json.Marshal(query)
	req, err := http.NewRequest("POST", c.url+"/graphql", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected HTTP response: %s: %s", res.Status, body)
	}
	var r graphql.Response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error decoding GraphQL response: %s", err)
	}
	if len(r.Errors) > 0 {
		return nil, fmt.Errorf("unexpected errors in GraphQL response: %v", r.Errors)
	}
	return r.Data, nil
}
