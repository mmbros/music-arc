//go:generate go run gen-templates.go

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mmbros/music-arc/model"
	"github.com/mmbros/music-arc/templates"
)

const urlPrefixAlbums = "/albums/"
const urlPrefixArtists = "/artists/"

var gMA *model.MusicArc

func viewAlbumHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s\n", r.Method, r.URL.Path)

	id := r.URL.Path[len(urlPrefixAlbums):]

	if len(id) == 0 {
		p := templates.PageAlbumList
		if err := p.Execute(w, &gMA); err != nil {
			panic(err)
		}
		return
	}

	album := gMA.Albums[id]
	if album == nil {
		fmt.Fprintf(w, "Album \"%s\" non trovato!", id)
		return
	}

	p := templates.PageAlbum
	if err := p.Execute(w, &album); err != nil {
		panic(err)
	}

}

func viewArtistHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s\n", r.Method, r.URL.Path)

	id := r.URL.Path[len(urlPrefixArtists):]

	if len(id) == 0 {
		p := templates.PageArtistList
		if err := p.Execute(w, &gMA); err != nil {
			panic(err)
		}
		return
	}

	artist := gMA.Artists.Map[id]
	if artist == nil {
		fmt.Fprintf(w, "Artist \"%s\" non trovato!", id)
		return
	}

	p := templates.PageArtist
	if err := p.Execute(w, &artist); err != nil {
		panic(err)
	}
}

func create() {
	arc := "./data/music-arc-inc.xml"
	err := model.CreateMusicArcInc("./data", arc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Creating %s ok\n", arc)
}

func main() {
	var err error

	//	create()

	gMA, err = model.LoadFromXMLFile("data/music-arc-inc.xml")
	if err != nil {
		panic(err)
	}

	addr := ":8080"
	log.Printf("listening to %s", addr)

	//	http.HandleFunc("/", handler)
	http.HandleFunc(urlPrefixAlbums, viewAlbumHandler)
	http.HandleFunc(urlPrefixArtists, viewArtistHandler)

	http.Handle("/img/album/", http.StripPrefix("/img/album/", http.FileServer(http.Dir("./data/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./dist/css"))))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panic(err)
	}
}