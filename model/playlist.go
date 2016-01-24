package model

import (
	"encoding/xml"
	"sort"
)

// Playlists type
type Playlists struct {
	Map  map[string]*Playlist
	List []*Playlist // ordered by ID
}

// PlaylistMap type
// type PlaylistMap map[string]*Playlist

// Playlist type
type Playlist struct {
	XMLName  xml.Name `xml:"playlist"`
	ID       string   `xml:"id,attr"`
	Type     string   `xml:"type,attr"`
	Title    string   `xml:"title"`
	Creator  string   `xml:"creator"`
	Date     string   `xml:"date"`
	Duration string   `xml:"duration"`
	Style    string   `xml:"style,attr"`
	AlbumIDs []string `xml:"list>album-ref"`

	Albums []*Album
}

// UnmarshalXML function unmarshal an <album-list> XML fragment to a Map[string]*Album
func (ps *Playlists) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// result
	m := map[string]*Playlist{}

LOOP:
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "playlist" {
				p := Playlist{}
				elmt := xml.StartElement(t)
				d.DecodeElement(&p, &elmt)
				m[p.ID] = &p
			}
		case xml.EndElement:
			if t.Name.Local == "playlist-list" {
				break LOOP
			}
		}
	}

	// Playlists.List
	l := make([]*Playlist, len(m))
	var j int
	for _, a := range m {
		l[j] = a
		j++
	}
	sort.Sort(byID(l))

	// result
	*ps = Playlists{Map: m, List: l}
	return nil

}

type byID []*Playlist

func (ps byID) Len() int           { return len(ps) }
func (ps byID) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps byID) Less(i, j int) bool { return ps[i].ID < ps[j].ID }

// Year returns the year of the album
func (p *Playlist) Year() string {
	if len(p.Date) < 4 {
		return ""
	}
	return p.Date[:4]
}
