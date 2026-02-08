package player

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"codeberg.org/dergs/tonearm/internal/settings"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"github.com/go-gst/go-gst/gst"
)

var (
	currentlyEnqueuedTrack *v1.PlaybackInfo = &v1.PlaybackInfo{}
)

type btsStream struct {
	MimeType       string   `json:"mimeType"`
	Codecs         string   `json:"codecs"`
	EncryptionType string   `json:"encryptionType"`
	URLs           []string `json:"urls"`
}

func enqueueBTSStream(info *v1.PlaybackInfo) error {
	if info.ManifestMimeType != v1.ManifestMimeTypeAudioBTS {
		return fmt.Errorf("unsupported manifest mime type: %s", info.ManifestMimeType)
	}
	logger.Debug("selected track is in BTS format")

	decodedManifest, err := base64.StdEncoding.DecodeString(info.Manifest)
	if err != nil {
		return fmt.Errorf("failed to decode BTS manifest: %w", err)
	}

	var stream btsStream
	if err := json.Unmarshal(decodedManifest, &stream); err != nil {
		return fmt.Errorf("failed to unmarshal BTS stream: %w", err)
	}

	playbin.SetArg("uri", stream.URLs[0])
	playbin.Set("volume", settings.Player().GetVolume())
	return nil
}

func enqueueMPDStream(info *v1.PlaybackInfo) error {
	if info.ManifestMimeType != v1.ManifestMimeTypeAudioMPD {
		return fmt.Errorf("unsupported manifest mime type: %s", info.ManifestMimeType)
	}
	logger.Debug("selected track is in MPD format")

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
	playbin.Set("volume", settings.Player().GetVolume())
	return nil
}

func enqueue(playbackInfo *v1.PlaybackInfo) error {
	switch playbackInfo.ManifestMimeType {
	case v1.ManifestMimeTypeAudioMPD:
		enqueueMPDStream(playbackInfo)
	case v1.ManifestMimeTypeAudioBTS:
		enqueueBTSStream(playbackInfo)
	default:
		logger.Error("unsupported manifest mime type", "mime_type", playbackInfo.ManifestMimeType)
		return fmt.Errorf("unsupported manifest mime type: %s", playbackInfo.ManifestMimeType)
	}
	currentlyEnqueuedTrack = playbackInfo
	if settings.Playback().NormalizeVolume() {
		applyReplayGain(playbackInfo)
	}
	return nil
}

func applyReplayGain(playbackInfo *v1.PlaybackInfo) {
	sinkPad := rgvolume().GetStaticPad("sink")
	sinkPad.AddProbe(gst.PadProbeTypeEventDownstream, func(pad *gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		event := info.GetEvent()
		if event == nil {
			return gst.PadProbeOK
		}

		if event.Type() == gst.EventTypeSegment {
			injectReplayGainTags(rgvolume(), playbackInfo)
			return gst.PadProbeRemove
		}

		return gst.PadProbeOK
	})
}
