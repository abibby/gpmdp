package gpmdp

func (g *GPMDP) pushTrack(track *Track) {
	g.trackCBsMtx.RLock()
	defer g.trackCBsMtx.RUnlock()

	for _, cb := range g.trackCBs {
		cb(track)
	}
}

func (g *GPMDP) Track() chan *Track {
	tracks := make(chan *Track, 10)
	go func() {
		g.trackCBsMtx.Lock()
		defer g.trackCBsMtx.Unlock()

		g.trackCBs = append(g.trackCBs, func(track *Track) {
			tracks <- track
		})
	}()
	return tracks
}
