package v1

type PlaybackInfo struct {
	AssetPresentation  AssetPresentation `json:"assetPresentation"`
	AudioMode          AudioMode         `json:"audioMode"`
	AudioQuality       AudioQuality      `json:"audioQuality"`
	BitDepth           int               `json:"bitDepth"`
	Manifest           string            `json:"manifest"`
	ManifestHash       string            `json:"manifestHash"`
	ManifestMimeType   ManifestMimeType  `json:"manifestMimeType"`
	SampleRate         int               `json:"sampleRate"`
	TrackID            int               `json:"trackId"`
	TrackPeakAmplitude float64           `json:"trackPeakAmplitude"`
	TrackReplayGain    float64           `json:"trackReplayGain"`
}
