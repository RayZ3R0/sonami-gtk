package player

import (
	"encoding/base64"
	"fmt"
	"os"

	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"
)

func enqueueMPDStream(info *v1.PlaybackInfo) error {
	if info.ManifestMimeType != v1.ManifestMimeTypeAudioMPD {
		return fmt.Errorf("unsupported manifest mime type: %s", info.ManifestMimeType)
	}

	file, err := os.CreateTemp("", "manifest-*.mpd")
	if err != nil {
		return err
	}
	defer file.Close()

	decoded, err := base64.StdEncoding.DecodeString(info.Manifest)
	if err != nil {
		return err
	}

	_, err = file.Write(decoded)
	if err != nil {
		return err
	}

	playbin.SetArg("uri", "file://"+file.Name())
	return nil
}
