package jellyfin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	ErrInvalidToken       = errors.New("invalid jellyfin api authorization token")
	ErrUnknownApiError    = errors.New("unknown error calling the jellyfin api")
	ErrUnknownApiResponse = errors.New("failed to parse the response body correctly")
	ErrTechnicalError     = errors.New("error with the request")
)

type Client struct {
	jHost      string
	jToken     string
	httpClient http.Client
}

func NewClient(jellyfinHost, jellyfinToken string) *Client {
	return &Client{
		jHost:  jellyfinHost,
		jToken: jellyfinToken,
		httpClient: http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (c *Client) GetActiveStreamsPerUser() (map[string]int, error) {
	sessions, err := queryJellyfinApi[session](fmt.Sprintf("%s/Sessions?ApiKey=%s", c.jHost, c.jToken), c.httpClient)
	if err != nil {
		return nil, err
	}

	userCountMap := make(map[string]int)
	for i := 0; i < len(sessions); i++ {
		if sessions[i].IsActive && sessions[i].PlayState.PlayMethod != "" && !sessions[i].PlayState.IsPaused {
			userCountMap[sessions[i].UserName]++
		}
	}

	return userCountMap, nil
}

func (c *Client) GetConnectedDevicesPerUser() (map[string]int, error) {
	sessions, err := queryJellyfinApi[session](fmt.Sprintf("%s/Sessions?ApiKey=%s", c.jHost, c.jToken), c.httpClient)
	if err != nil {
		return nil, err
	}

	userCountMap := make(map[string]int)
	for i := 0; i < len(sessions); i++ {
		if sessions[i].IsActive {
			userCountMap[sessions[i].UserName]++
		}
	}

	return userCountMap, nil
}

func (c *Client) GetMediaByType() (map[string]int, error) {
	counts, err := queryJellyfinApi[mediaCounts](fmt.Sprintf("%s/Items/Counts?ApiKey=%s", c.jHost, c.jToken), c.httpClient)
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	// Only add counts over zero to reduce metric cardinality
	addIfCountOverZero("albums", counts.AlbumCount, countMap)
	addIfCountOverZero("artists", counts.ArtistCount, countMap)
	addIfCountOverZero("books", counts.BookCount, countMap)
	addIfCountOverZero("boxSets", counts.BoxSetCount, countMap)
	addIfCountOverZero("episodes", counts.EpisodeCount, countMap)
	addIfCountOverZero("items", counts.ItemCount, countMap)
	addIfCountOverZero("movies", counts.MovieCount, countMap)
	addIfCountOverZero("music_videos", counts.MusicVideoCount, countMap)
	addIfCountOverZero("programs", counts.ProgramCount, countMap)
	addIfCountOverZero("series", counts.SeriesCount, countMap)
	addIfCountOverZero("songs", counts.SongCount, countMap)
	addIfCountOverZero("trailers", counts.TrailerCount, countMap)

	return countMap, nil
}

func (c *Client) ValidateToken() error {
	_, err := queryJellyfinApi[mediaCounts](fmt.Sprintf("%s/Items/Counts?ApiKey=%s", c.jHost, c.jToken), c.httpClient)
	if errors.Is(ErrInvalidToken, err) {
		return err
	}
	return nil
}

func addIfCountOverZero(key string, count int, countMap map[string]int) {
	if count > 0 {
		countMap[key] = count
	}
}

func queryJellyfinApi[response any](url string, client http.Client) (r response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return r, errors.Join(ErrTechnicalError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return r, errors.Join(ErrTechnicalError, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to close body, possible memory leak: %v", err))
		}
	}(resp.Body)

	if resp.StatusCode == 401 {
		return r, ErrInvalidToken
	}

	if resp.StatusCode >= 300 {
		// The jellyfin API gives no info for errors in API responses,
		// so all we can say is something went wrong.
		return r, ErrUnknownApiError
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&r)
	if err != nil {
		return r, errors.Join(ErrUnknownApiResponse, err)
	}

	return r, nil
}
