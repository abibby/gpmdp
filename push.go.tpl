package gpmdp

func (g *GPMDP) push{{.FuncName}}({{.Name}} {{.Type}}) {
	g.{{.Name}}CBsMtx.RLock()
	defer g.{{.Name}}CBsMtx.RUnlock()

	for _, cb := range g.{{.Name}}CBs {
		cb({{.Name}})
	}
}

func (g *GPMDP) {{.FuncName}}() chan {{.Type}} {
	{{.Name}}s := make(chan {{.Type}}, 10)
	go func() {
		g.{{.Name}}CBsMtx.Lock()
		defer g.{{.Name}}CBsMtx.Unlock()

		g.{{.Name}}CBs = append(g.{{.Name}}CBs, func({{.Name}} {{.Type}}) {
			{{.Name}}s <- {{.Name}}
		})
	}()
	return {{.Name}}s
}
