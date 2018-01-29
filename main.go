//
// main.go
// Copyright (C) 2018 pietro <pietro@the-arch>
//
// Distributed under terms of the MIT license.
//
package main

import (
	"github.com/gorilla/mux"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	"net/http"
	"os"
)

type ScaleRatio struct {
	Width  int
	Height int
}

var window *gtk.Window
var buf *gdkpixbuf.Pixbuf
var image *gtk.Image
var scale ScaleRatio

func main() {
	gtk.Init(nil)
	window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Go View")
	window.Connect("destroy", gtk.MainQuit)

	box := gtk.NewHBox(false, 0)

	image = gtk.NewImage()

	refresh()

	window.Connect("configure-event", resize)

	box.Add(image)
	window.Add(box)
	window.ShowAll()

	go server()

	gtk.Main()
}

func refresh() {
	buf, _ = gdkpixbuf.NewPixbufFromFile(os.Args[1])

	// Getting the scale
	scale = ScaleRatio{Height: buf.GetHeight(), Width: buf.GetWidth()}
	for i := 2; scale.IsLessThan(i); {
		if !scale.Divide(i) {
			scale = ScaleRatio{Width: scale.Width / i, Height: scale.Height / i}
			i++
		}
	}

	resize()
}

func resize() {

	w := window.GetAllocation().Width
	h := window.GetAllocation().Height

	smallest := 0
	if w/scale.Width > h/scale.Height {
		smallest = h / scale.Height
	} else {
		smallest = w / scale.Width
	}

	finalW := scale.Width * smallest
	finalH := scale.Height * smallest

	image.SetFromPixbuf(buf.ScaleSimple(finalW, finalH, gdkpixbuf.INTERP_BILINEAR))
}

func server() {

	r := mux.NewRouter()

	r.HandleFunc("/refresh", refreshReq)

	http.ListenAndServe(":6969", r)
}

func refreshReq(w http.ResponseWriter, r *http.Request) {
	refresh()
}

func (s *ScaleRatio) Divide(i int) bool {

	if s.Width%i != 0 || s.Height%i != 0 {
		return false
	}

	s.Width = s.Width / i
	s.Height = s.Height / i

	return true
}

func (s ScaleRatio) IsLessThan(i int) bool {
	return s.Width/i >= 2 && s.Height/i >= 2
}
