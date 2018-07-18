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

func main() {
	gtk.Init(nil)
	window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Go View")
	window.Connect("destroy", gtk.MainQuit)

	box := gtk.NewHBox(false, 0)

	image = gtk.NewImage()

	refresh()

	window.Connect("configure-event", refresh)

	box.Add(image)
	window.Add(box)
	window.ShowAll()

	go server()

	gtk.Main()
}

func refresh() {

	file, _ := imaging.Open(os.Args[1])

	scaled := imaging.Fit(file, window.GetAllocation().Width, window.GetAllocation().Height, imaging.Box)

	var bts []byte

	buffer := bytes.NewBuffer(bts)

	imaging.Encode(buffer, scaled, imaging.PNG)

	buf, _ := gdkpixbuf.NewPixbufFromBytes(buffer.Bytes())

	image.SetFromPixbuf(buf)
}

func server() {

	r := mux.NewRouter()

	r.HandleFunc("/refresh", refreshReq)

	http.ListenAndServe(":6969", r)
}

func refreshReq(w http.ResponseWriter, r *http.Request) {
	glib.IdleAdd(refresh)
}
