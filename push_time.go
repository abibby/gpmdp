package gpmdp

func (g *GPMDP) pushTime(time *Time) {
	g.timeCBsMtx.RLock()
	defer g.timeCBsMtx.RUnlock()

	for _, cb := range g.timeCBs {
		cb(time)
	}
}

func (g *GPMDP) Time() chan *Time {
	times := make(chan *Time, 10)
	go func() {
		g.timeCBsMtx.Lock()
		defer g.timeCBsMtx.Unlock()

		g.timeCBs = append(g.timeCBs, func(time *Time) {
			times <- time
		})
	}()
	return times
}
