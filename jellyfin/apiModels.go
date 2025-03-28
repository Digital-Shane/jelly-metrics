package jellyfin

import "time"

type session struct {
	LastActivityDate    time.Time `json:"LastActivityDate"`
	LastPlaybackCheckIn time.Time `json:"LastPlaybackCheckIn"`
	NowPlayingItem      *struct {
		ID string `json:"Id"`
	} `json:"NowPlayingItem,omitempty"`
	PlayState struct {
		IsPaused bool `json:"IsPaused"`
	} `json:"PlayState"`
	UserName string `json:"UserName"`
}

type mediaCounts struct {
	AlbumCount      int `json:"AlbumCount"`
	ArtistCount     int `json:"ArtistCount"`
	BookCount       int `json:"BookCount"`
	BoxSetCount     int `json:"BoxSetCount"`
	EpisodeCount    int `json:"EpisodeCount"`
	ItemCount       int `json:"ItemCount"`
	MovieCount      int `json:"MovieCount"`
	MusicVideoCount int `json:"MusicVideoCount"`
	ProgramCount    int `json:"ProgramCount"`
	SeriesCount     int `json:"SeriesCount"`
	SongCount       int `json:"SongCount"`
	TrailerCount    int `json:"TrailerCount"`
}
