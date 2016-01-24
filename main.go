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
const urlPrefixPlaylists = "/playlists/"
const urlPrefixFrontCover = "/front-cover/"

var gMA *model.MusicArc

// PageData is
type PageData struct {
	MusicArc *model.MusicArc
	Artist   *model.Artist
	Album    *model.Album
	Playlist *model.Playlist
}

func viewFrontCoverHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s\n", r.Method, r.URL.Path)

	// init page data
	data := PageData{MusicArc: gMA}

	// get the id
	id := r.URL.Path[len(urlPrefixFrontCover):]

	if len(id) == 0 {
		// error page
		fmt.Fprintf(w, "Must have a Playlist!")
		return
	}

	// get playlist by id
	playlist := gMA.Playlists.Map[id]
	if playlist == nil {
		// error page
		fmt.Fprintf(w, "Playlist \"%s\" not found!", id)
		return
	}
	data.Playlist = playlist

	// return playlist page detail
	p := templates.PageCoverFront10
	if err := p.Execute(w, &data); err != nil {
		panic(err)
	}

}

func viewPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s\n", r.Method, r.URL.Path)

	// init page data
	data := PageData{MusicArc: gMA}

	// get the id
	id := r.URL.Path[len(urlPrefixPlaylists):]

	if len(id) == 0 {
		// return list of playlists
		p := templates.PagePlaylistsList
		if err := p.Execute(w, &data); err != nil {
			panic(err)
		}
		return
	}

	// get playlist by id
	playlist := gMA.Playlists.Map[id]
	if playlist == nil {
		// error page
		fmt.Fprintf(w, "Playlist \"%s\" not found!", id)
		return
	}
	data.Playlist = playlist

	// return playlist page detail
	p := templates.PagePlaylist
	if err := p.Execute(w, &data); err != nil {
		panic(err)
	}

}

func viewAlbumHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s\n", r.Method, r.URL.Path)

	data := PageData{MusicArc: gMA}

	id := r.URL.Path[len(urlPrefixAlbums):]

	if len(id) == 0 {
		p := templates.PageAlbumList
		if err := p.Execute(w, &data.MusicArc); err != nil {
			panic(err)
		}
		return
	}

	album := gMA.Albums[id]
	if album == nil {
		fmt.Fprintf(w, "Album \"%s\" non trovato!", id)
		return
	}

	data.Album = album
	data.Artist = album.ArtistRefs[0].Artist

	p := templates.PageAlbum
	if err := p.Execute(w, &data); err != nil {
		panic(err)
	}

}

func viewArtistHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s\n", r.Method, r.URL.Path)

	data := PageData{MusicArc: gMA}

	id := r.URL.Path[len(urlPrefixArtists):]

	if len(id) == 0 {
		p := templates.PageArtistsList
		if err := p.Execute(w, &data); err != nil {
			panic(err)
		}
		return
	}

	artist := gMA.Artists.Map[id]
	if artist == nil {
		fmt.Fprintf(w, "Artist \"%s\" non trovato!", id)
		return
	}

	data.Artist = artist

	p := templates.PageArtist
	if err := p.Execute(w, &data); err != nil {
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
	http.HandleFunc(urlPrefixPlaylists, viewPlaylistHandler)
	http.HandleFunc(urlPrefixFrontCover, viewFrontCoverHandler)

	http.Handle("/img/album/", http.StripPrefix("/img/album/", http.FileServer(http.Dir("./data/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./dist/css"))))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panic(err)
	}
}
