package lyrics

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"time"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

type netEaseSearchResponse struct {
	Result *netEaseSearchResult `json:"result"`
	Code   int                  `json:"code"`
}

type netEaseSearchResult struct {
	Songs []netEaseSong `json:"songs"`
}

type netEaseSong struct {
	ID uint64 `json:"id"`
}

type netEaseLyricResponse struct {
	LRC  *netEaseLyricContent `json:"lrc"`
	Code int                  `json:"code"`
}

type netEaseLyricContent struct {
	Lyric string `json:"lyric"`
}

var netEaseHTTP = &http.Client{Timeout: 8 * time.Second}

const netEaseUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

func netEaseGet(fakeIP, rawURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", netEaseUserAgent)
	req.Header.Set("X-Real-IP", fakeIP)
	req.Header.Set("Referer", "http://music.163.com/")
	return netEaseHTTP.Do(req)
}

func getNetEaseLyrics(track sonami.Track) (lyrics string, isTimestamped bool, err error) {
	title := sonami.FormatTitle(track)
	artists := track.Artists().Names()
	artist := ""
	if len(artists) > 0 {
		artist = artists[0]
	}

	fakeIP := fmt.Sprintf("220.181.%d.%d", rand.IntN(255), rand.IntN(255))

	songID, err := netEaseSearchSong(fakeIP, title, artist)
	if err != nil || songID == 0 {
		return
	}

	lyrics, err = netEaseFetchLyrics(fakeIP, songID)
	if err != nil || lyrics == "" {
		return
	}
	isTimestamped = true
	return
}

func netEaseSearchSong(fakeIP, title, artist string) (uint64, error) {
	base, _ := url.Parse("http://music.163.com/api/search/get")
	q := base.Query()
	q.Set("s", artist+" "+title)
	q.Set("type", "1")
	q.Set("offset", "0")
	q.Set("limit", "5")
	base.RawQuery = q.Encode()

	resp, err := netEaseGet(fakeIP, base.String())
	if err != nil {
		return 0, fmt.Errorf("netease search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("netease search: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data netEaseSearchResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("netease search parse: %w", err)
	}

	if data.Result == nil || len(data.Result.Songs) == 0 {
		return 0, nil
	}

	return data.Result.Songs[0].ID, nil
}

func netEaseFetchLyrics(fakeIP string, songID uint64) (string, error) {
	base, _ := url.Parse("http://music.163.com/api/song/lyric")
	q := base.Query()
	q.Set("id", fmt.Sprintf("%d", songID))
	q.Set("lv", "-1")
	q.Set("kv", "-1")
	q.Set("tv", "-1")
	base.RawQuery = q.Encode()

	resp, err := netEaseGet(fakeIP, base.String())
	if err != nil {
		return "", fmt.Errorf("netease lyrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("netease lyrics: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data netEaseLyricResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("netease lyrics parse: %w", err)
	}

	if data.LRC != nil {
		return data.LRC.Lyric, nil
	}
	return "", nil
}
