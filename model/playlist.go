package model

import "encoding/xml"

// PlaylistMap type
type PlaylistMap map[string]*Playlist

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
}

// UnmarshalXML function unmarshal an <album-list> XML fragment to a Map[string]*Album
func (pm *PlaylistMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

	*pm = PlaylistMap(m)
	return nil
}
