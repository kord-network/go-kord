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
	"github.com/meta-network/go-meta/graphql"
)

type Client struct {
	*graphql.Client

	source *Source
}

func NewClient(url string, source *Source) *Client {
	return &Client{graphql.NewClient(url), source}
}

func (c *Client) Query(query string, variables graphql.Variables, out interface{}) error {
	return c.Do(query, variables, out)
}

func (c *Client) Performer(identifier *Identifier) (*Performer, error) {
	var v struct {
		Performer struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"performer"`
	}
	if err := c.Query(
		getPerformerQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Performer{Name: v.Performer.Name.Value}, nil
}

func (c *Client) Composer(identifier *Identifier) (*Composer, error) {
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
	if err := c.Query(
		getComposerQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Composer{
		FirstName: v.Composer.FirstName.Value,
		LastName:  v.Composer.LastName.Value,
	}, nil
}

func (c *Client) RecordLabel(identifier *Identifier) (*RecordLabel, error) {
	var v struct {
		RecordLabel struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"record_label"`
	}
	if err := c.Query(
		getRecordLabelQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &RecordLabel{Name: v.RecordLabel.Name.Value}, nil
}

func (c *Client) Publisher(identifier *Identifier) (*Publisher, error) {
	var v struct {
		Publisher struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"publisher"`
	}
	if err := c.Query(
		getPublisherQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Publisher{Name: v.Publisher.Name.Value}, nil
}

func (c *Client) Recording(identifier *Identifier) (*Recording, error) {
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
	if err := c.Query(
		getRecordingQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Recording{
		Title:    v.Recording.Title.Value,
		Duration: v.Recording.Duration.Value,
	}, nil
}

func (c *Client) Work(identifier *Identifier) (*Work, error) {
	var v struct {
		Work struct {
			Title struct {
				Value string `json:"value"`
			} `json:"title"`
		} `json:"work"`
	}
	if err := c.Query(
		getWorkQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Work{Title: v.Work.Title.Value}, nil
}

func (c *Client) Song(identifier *Identifier) (*Song, error) {
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
	if err := c.Query(
		getSongQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Song{
		Title:    v.Song.Title.Value,
		Duration: v.Song.Duration.Value,
	}, nil
}

func (c *Client) Release(identifier *Identifier) (*Release, error) {
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
	if err := c.Query(
		getReleaseQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Release{
		Type:  v.Recording.Type.Value,
		Title: v.Recording.Title.Value,
		Date:  v.Recording.Date.Value,
	}, nil
}

func (c *Client) Organisation(identifier *Identifier) (*Organisation, error) {
	var v struct {
		Organisation struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"organisation"`
	}
	if err := c.Query(
		getOrganisationQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Organisation{Name: v.Organisation.Name.Value}, nil
}

func (c *Client) Series(identifier *Identifier) (*Series, error) {
	var v struct {
		Series struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"series"`
	}
	if err := c.Query(
		getSeriesQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Series{Name: v.Series.Name.Value}, nil
}

func (c *Client) Season(identifier *Identifier) (*Season, error) {
	var v struct {
		Season struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"season"`
	}
	if err := c.Query(
		getSeasonQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Season{Name: v.Season.Name.Value}, nil
}

func (c *Client) Episode(identifier *Identifier) (*Episode, error) {
	var v struct {
		Episode struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"episode"`
	}
	if err := c.Query(
		getEpisodeQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Episode{Name: v.Episode.Name.Value}, nil
}

func (c *Client) Supplemental(identifier *Identifier) (*Supplemental, error) {
	var v struct {
		Supplemental struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
		} `json:"supplemental"`
	}
	if err := c.Query(
		getSupplementalQuery,
		graphql.Variables{"identifier": identifier},
		&v,
	); err != nil {
		return nil, err
	}
	return &Supplemental{Name: v.Supplemental.Name.Value}, nil
}

func (c *Client) CreatePerformer(performer *Performer, identifier *Identifier) error {
	return c.createResource(
		"performer",
		createPerformerQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       performer.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateContributor(contributor *Contributor, identifier *Identifier) error {
	return c.createResource(
		"contributor",
		createContributorQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       contributor.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateComposer(composer *Composer, identifier *Identifier) error {
	return c.createResource(
		"composer",
		createComposerQuery,
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
			"identifier": identifier,
			"type":       release.Type,
			"title":      release.Title,
			"date":       release.Date,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateOrganisation(organisation *Organisation, identifier *Identifier) error {
	return c.createResource(
		"organisation",
		createOrganisationQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       organisation.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateSeries(series *Series, identifier *Identifier) error {
	return c.createResource(
		"series",
		createSeriesQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       series.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateSeason(season *Season, identifier *Identifier) error {
	return c.createResource(
		"season",
		createSeasonQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       season.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateEpisode(episode *Episode, identifier *Identifier) error {
	return c.createResource(
		"episode",
		createEpisodeQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       episode.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateSupplemental(supplemental *Supplemental, identifier *Identifier) error {
	return c.createResource(
		"supplemental",
		createSupplementalQuery,
		graphql.Variables{
			"identifier": identifier,
			"name":       supplemental.Name,
			"source":     c.source,
		},
	)
}

func (c *Client) CreatePerformerRecordingLink(link *PerformerRecordingLink) error {
	return c.createResource(
		"link",
		createPerformerRecordingLinkQuery,
		graphql.Variables{
			"performer_id": link.Performer,
			"recording_id": link.Recording,
			"role":         link.Role,
			"source":       c.source,
		},
	)
}

func (c *Client) CreatePerformerSongLink(link *PerformerSongLink) error {
	return c.createResource(
		"link",
		createPerformerSongLinkQuery,
		graphql.Variables{
			"performer_id": link.Performer,
			"song_id":      link.Song,
			"role":         link.Role,
			"source":       c.source,
		},
	)
}

func (c *Client) CreatePerformerReleaseLink(link *PerformerReleaseLink) error {
	return c.createResource(
		"link",
		createPerformerReleaseLinkQuery,
		graphql.Variables{
			"performer_id": link.Performer,
			"release_id":   link.Release,
			"role":         link.Role,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateContributorRecordingLink(link *ContributorRecordingLink) error {
	return c.createResource(
		"link",
		createContributorRecordingLinkQuery,
		graphql.Variables{
			"contributor_id": link.Contributor,
			"recording_id":   link.Recording,
			"role":           link.Role,
			"source":         c.source,
		},
	)
}

func (c *Client) CreateComposerWorkLink(link *ComposerWorkLink) error {
	return c.createResource(
		"link",
		createComposerWorkLinkQuery,
		graphql.Variables{
			"composer_id": link.Composer,
			"work_id":     link.Work,
			"role":        link.Role,
			"pr_share":    link.PerformanceRightsShare,
			"mr_share":    link.MechanicalRightsShare,
			"sr_share":    link.SynchronizationRightsShare,
			"source":      c.source,
		},
	)
}

func (c *Client) CreateRecordLabelRecordingLink(link *RecordLabelRecordingLink) error {
	return c.createResource(
		"link",
		createRecordLabelRecordingLinkQuery,
		graphql.Variables{
			"record_label_id": link.RecordLabel,
			"recording_id":    link.Recording,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateRecordLabelSongLink(link *RecordLabelSongLink) error {
	return c.createResource(
		"link",
		createRecordLabelSongLinkQuery,
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
			"publisher_id": link.Publisher,
			"work_id":      link.Work,
			"role":         link.Role,
			"pr_share":     link.PerformanceRightsShare,
			"mr_share":     link.MechanicalRightsShare,
			"sr_share":     link.SynchronizationRightsShare,
			"source":       c.source,
		},
	)
}

func (c *Client) CreateSongRecordingLink(link *SongRecordingLink) error {
	return c.createResource(
		"link",
		createSongRecordingLinkQuery,
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
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
		graphql.Variables{
			"release_id": link.Release,
			"song_id":    link.Song,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateOrganisationSeriesLink(link *OrganisationSeriesLink) error {
	return c.createResource(
		"link",
		createOrganisationSeriesLinkQuery,
		graphql.Variables{
			"organisation_id": link.Organisation,
			"series_id":       link.Series,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateOrganisationSeasonLink(link *OrganisationSeasonLink) error {
	return c.createResource(
		"link",
		createOrganisationSeasonLinkQuery,
		graphql.Variables{
			"organisation_id": link.Organisation,
			"season_id":       link.Season,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateOrganisationEpisodeLink(link *OrganisationEpisodeLink) error {
	return c.createResource(
		"link",
		createOrganisationEpisodeLinkQuery,
		graphql.Variables{
			"organisation_id": link.Organisation,
			"episode_id":      link.Episode,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateOrganisationSupplementalLink(link *OrganisationSupplementalLink) error {
	return c.createResource(
		"link",
		createOrganisationSupplementalLinkQuery,
		graphql.Variables{
			"organisation_id": link.Organisation,
			"supplemental_id": link.Supplemental,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateSeriesSeasonLink(link *SeriesSeasonLink) error {
	return c.createResource(
		"link",
		createSeriesSeasonLinkQuery,
		graphql.Variables{
			"series_id": link.Series,
			"season_id": link.Season,
			"source":    c.source,
		},
	)
}

func (c *Client) CreateSeriesEpisodeLink(link *SeriesEpisodeLink) error {
	return c.createResource(
		"link",
		createSeriesEpisodeLinkQuery,
		graphql.Variables{
			"series_id":  link.Series,
			"episode_id": link.Episode,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateSeriesSupplementalLink(link *SeriesSupplementalLink) error {
	return c.createResource(
		"link",
		createSeriesSupplementalLinkQuery,
		graphql.Variables{
			"series_id":       link.Series,
			"supplemental_id": link.Supplemental,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateSeasonEpisodeLink(link *SeasonEpisodeLink) error {
	return c.createResource(
		"link",
		createSeasonEpisodeLinkQuery,
		graphql.Variables{
			"season_id":  link.Season,
			"episode_id": link.Episode,
			"source":     c.source,
		},
	)
}

func (c *Client) CreateSeasonSupplementalLink(link *SeasonSupplementalLink) error {
	return c.createResource(
		"link",
		createSeasonSupplementalLinkQuery,
		graphql.Variables{
			"season_id":       link.Season,
			"supplemental_id": link.Supplemental,
			"source":          c.source,
		},
	)
}

func (c *Client) CreateEpisodeSupplementalLink(link *EpisodeSupplementalLink) error {
	return c.createResource(
		"link",
		createEpisodeSupplementalLinkQuery,
		graphql.Variables{
			"episode_id":      link.Episode,
			"supplemental_id": link.Supplemental,
			"source":          c.source,
		},
	)
}

func (c *Client) createResource(name, query string, variables graphql.Variables) error {
	return c.Do(query, graphql.Variables{name: variables}, nil)
}
