package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"os"
	"image"
	"strconv"
	"image/png"
	"bytes"
	"log"
	"github.com/go-chi/chi/middleware"
	"strings"
	"fmt"
)

func webserv() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/{genID}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "genID")
		g, exists := generators[id]
		if !exists {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "invalid id")

			return
		}
		fargs := make([]string, 0)
		for _,y := range r.URL.Query() {
			fargs = append(fargs, y...)
		}
		args := make([]string, 0)
		for _, z := range fargs {
			for _, v := range strings.Split(z, "|") {
				args = append(args, v)
			}
		}
		i, err := g.Render(args)
		if err != nil {
			panic(err)
		}
		writeImage(w, &i)
	})
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}