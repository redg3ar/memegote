package main

import (
	"encoding/json"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image/draw"
	"image/color"
	"golang.org/x/image/math/fixed"
	"github.com/davecgh/go-spew/spew"
	"log"
)

type Generator struct {
	Name   string
	Labels []Label
}
type Label struct {
	Point image.Point
	Size float64
	Color color.RGBA
	FontFamily     string
}

var cac *cache.Cache

var generators map[string]Generator

func main() {
	generators = make(map[string]Generator)

	cac = cache.New(5*time.Minute, 10*time.Minute)

	//Start with jsons
	{
		jsons, err := ioutil.ReadDir("json")
		if err != nil {
			panic(err)
		}
		for _, x := range jsons {
			b, err := ioutil.ReadFile("json/" + x.Name())
			if err != nil {
				panic(err)
			}
			newGen := Generator{}
			json.Unmarshal(b, &newGen)
			generators[newGen.Name] = newGen
		}
	}
	webserv()
}

func (g *Generator) Render(args []string) (image.Image, error) {

	//Phase one: grab image

	var i image.Image
	{
		//Prepare file for reading
		b, err := os.Open("images/" + g.Name)
		if err != nil {
			panic(err)
		}

		//do it
		i, _, err = image.Decode(b)
		if err != nil {
			panic(err)
		}
	}

	//Phase two: add labels

	for x := range g.Labels {

			interf, found := cac.Get(g.Labels[x].FontFamily)
			if !found {
				file, err := os.Open("fonts/" + g.Labels[x].FontFamily + ".ttf")
				if err != nil {
					panic(err)
				}
				b, err := ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}
				interf, err = truetype.Parse(b)
				if err != nil {
					panic(err)
				}
				cac.Set(g.Labels[x].FontFamily, interf, cache.DefaultExpiration)
			}
			f := truetype.NewFace(interf.(*truetype.Font), &truetype.Options{Size: g.Labels[x].Size})
			drawer := font.Drawer{
				Dst:  i.(draw.Image),
				Src:  image.NewUniform(g.Labels[x].Color),
				Face: f,
				Dot:  fixed.P(g.Labels[x].Point.X, g.Labels[x].Point.Y),
			}
			drawer.DrawString(args[x])
	}

	return i, nil
}
