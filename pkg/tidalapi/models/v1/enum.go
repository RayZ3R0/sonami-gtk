package v1

type AssetPresentation string

const (
	AssetPresentationFull    AssetPresentation = "FULL"
	AssetPresentationPreview AssetPresentation = "PREVIEW"
)

type AudioMode string

const (
	AudioModeStereo AudioMode = "STEREO"
)

type AudioQuality string

const (
	AudioQualityHighResLossess AudioQuality = "HI_RES_LOSSLESS"
)

type ManifestMimeType string

const (
	ManifestMimeTypeAudioBTS  ManifestMimeType = "application/vnd.tidal.bts"
	ManifestMimeTypeAudioMPD  ManifestMimeType = "application/dash+xml"
	ManifestMimeTypeVideoMP2T ManifestMimeType = "video/mp2t"
)

type PlaybackMode string

const (
	PlaybackModeStream PlaybackMode = "STREAM"
)
