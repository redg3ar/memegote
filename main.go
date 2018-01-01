package main

import (
	"io/ioutil"
	"encoding/json"
	"image"
	"strings"
	"errors"
	"image/jpeg"
	"os"
	"image/png"
)

type Generator struct {
	Name string
	Labels []Label
}
type Label struct {
	X1, X2, Y1, Y2 uint
	FontFamily string
}

var generators map[string]Generator

func main() {
	generators = make(map[string]Generator)

	//Start with jsons
	{
		jsons, err := ioutil.ReadDir("json")
		if err != nil {
			panic(err)
		}
		for _,x := range jsons {
			b, err := ioutil.ReadFile("json/"+x.Name())
			if err != nil {
				panic(err)
			}
			newGen := Generator{}
			json.Unmarshal(b, &newGen)
			generators[newGen.Name] = newGen
		}
	}
}

func (g *Generator) Render(args... string) (image.Image, error) {

	//Phase one: grab image
	var i image.Image
	{
		//Find out extension
		var ext string
		if n := strings.Split(g.Name, "."); len(n) != 2 {
			ext = n[1]
		} else {
			return nil, errors.New("Invalid generator file")
		}
		//Prepare file for reading
		b, err := os.Open("images/" + g.Name)
		if err != nil {
			panic(err)
		}

		//do it
		switch ext {
		case "jpeg", "jpg":
			i, err = jpeg.Decode(b)
			if err != nil {
				return nil, err
			}
		case "png":
			i, err = png.Decode(b)
			//HOOEY I LOVE GOLANG'S BOILERPLATE FREE WRIVDSICOAS
			if err != nil {
				return nil, err
			}
		}
	}

	//Phase two: add labels
	{
		for x := range g.Labels {
			if (len(args) - 1) < x {

			}
		}
	}
}