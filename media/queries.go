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

const createPerformerQuery = `
mutation CreatePerformer($performer: PerformerInput!) {
  createPerformer(performer: $performer) {
    identifiers { value { type value } }
    name { value }
  }
}
`

const getPerformerQuery = `
query GetPerformer($identifier: IdentifierInput!) {
  performer(identifier: $identifier) {
    identifiers { value { type value } }
    name { value }
  }
}
`

const createComposerQuery = `
mutation CreateComposer($composer: ComposerInput!) {
  createComposer(composer: $composer) {
    identifiers { value { type value } }
    firstName { value }
    lastName { value }
  }
}
`

const getComposerQuery = `
query GetComposer($identifier: IdentifierInput!) {
  composer(identifier: $identifier) {
    identifiers { value { type value } }
    first_name { value }
    last_name { value }
  }
}
`

const createRecordLabelQuery = `
mutation CreateRecordLabel($record_label: RecordLabelInput!) {
  createRecordLabel(record_label: $record_label) {
    identifiers { value { type value } }
    name { value }
  }
}
`

const getRecordLabelQuery = `
query GetRecordLabel($identifier: IdentifierInput!) {
  record_label(identifier: $identifier) {
    identifiers { value { type value } }
    name { value }
  }
}
`

const createPublisherQuery = `
mutation CreatePublisher($publisher: PublisherInput!) {
  createPublisher(publisher: $publisher) {
    identifiers { value { type value } }
    name { value }
  }
}
`

const getPublisherQuery = `
query GetPublisher($identifier: IdentifierInput!) {
  publisher(identifier: $identifier) {
    identifiers { value { type value } }
    name { value }
  }
}
`

const createRecordingQuery = `
mutation CreateRecording($recording: RecordingInput!) {
  createRecording(recording: $recording) {
    identifiers { value { type value } }
    title { value }
  }
}
`

const getRecordingQuery = `
query GetRecording($identifier: IdentifierInput!) {
  recording(identifier: $identifier) {
    identifiers { value { type value } }
    title { value }
    duration { value }
  }
}
`

const createWorkQuery = `
mutation CreateWork($work: WorkInput!) {
  createWork(work: $work) {
    identifiers { value { type value } }
    title { value }
  }
}
`

const getWorkQuery = `
query GetWork($identifier: IdentifierInput!) {
  work(identifier: $identifier) {
    identifiers { value { type value } }
    title { value }
  }
}
`
const createSongQuery = `
mutation CreateSong($song: SongInput!) {
  createSong(song: $song) {
    identifiers { value { type value } }
    title { value }
  }
}
`

const getSongQuery = `
query GetSong($identifier: IdentifierInput!) {
  song(identifier: $identifier) {
    identifiers { value { type value } }
    title { value }
    duration { value }
  }
}
`

const createReleaseQuery = `
mutation CreateRelease($release: ReleaseInput!) {
  createRelease(release: $release) {
    identifiers { value { type value } }
    title { value }
  }
}
`

const getReleaseQuery = `
query GetRelease($identifier: IdentifierInput!) {
  release(identifier: $identifier) {
    identifiers { value { type value } }
    type { value }
    title { value }
    date { value }
  }
}
`

const createPerformerRecordingLinkQuery = `
mutation CreatePerformerRecordingLink($link: PerformerRecordingLinkInput!) {
  createPerformerRecordingLink(link: $link) { source { name } }
}
`

const createComposerWorkLinkQuery = `
mutation CreateComposerWorkLink($link: ComposerWorkLinkInput!) {
  createComposerWorkLink(link: $link) { source { name } }
}
`

const createRecordLabelSongLinkQuery = `
mutation CreateRecordLabelSongLink($link: RecordLabelSongLinkInput!) {
  createRecordLabelSongLink(link: $link) { source { name } }
}
`

const createRecordLabelReleaseLinkQuery = `
mutation CreateRecordLabelReleaseLink($link: RecordLabelReleaseLinkInput!) {
  createRecordLabelReleaseLink(link: $link) { source { name } }
}
`

const createPublisherWorkLinkQuery = `
mutation CreatePublisherWorkLink($link: PublisherWorkLinkInput!) {
  createPublisherWorkLink(link: $link) { source { name } }
}
`

const createSongRecordingLinkQuery = `
mutation CreateSongRecordingLink($link: SongRecordingLinkInput!) {
  createSongRecordingLink(link: $link) { source { name } }
}
`

const createReleaseRecordingLinkQuery = `
mutation CreateReleaseRecordingLink($link: ReleaseRecordingLinkInput!) {
  createReleaseRecordingLink(link: $link) { source { name } }
}
`

const createRecordingWorkLinkQuery = `
mutation CreateRecordingWorkLink($link: RecordingWorkLinkInput!) {
  createRecordingWorkLink(link: $link) { source { name } }
}
`

const createReleaseSongLinkQuery = `
mutation CreateReleaseSongLink($link: ReleaseSongLinkInput!) {
  createReleaseSongLink(link: $link) { source { name } }
}
`
