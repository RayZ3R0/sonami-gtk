package player

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
)

type BTSStream struct {
	MimeType       string   `json:"mimeType"`
	Codecs         string   `json:"codecs"`
	EncryptionType string   `json:"encryptionType"`
	URLs           []string `json:"urls"`
}

func enqueueBTSStream(info *v1.PlaybackInfo) error {
	if info.ManifestMimeType != v1.ManifestMimeTypeAudioBTS {
		return fmt.Errorf("unsupported manifest mime type: %s", info.ManifestMimeType)
	}

	decodedManifest, err := base64.StdEncoding.DecodeString(info.Manifest)
	if err != nil {
		return fmt.Errorf("failed to decode BTS manifest: %w", err)
	}

	var stream BTSStream
	if err := json.Unmarshal(decodedManifest, &stream); err != nil {
		return fmt.Errorf("failed to unmarshal BTS stream: %w", err)
	}

	playbin.SetArg("uri", stream.URLs[0])
	return nil
}
