package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
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
	"github.com/icza/mjpeg"
)

var (
	port       = flag.Int("port", 8080, "Webserver port to listen on.")
	hostname   = flag.String("hostname", "localhost", "Host to listen on.")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func generateImage(params params.ImageParams, comp fractales.Computation) image.Image {
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

	return img
}

func generateVideo(w io.Writer, imageParams params.ImageParams) error {
	aw, err := mjpeg.New("test.avi", int32(imageParams.Width), int32(imageParams.Height), 25)

	if err != nil {
		return err
	}

	left := imageParams.Left
	right := imageParams.Right
	top := imageParams.Top
	bottom := imageParams.Bottom

	deltaX := math.Abs(left - right)
	deltaY := math.Abs(top - bottom)

	for  i := 0; i < 200; i++ {
		imageParams.Left = left + float64(i)  / 2 * deltaX / 200
		imageParams.Right = right - float64(i) / 2 * deltaX / 200
		imageParams.Top = top + float64(i) / 2 * deltaY / 200
		imageParams.Bottom = bottom - float64(i) / 2 * deltaY / 200
		comp := parsing.ComputerFromParameters(imageParams)
		img := generateImage(imageParams, comp)


		buf := &bytes.Buffer{}
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			return err
		}

		err = aw.AddFrame(buf.Bytes())
		if err != nil {
			return err
		}
	}

	err = aw.Close()
	if err != nil {
		return err
	}

	return nil
}

func fractale(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	imageParams := parsing.ParseImageParams(r)
	comp := parsing.ComputerFromParameters(imageParams)

	w.Header().Set("Content-Type", "image/png")
	img := generateImage(imageParams, comp)
	err := png.Encode(w, img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Print("Image served", imageParams)
	fmt.Printf("in %s\n", time.Since(start))
}

func video(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	imageParams := parsing.ParseImageParams(r)

	w.Header().Set("Content-Type", "image/png")
	err := generateVideo(w, imageParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Print("Video served", imageParams)
	fmt.Printf("in %s\n", time.Since(start))
}

func main() {
	flag.Parse()
	http.HandleFunc("/", fractale)
	http.HandleFunc("/video/", video)
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
