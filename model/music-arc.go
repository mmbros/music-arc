package model

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// MusicArc type
type MusicArc struct {
	XMLName   xml.Name  `xml:"music-arc"`
	Artists   Artists   `xml:"artist-list"`
	Albums    AlbumMap  `xml:"album-list"`
	Playlists Playlists `xml:"playlist-list"`
}

var entities = map[string]string{
	"sep":    "|",      // bullet = black small circle
	"bull":   "\u2022", // bullet = black small circle
	"hellip": "\u2026", // horizontal ellipsis = three dot leader
	"nbsp":   "\u00A0", // non-breaking space
}

// LoadFromXMLFile function
func LoadFromXMLFile(path string) (*MusicArc, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	d := xml.NewDecoder(bytes.NewReader(xmlData))
	d.Entity = entities

	ma := MusicArc{}

	if err := d.Decode(&ma); err != nil {
		return nil, err
	}

	if err := ma.postLoad(); err != nil {
		return nil, err
	}

	return &ma, nil
}

// UnmarshalXML function unmarshal an <album-list> XML fragment to a Map[string]*Album
func (ma *MusicArc) postLoad() error {

	// update each Album.ArtistRefs and Artist.Albums
	for _, album := range ma.Albums {
		for _, artistRef := range album.ArtistRefs {
			artist, ok := ma.Artists.Map[artistRef.ArtistID]
			if !ok {
				return fmt.Errorf("Invalid ArtistRef \"%s\" in album \"%s\"\n", artistRef.ArtistID, album.ID)
			}
			// update artistRef.Artist
			artistRef.Artist = artist
			// add album to artist.Albums
			artist.Albums = append(artist.Albums, album)
		}
	}

	// sort by Date artist.Albums
	for _, artist := range ma.Artists.List {
		sort.Sort(byDate(artist.Albums))
	}

	// init each Playlist.Albums
	for _, playlist := range ma.Playlists.List {
		for _, albumID := range playlist.AlbumIDs {
			// get album by id
			album, ok := ma.Albums[albumID]
			if !ok {
				return fmt.Errorf("Invalid album \"%s\" in playlist \"%s\"\n", albumID, playlist.ID)
			}
			// add album to playlist.Albums
			playlist.Albums = append(playlist.Albums, album)
		}
	}

	return nil
}

type byDate []*Album

func (s byDate) Len() int      { return len(s) }
func (s byDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byDate) Less(i, j int) bool {
	return (s[i].Date < s[j].Date) || ((s[i].Date == s[j].Date) && (s[i].ID < s[j].ID))
}
