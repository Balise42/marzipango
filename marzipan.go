package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/Balise42/marzipango/fractales"
	"github.com/Balise42/marzipango/params"
	"github.com/Balise42/marzipango/parsing"
)

func generateImage(w io.Writer, params params.ImageParams, comp fractales.Computation) error {
	var wg sync.WaitGroup
	img := image.NewRGBA64(image.Rect(0, 0, params.Width, params.Height))
	for x := 0; x < params.Width; x++ {
		var numRows = params.Height / runtime.NumCPU()
		for cpu := 0; cpu < runtime.NumCPU()-1; cpu++ {
			wg.Add(1)
			go comp(x, numRows*cpu, numRows*(cpu+1), img, &wg)
		}
		wg.Add(1)
		go comp(x, numRows*(runtime.NumCPU()-1), params.Height, img, &wg)
	}
	wg.Wait()

	return png.Encode(w, img)
}

func fractale(w http.ResponseWriter, r *http.Request) {
	comp, imageParams := parsing.ParseComputation(r)

	w.Header().Set("Content-Type", "image/png")
	err := generateImage(w, imageParams, comp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Image served", imageParams)
}

func main() {
	http.HandleFunc("/", fractale)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
