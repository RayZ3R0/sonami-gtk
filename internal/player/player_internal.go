package player

import (
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/infinytum/injector"
)

func playAlbum(albumId string, shuffle bool, startingPosition int) error {
	service, err := injector.Inject[tonearm.Service]()
	if err != nil {
		return err
	}

	// We request the album from the service, hinting that we will query the tracks, to start playback
	// and the cover art for the SourceChanged signal.
	album, err := service.GetAlbumInfo(albumId)
	if err != nil {
		return err
	}

	trackPaginator, err := service.GetAlbumTracks(albumId)
	if err != nil {
		return err
	}

	// We need to get all tracks in the album to be able to properly display them in a queue
	// and in case the user wants to shuffle-play we need them anyways so they can be shuffled.
	tracks, err := trackPaginator.GetAll()
	if err != nil {
		return err
	}

	if _, err := playTracklist(tracks, shuffle, startingPosition); err != nil {
		return err
	}

	// The tracklist has successfully started playing, we can now notify subscribers
	// where the playback originated from.
	SourceChanged.Notify(func(oldValue tonearm.PlaybackSource) tonearm.PlaybackSource {
		return album
	})

	return nil
}

func playPlaylist(playlistId string, shuffle bool, startingPosition int) error {
	service, err := injector.Inject[tonearm.Service]()
	if err != nil {
		return err
	}

	// We request the playlist from the service, hinting that we will query the tracks, to start playback
	// and the cover art for the SourceChanged signal.
	playlist, err := service.GetPlaylistInfo(playlistId)
	if err != nil {
		return err
	}

	trackPaginator, err := service.GetPlaylistTracks(playlistId)
	if err != nil {
		return err
	}

	// We need to get all tracks in the playlist to be able to properly display them in a queue
	// and in case the user wants to shuffle-play we need them anyways so they can be shuffled.
	tracks, err := trackPaginator.GetAll()
	if err != nil {
		return err
	}

	if _, err := playTracklist(tracks, shuffle, startingPosition); err != nil {
		return err
	}

	// The tracklist has successfully started playing, we can now notify subscribers
	// where the playback originated from.
	SourceChanged.Notify(func(oldValue tonearm.PlaybackSource) tonearm.PlaybackSource {
		return playlist
	})

	return nil
}

func playTracklist(tracks []tonearm.Track, shuffle bool, startingPosition int) (tonearm.Track, error) {
	// When a tracklist change is initiated, all queues are cleared
	clearQueues()

	// We want to notify everyone that we have stopped playing whatever was currently playing
	// as the current queue has just been invalidated
	Stop()
	TrackChanged.Notify(func(oldValue tonearm.Track) tonearm.Track {
		return nil
	})

	// Now we want to register the given tracklist as the base queue
	BaseQueue.Set(tracks)

	// If the user wants to start playing the tracklist at an offset
	// we skip the queue to the desired position
	if startingPosition > 0 {
		BaseQueue.Skip(startingPosition)
	}

	// If the user wants to play the tracklist in shuffle mode, we shuffle the queue
	// and enable the shuffle mode for the player.
	if shuffle {
		SetShuffle(true)
	}

	// Now we need to kick-off playback by fetching the next track from the queue
	// and starting playback.
	nextTrack := getNextTrackFromQueue(false)
	if nextTrack != nil {
		logger.Info("starting tracklist playback", "track_id", nextTrack.ID())
		if err := playTrack(nextTrack); err != nil {
			return nil, err
		}

		// The track has started playing, so we can add it to the playback history
		history.Push(&HistoryEntry{
			TrackID: nextTrack.ID(),
		})
	}

	return nextTrack, nil
}
