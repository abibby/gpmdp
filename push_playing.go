package gpmdp

func (g *GPMDP) pushPlaying(playing bool) {
	g.playingCBsMtx.RLock()
	defer g.playingCBsMtx.RUnlock()

	for _, cb := range g.playingCBs {
		cb(playing)
	}
}

func (g *GPMDP) Playing() chan bool {
	playings := make(chan bool, 10)
	go func() {
		g.playingCBsMtx.Lock()
		defer g.playingCBsMtx.Unlock()

		g.playingCBs = append(g.playingCBs, func(playing bool) {
			playings <- playing
		})
	}()
	return playings
}
