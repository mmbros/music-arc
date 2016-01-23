package model

import (
	"bytes"
	"encoding/xml"
)

// AlbumMap type
type AlbumMap map[string]*Album

/*
type AlbumMap struct {
	Map  map[string]*Album
	List []*Album  // order by name
}
*/

// Album type
type Album struct {
	XMLName       xml.Name     `xml:"album"`
	ID            string       `xml:"id,attr"`
	Title         string       `xml:"title"`
	ArtistDisplay string       `xml:"artist-display"`
	ArtistRefs    []*ArtistRef `xml:"artist-ref"`
	Date          string       `xml:"date"`
	Duration      string       `xml:"duration"`
	Imgs          []string     `xml:"img"`
	Tracklists    []Tracklist  `xml:"track-list"`
}

// Tracklist type
type Tracklist struct {
	XMLName       xml.Name     `xml:"track-list"`
	Title         string       `xml:"title"`
	ArtistDisplay string       `xml:"artist-display"`
	ArtistRefs    []*ArtistRef `xml:"artist-ref"`
	Duration      string       `xml:"duration"`
	Tracks        []Track      `xml:"track"`
	ShowArtist    bool         `xml:"show-artist,attr"`
}

// Track type
type Track struct {
	XMLName       xml.Name     `xml:"track"`
	Title         string       `xml:"title"`
	ArtistDisplay string       `xml:"artist-display"`
	ArtistRefs    []*ArtistRef `xml:"artist-ref"`
	Duration      string       `xml:"duration"`
}

// UnmarshalXML function unmarshal an <album-list> XML fragment to a Map[string]*Album
func (am *AlbumMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// result
	m := map[string]*Album{}

LOOP:
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:

			if t.Name.Local == "album" {
				a := Album{}
				elmt := xml.StartElement(t)
				d.DecodeElement(&a, &elmt)
				m[a.ID] = &a
			}
		case xml.EndElement:
			if t.Name.Local == "album-list" {
				break LOOP
			}
		}
	}

	*am = AlbumMap(m)
	return nil
}

// Year returns the year of the album
func (a *Album) Year() string {
	if len(a.Date) < 4 {
		return ""
	}
	return a.Date[:4]
}

// ArtistName returns the artist name of the album
func (a *Album) ArtistName() string {
	if len(a.ArtistDisplay) > 0 {
		return a.ArtistDisplay
	}

	var buffer bytes.Buffer

	for i, ar := range a.ArtistRefs {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(ar.Artist.Name)
	}

	return buffer.String()
}
