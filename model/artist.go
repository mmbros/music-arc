package model

import (
	"encoding/xml"
	"sort"
)

// ArtistMap type
//type ArtistMap map[string]*Artist

// Artists type
type Artists struct {
	Map  map[string]*Artist
	List []*Artist // ordered by SortName
}

// Artist type
type Artist struct {
	XMLName   xml.Name `xml:"artist"`
	ID        string   `xml:"id,attr"`
	Name      string   `xml:"name"`
	SortName  string   `xml:"sortName"`
	BeginDate string   `xml:"beginDate"`
	EndDate   string   `xml:"endDate"`
	Country   string   `xml:"country,attr"`
	Type      string   `xml:"type,attr"`

	Albums []*Album
}

/*
// ArtistRole type
type ArtistRole int

const (
	Default ArtistRole = iota
	Featuring
	With
)
*/

// ArtistRef type
type ArtistRef struct {
	ArtistID string `xml:",chardata"`
	Role     string `xml:"role,attr"`

	Artist *Artist
}

// Name returns Artist.Name
func (ar *ArtistRef) Name() string {
	return ar.Artist.Name
}

// RoleDefault returns true if artist role is `default` (main)
func (ar *ArtistRef) RoleDefault() bool {
	return len(ar.Role) == 0 || ar.Role == "default"
}

// RoleWith returns true if artist role is `with`
func (ar *ArtistRef) RoleWith() bool {
	return ar.Role == "with"
}

// RoleFeaturing returns true if artist role is `featuring`
func (ar *ArtistRef) RoleFeaturing() bool {
	return ar.Role == "featuring"
}

// UnmarshalXML function unmarshal an <artist-list> XML fragment to a Map[string]*Artist
func (am *Artists) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	/*
		     https://groups.google.com/forum/#!topic/golang-nuts/Y6KitymP7Dg
			   https://gist.github.com/chrisfarms/1377218
			   http://play.golang.org/p/v9Nt57Fj03
	*/

	// Artists.Map
	m := map[string]*Artist{}

LOOP:
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:

			if t.Name.Local == "artist" {
				a := Artist{}
				elmt := xml.StartElement(t)
				d.DecodeElement(&a, &elmt)
				m[a.ID] = &a
			}
		case xml.EndElement:
			if t.Name.Local == "artist-list" {
				break LOOP
			}
		}
	}

	// Artists.List
	l := make([]*Artist, len(m))
	var j int
	for _, a := range m {
		l[j] = a
		j++
	}
	sort.Sort(bySortName(l))

	// result
	*am = Artists{Map: m, List: l}
	return nil
}

type bySortName []*Artist

func (sn bySortName) Len() int           { return len(sn) }
func (sn bySortName) Swap(i, j int)      { sn[i], sn[j] = sn[j], sn[i] }
func (sn bySortName) Less(i, j int) bool { return sn[i].SortName < sn[j].SortName }
