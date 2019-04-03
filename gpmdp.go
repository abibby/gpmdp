package gpmdp

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// ws://localhost:5672

type Track struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	AlbumArt string `json:"albumArt"`
}

type Rating struct {
	Liked    bool `json:"liked"`
	Disliked bool `json:"disliked"`
}

type Time struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type Info struct {
	Playing    bool    `json:"playing"`
	Song       *Track  `json:"song"`
	Rating     *Rating `json:"rating"`
	Time       *Time   `json:"time"`
	SongLyrics string  `json:"songLyrics"`
	Shuffle    string  `json:"shuffle"`
	Repeat     string  `json:"repeat"`
	Volume     int     `json:"volume"`
}

type message struct {
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
}

// https://github.com/MarshallOfSound/Google-Play-Music-Desktop-Player-UNOFFICIAL-/blob/master/docs/PlaybackAPI_WebSocket.md

//go:generate go run gen.go error Error error

//go:generate go run gen.go bool Playing playing
//go:generate go run gen.go *Track Track track
//go:generate go run gen.go *Rating Rating rating
//go:generate go run gen.go *Time Time time

type GPMDP struct {
	errorCBs    []func(error)
	errorCBsMtx sync.RWMutex

	playingCBs    []func(bool)
	playingCBsMtx sync.RWMutex

	trackCBs    []func(*Track)
	trackCBsMtx sync.RWMutex

	ratingCBs    []func(*Rating)
	ratingCBsMtx sync.RWMutex

	timeCBs    []func(*Time)
	timeCBsMtx sync.RWMutex
}

func Connect() (*GPMDP, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:5672", nil)
	if err != nil {
		return nil, err
	}
	msg := &message{}
	g := &GPMDP{}
	go func() {
		for {
			err := conn.ReadJSON(msg)
			if err != nil {
				g.pushError(err)
			}

			switch msg.Channel {
			case "track":
				track := &Track{}
				err := json.Unmarshal(msg.Payload, track)
				if err != nil {
					g.pushError(err)
				}
				g.pushTrack(track)
			case "playState":
				var playing bool
				err := json.Unmarshal(msg.Payload, &playing)
				if err != nil {
					g.pushError(err)
				}
				g.pushPlaying(playing)
			}
		}
	}()
	return g, nil
}
