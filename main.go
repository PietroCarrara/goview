//
// main.go
// Copyright (C) 2018 pietro <pietro@the-arch>
//
// Distributed under terms of the MIT license.
//
package main

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"net/http"
	"os"
)

var window *gtk.Window
var buf *gdkpixbuf.Pixbuf
var image *gtk.Image

var path string

func main() {
	gtk.Init(nil)
	window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Go View")
	window.Connect("destroy", gtk.MainQuit)

	box := gtk.NewHBox(false, 0)

	image = gtk.NewImage()

	if len(os.Args) > 1 {
		setImage(os.Args[1])
	}

	window.Connect("configure-event", refresh)

	box.Add(image)
	window.Add(box)
	window.ShowAll()

	go server()

	gtk.Main()
}

func refresh() {
	setImage(path)
}

func setImage(img string) {

	if img == "" {
		return
	}

	path = img

	if buf != nil {
		buf.Unref()
	}

	file, _ := imaging.Open(img)

	scaled := imaging.Fit(file, window.GetAllocation().Width, window.GetAllocation().Height, imaging.Box)

	var bts []byte

	buffer := bytes.NewBuffer(bts)

	imaging.Encode(buffer, scaled, imaging.PNG)

	buf, _ = gdkpixbuf.NewPixbufFromBytes(buffer.Bytes())

	image.SetFromPixbuf(buf)
}

func server() {

	r := mux.NewRouter()

	r.HandleFunc("/refresh", refreshReq)
	r.HandleFunc("/setImage", setImageReq)

	http.ListenAndServe(":6969", r)
}

func refreshReq(w http.ResponseWriter, r *http.Request) {
	glib.IdleAdd(refresh)
}

func setImageReq(w http.ResponseWriter, r *http.Request) {
	img := r.PostFormValue("image")

	glib.IdleAdd(func() { setImage(img) })
}
