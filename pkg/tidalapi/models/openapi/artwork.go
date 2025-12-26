package openapi

import "encoding/json"

type Artwork Response[ArtworkData]

const ObjectTypeArtworks = "artworks"

type ArtworkData struct {
	Attributes   ArtworkAttributes    `json:"attributes"`
	ID           string               `json:"id"`
	Relationship ArtworkRelationships `json:"relationships"`
	Type         string               `json:"type"`
}

type ArtworkAttributes struct {
	Files     ArtworkFiles `json:"files"`
	MediaType string       `json:"mediaType"`
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
	for _, file := range files {
		if file.Meta.Height >= size || file.Meta.Width >= size {
			return file
		}
	}
	return files[0]
}

func (i IncludedObjects) PlainArtworks(relationships ...Relationship) []ArtworkData {
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
