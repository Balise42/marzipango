package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/Balise42/marzipango/fractales"
	"github.com/Balise42/marzipango/params"
	"github.com/Balise42/marzipango/parsing"
)

var (
	port       = flag.Int("port", 8080, "Webserver port to listen on.")
	hostname   = flag.String("hostname", "localhost", "Host to listen on.")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
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
	start := time.Now()
	comp, imageParams := parsing.ParseComputation(r)

	w.Header().Set("Content-Type", "image/png")
	err := generateImage(w, imageParams, comp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Print("Image served", imageParams)
	fmt.Printf("in %s\n", time.Since(start))
}

func main() {
	flag.Parse()
	http.HandleFunc("/", fractale)
	address := fmt.Sprintf("%s:%d", *hostname, *port)
	fmt.Printf("Listening on http://%s ...\n", address)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pprof.StopCPUProfile()
			os.Exit(1)
		}
	}()

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
