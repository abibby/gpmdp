// +build !

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type data struct {
	Type     string
	FuncName string
	Name     string
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 4 {
		check(fmt.Errorf("must have 4 arguments"))
	}

	typeName := os.Args[1]
	funcName := os.Args[2]
	name := os.Args[3]

	b, err := ioutil.ReadFile("push.go.tpl")
	check(err)
	tmpl, err := template.New("push").Parse(string(b))
	check(err)
	fName := fmt.Sprintf("push_%s.go", name)
	os.Remove(fName)
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_RDWR, 0644)
	check(err)

	tmpl.Execute(f, data{
		Type:     typeName,
		FuncName: funcName,
		Name:     name,
	})

}
