package feed

import (
	"time"
)

type ActivityType string

const (
	ActivityTypeNewAlbumRelease ActivityType = "NEW_ALBUM_RELEASE"
	ActivityTypeNewHistoryMix   ActivityType = "NEW_HISTORY_MIX"
)

type FollowableActivity struct {
	HistoryMix *HistoryMix `json:"historyMix"`
	Album      *Album      `json:"album"`

	ActivityType ActivityType `json:"activityType"`
	OccuredAt    time.Time    `json:"occurredAt"`
}

type Activity struct {
	FollowableActivity FollowableActivity
	Seen               bool
}

type MixType string

const (
	MixTypeHistoryMonthly MixType = "HISTORY_MONTHLY_MIX"
)

type imageData struct {
	Width, Height int
	Url           string
}

type HistoryMix struct {
	Id string
	MixType
	TitleTextInfo, SubtitleTextInfo struct {
		text, color string
	}
	Updated              time.Time
	Images, DetailImages struct {
		Small, Medium, Large imageData
	}
	Master          bool
	Title, Subtitle string
}

type AlbumType string

const (
	AlbumTypeSingle AlbumType = "SINGLE"
	AlbumTypeAlbum  AlbumType = "ALBUM"
	AlbumTypeEP     AlbumType = "EP"
)

type Album struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Artists []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Main bool   `json:"main"`
	} `json:"artists"`
	MainArtists []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"mainArtists"`
	Type                   AlbumType `json:"type"`
	AudioQuality           string    `json:"audioQuality"`
	StreamStartDate        time.Time `json:"streamStartDate"`
	ReleaseDate            string    `json:"releaseDate"`
	AllowStreaming         bool      `json:"allowStreaming"`
	StreamReady            bool      `json:"streamReady"`
	Cover                  string    `json:"cover"`
	NumberOfVolumes        int       `json:"numberOfVolumes"`
	NumberOfTracks         int       `json:"numberOfTracks"`
	NumberOfVideos         int       `json:"numberOfVideos"`
	Explicit               bool      `json:"explicit"`
	AudioModes             []string  `json:"audioModes"`
	AdSupportedStreamReady bool      `json:"adSupportedStreamReady"`
	MediaMetadata          struct {
		Tags []string `json:"tags"`
	} `json:"mediaMetadata"`
	ProviderID int  `json:"providerId"`
	DjReady    bool `json:"djReady"`
	StemReady  bool `json:"stemReady"`
	Upload     bool `json:"upload"`
}
