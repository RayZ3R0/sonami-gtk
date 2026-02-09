package player

import (
	"fmt"
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/g"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"github.com/go-gst/go-gst/gst"
)

var (
	rgvolume = g.Lazy(func() *gst.Element {
		rgvolume, err := gst.NewElementWithName("rgvolume", "rgvolume")
		if err != nil {
			panic(err)
		}
		rgvolume.Set("album-mode", false)
		rgvolume.Set("pre-amp", 0.0)
		rgvolume.Set("fallback-gain", 0.0)
		return rgvolume
	})
)

func buildReplayGainFilterBin() (*gst.Bin, error) {
	bin := gst.NewBin("replaygain-bin")

	rglimiter, err := gst.NewElementWithName("rglimiter", "rglimiter")
	if err != nil {
		return nil, err
	}

	bin.Add(rgvolume())
	bin.Add(rglimiter)

	rgvolume().Link(rglimiter)

	sinkPad := rgvolume().GetStaticPad("sink")
	ghostSink := gst.NewGhostPad("sink", sinkPad)
	bin.AddPad(ghostSink.Pad)

	srcPad := rglimiter.GetStaticPad("src")
	ghostSrc := gst.NewGhostPad("src", srcPad)
	bin.AddPad(ghostSrc.Pad)

	return bin, nil
}

func injectReplayGainTags(rgvolume *gst.Element, info *v1.PlaybackInfo) {
	tagList := gst.NewEmptyTagList()

	tagList.AddValue(gst.TagMergeReplaceAll, gst.TagTrackGain, info.TrackReplayGain)
	tagList.AddValue(gst.TagMergeAppend, gst.TagTrackPeak, info.TrackPeakAmplitude)

	tagEvent := gst.NewTagEvent(tagList)
	sinkPad := rgvolume.GetStaticPad("sink")
	if ok := sinkPad.SendEvent(tagEvent); !ok {
		slog.Warn("failed to send ReplayGain tag event")
	} else {
		slog.Info("injected ReplayGain tags", "track", fmt.Sprintf("%.2f dB", info.TrackReplayGain), "peak", fmt.Sprintf("%.6f", info.TrackPeakAmplitude))
	}
}
