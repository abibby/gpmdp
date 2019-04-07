package gpmdp

import (
	"encoding/json"

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

type rawEvent struct {
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
}
type Event struct {
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
}

func (e *Event) Track() (track *Track, ok bool) {
	if e.Channel != "track" {
		return nil, false
	}
	track, ok = e.Payload.(*Track)
	return track, ok
}

func (e *Event) Playing() (playing bool, ok bool) {
	if e.Channel != "playState" {
		return false, false
	}
	playing, ok = e.Payload.(bool)
	return playing, ok
}

// https://github.com/MarshallOfSound/Google-Play-Music-Desktop-Player-UNOFFICIAL-/blob/master/docs/PlaybackAPI_WebSocket.md

type GPMDP struct {
	Error chan error
	Event chan *Event

	open bool
}

func Connect() (*GPMDP, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:5672", nil)
	if err != nil {
		return nil, err
	}
	msg := &rawEvent{}
	g := &GPMDP{
		Error: make(chan error),
		Event: make(chan *Event),
		open:  true,
	}
	go func() {
		for g.open {
			err := conn.ReadJSON(msg)
			if err != nil {
				go g.pushError(err)
				continue
			}

			var payload interface{}
			switch msg.Channel {
			case "track":
				track := &Track{}
				err := json.Unmarshal(msg.Payload, track)
				if err != nil {
					go g.pushError(err)
					continue
				}
				payload = track
			case "playState":
				var playing bool
				err := json.Unmarshal(msg.Payload, &playing)
				if err != nil {
					go g.pushError(err)
					continue
				}
				payload = playing
			default:
				continue
			}
			g.Event <- &Event{
				Channel: msg.Channel,
				Payload: payload,
			}
		}
	}()
	return g, nil
}

func (g *GPMDP) Close() {
	g.open = false
}

func (g *GPMDP) pushError(err error) {
	g.Error <- err
}
