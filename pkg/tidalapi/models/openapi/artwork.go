package openapi

import (
	"encoding/json"
	"sort"
)

type Artwork Response[ArtworkData]

const ObjectTypeArtworks = "artworks"

type ArtworkMediaType string

const (
	ArtworkMediaTypeImage ArtworkMediaType = "IMAGE"
	ArtworkMediaTypeVideo ArtworkMediaType = "VIDEO"
)

type ArtworkData struct {
	Attributes   ArtworkAttributes    `json:"attributes"`
	ID           string               `json:"id"`
	Relationship ArtworkRelationships `json:"relationships"`
	Type         string               `json:"type"`
}

type ArtworkAttributes struct {
	Files     ArtworkFiles     `json:"files"`
	MediaType ArtworkMediaType `json:"mediaType"`
}

func (a ArtworkAttributes) IsPicture() bool {
	return a.MediaType == ArtworkMediaTypeImage
}

type ArtworkRelationships struct {
	Owners Response[[]Relationship] `json:"owners"`
}

type ArtworkFile struct {
	Href string `json:"href"`
	Meta struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"meta"`
}

type ArtworkFiles []ArtworkFile

func (files ArtworkFiles) AtLeast(size int) ArtworkFile {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Meta.Height < files[j].Meta.Height || files[i].Meta.Width < files[j].Meta.Width
	})
	for _, file := range files {
		if min(file.Meta.Height, file.Meta.Width) >= size {
			return file
		}
	}
	return files[0]
}

type Artworks []ArtworkData

func (artworks Artworks) AtLeast(size int) string {
	coverUrl := ""
	for _, artwork := range artworks {
		if artwork.Attributes.IsPicture() {
			coverUrl = artwork.Attributes.Files.AtLeast(size).Href
			break
		}
	}
	return coverUrl
}

func (i IncludedObjects) PlainArtworks(relationships ...Relationship) Artworks {
	var objects = i.FromRelationships(relationships, ObjectTypeArtworks)

	var artworks []ArtworkData
	for _, obj := range objects {
		var artwork ArtworkData
		if err := json.Unmarshal(obj.Raw, &artwork); err != nil {
			continue
		}
		artworks = append(artworks, artwork)
	}
	return artworks
}
