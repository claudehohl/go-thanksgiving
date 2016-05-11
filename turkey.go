package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	backgroundImageFile = "img/bg/thanksgiving-2011-actual-size-bg-nolinks.png"
	defaultImageFile    = "img/00000000.png"
	prefix              = "/thumb/"
)

var (
	// dirs maps each layout element to its location on disk.
	dirs = map[string]string{
		"h": "img/heads",
		"b": "img/eyes_beak",
		"i": "img/index_feathers",
		"m": "img/middle_feathers",
		"r": "img/ring_feathers",
		"p": "img/pinky_feathers",
		"f": "img/feet",
		"w": "img/wing",
	}

	// urlMap maps each URL character position to
	// its corresponding layout element.
	urlMap = [...]string{"b", "h", "i", "m", "r", "p", "f", "w"}

	// layoutMap maps each layout element to its position
	// on the background image.
	layoutMap = map[string]image.Rectangle{
		"h": {image.Pt(109, 50), image.Pt(166, 152)},
		"i": {image.Pt(136, 21), image.Pt(180, 131)},
		"m": {image.Pt(159, 7), image.Pt(201, 126)},
		"r": {image.Pt(188, 20), image.Pt(230, 125)},
		"p": {image.Pt(216, 48), image.Pt(258, 134)},
		"f": {image.Pt(155, 176), image.Pt(243, 213)},
		"w": {image.Pt(169, 118), image.Pt(250, 197)},
		"b": {image.Pt(105, 104), image.Pt(145, 148)},
	}

	// elements maps each layout element to its images.
	elements = make(map[string][]*image.RGBA)

	// backgroundImage contains the background image data.
	backgroundImage *image.RGBA

	// defaultImage is the image that is served if an error occurs.
	defaultImage *image.RGBA

	// imageQuality is the encoding quality setting for the output image.
	imageQuality = jpeg.Options{95}
)

func main() {
	fmt.Printf("Server running at 8080\n")
	http.HandleFunc(prefix, handler)
	http.ListenAndServe(":8080", nil)
}

// handler serves a turkey snapshot for the given request.
func handler(w http.ResponseWriter, r *http.Request) {

	// Load images from disk on the first request.
	loadOnce.Do(load)

	// Make a copy of the background to draw into.
	bgRect := backgroundImage.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, bgRect.Dx(), bgRect.Dy()))
	draw.Draw(m, m.Bounds(), backgroundImage, image.ZP, draw.Over)

	// Process each character of the request string.
	code := strings.ToLower(r.URL.Path[len(prefix):])
	for i, p := range code {
		// Decode hex character p in place.
		if p < 'a' {
			// it's a digit
			p = p - '0'
		} else {
			// it's a letter
			p = p - 'a' + 10
		}

		t := urlMap[i]    // element type by index
		em := elements[t] // element images by type
		if int(p) >= len(em) {
			panic(fmt.Sprintf("element index out of range %s: "+
				"%d >= %d", t, p, len(em)))
		}

		// Draw the element to m,
		// using the layoutMap to specify its position.
		draw.Draw(m, layoutMap[t], em[p], image.ZP, draw.Over)
	}

	// Encode JPEG image and write it as the response.
	w.Header().Set("Content-type", "image/jpeg")
	w.Header().Set("cache-control", "public, max-age=259200")
	jpeg.Encode(w, m, &imageQuality)
	fmt.Printf("Served a png\n")
}

// loadOnce is used to call the load function only on the first request.
var loadOnce sync.Once

// load reads the various PNG images from disk and stores them in their
// corresponding global variables.
func load() {
	defaultImage = loadPNG(defaultImageFile)
	backgroundImage = loadPNG(backgroundImageFile)
	for dirKey, dir := range dirs {
		paths, err := filepath.Glob(dir + "/*.png")
		if err != nil {
			panic(err)
		}
		for _, p := range paths {
			elements[dirKey] = append(elements[dirKey], loadPNG(p))
		}
	}
}

// loadPNG loads a PNG image from disk and returns it as an *image.RGBA.
func loadPNG(filename string) *image.RGBA {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	m, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return rgba(m)
}

// rgba converts an image.Image to an *image.RGBA.
func rgba(m image.Image) *image.RGBA {
	if r, ok := m.(*image.RGBA); ok {
		return r
	}
	b := m.Bounds()
	r := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(r, b, m, image.ZP, draw.Src)
	return r
}
