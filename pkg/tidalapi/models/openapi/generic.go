package openapi

import (
	"encoding/json"
	"fmt"
)

type ObjectType string

const (
	ObjectTypeAlbums    = "albums"
	ObjectTypeArtists   = "artists"
	ObjectTypeArtworks  = "artworks"
	ObjectTypeLyrics    = "lyrics"
	ObjectTypePlaylists = "playlists"
	ObjectTypeTracks    = "tracks"
)

type baseObject struct {
	ID            string          `json:"id"`
	RawAttributes json.RawMessage `json:"attributes"`
	Relationships map[string]struct {
		Links map[string]string `json:"links"`
	} `json:"relationships"`
	Type string `json:"type"`
}

type TypedObject[T any] struct {
	baseObject
	Attributes T `json:"attributes"`
}

type Object struct {
	baseObject
	Attributes Attributes
}

func (o *Object) MarshalJSON() ([]byte, error) {
	switch o.Type {
	case ObjectTypeAlbums:
		if raw, err := json.Marshal(o.Attributes.Album); err != nil {
			return nil, err
		} else {
			o.RawAttributes = raw
		}
	case ObjectTypeArtists:
		if raw, err := json.Marshal(o.Attributes.Artist); err != nil {
			return nil, err
		} else {
			o.RawAttributes = raw
		}
	case ObjectTypeArtworks:
		if raw, err := json.Marshal(o.Attributes.Artworks); err != nil {
			return nil, err
		} else {
			o.RawAttributes = raw
		}
	case ObjectTypeLyrics:
		if raw, err := json.Marshal(o.Attributes.Lyrics); err != nil {
			return nil, err
		} else {
			o.RawAttributes = raw
		}
	case ObjectTypePlaylists:
		if raw, err := json.Marshal(o.Attributes.Playlist); err != nil {
			return nil, err
		} else {
			o.RawAttributes = raw
		}
	case ObjectTypeTracks:
		if raw, err := json.Marshal(o.Attributes.Track); err != nil {
			return nil, err
		} else {
			o.RawAttributes = raw
		}
	default:
		return nil, fmt.Errorf("Unknown object type %s", o.Type)
	}
	return json.Marshal(o.baseObject)
}

func (o *Object) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &o.baseObject); err != nil {
		return err
	}

	o.Attributes = Attributes{}
	o.ID = o.baseObject.ID
	o.RawAttributes = o.baseObject.RawAttributes
	o.Relationships = o.baseObject.Relationships
	o.Type = o.baseObject.Type

	switch o.baseObject.Type {
	case ObjectTypeAlbums:
		if err := json.Unmarshal(o.baseObject.RawAttributes, &o.Attributes.Album); err != nil {
			return err
		}
	case ObjectTypeArtists:
		if err := json.Unmarshal(o.baseObject.RawAttributes, &o.Attributes.Artist); err != nil {
			return err
		}
	case ObjectTypeArtworks:
		if err := json.Unmarshal(o.baseObject.RawAttributes, &o.Attributes.Artworks); err != nil {
			return err
		}
	case ObjectTypeLyrics:
		if err := json.Unmarshal(o.baseObject.RawAttributes, &o.Attributes.Lyrics); err != nil {
			return err
		}
	case ObjectTypePlaylists:
		if err := json.Unmarshal(o.baseObject.RawAttributes, &o.Attributes.Playlist); err != nil {
			return err
		}
	case ObjectTypeTracks:
		if err := json.Unmarshal(o.baseObject.RawAttributes, &o.Attributes.Track); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown object type %s", o.baseObject.Type)
	}
	o.RawAttributes = nil
	o.baseObject.RawAttributes = nil
	return nil
}

type Attributes struct {
	Album    *AlbumAttributes
	Artist   *ArtistAttributes
	Artworks *ArtworkAttributes
	Lyrics   *LyricsAttributes
	Playlist *PlaylistAttributes
	Track    *TrackAttributes
}

type LyricsAttributes struct {
	Direction string `json:"direction"`
	LRCText   string `json:"lrcText"`
	Text      string `json:"text"`
}
