package gpmdp

func (g *GPMDP) pushError(error error) {
	g.errorCBsMtx.RLock()
	defer g.errorCBsMtx.RUnlock()

	for _, cb := range g.errorCBs {
		cb(error)
	}
}

func (g *GPMDP) Error() chan error {
	errors := make(chan error, 10)
	go func() {
		g.errorCBsMtx.Lock()
		defer g.errorCBsMtx.Unlock()

		g.errorCBs = append(g.errorCBs, func(error error) {
			errors <- error
		})
	}()
	return errors
}
