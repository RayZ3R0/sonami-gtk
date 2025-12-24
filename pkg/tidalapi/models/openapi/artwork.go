package openapi

type ArtworkAttributes struct {
	Files     ArtworkFiles `json:"files"`
	MediaType string       `json:"mediaType"`
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

type ArtworkFile struct {
	Href string `json:"href"`
	Meta struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"meta"`
}
