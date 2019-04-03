package gpmdp

func (g *GPMDP) pushRating(rating *Rating) {
	g.ratingCBsMtx.RLock()
	defer g.ratingCBsMtx.RUnlock()

	for _, cb := range g.ratingCBs {
		cb(rating)
	}
}

func (g *GPMDP) Rating() chan *Rating {
	ratings := make(chan *Rating, 10)
	go func() {
		g.ratingCBsMtx.Lock()
		defer g.ratingCBsMtx.Unlock()

		g.ratingCBs = append(g.ratingCBs, func(rating *Rating) {
			ratings <- rating
		})
	}()
	return ratings
}
