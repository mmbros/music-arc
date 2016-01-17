package model

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CreateMusicArcInc function
func CreateMusicArcInc(srcDir, destPath string) error {

	fo, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fo.Close()

	// make a write buffer
	w := bufio.NewWriter(fo)

	if err := appendString(`
<!DOCTYPE music-arc [
  <!ENTITY sep    "&#x2022;"><!-- Bullet -->
	<!ENTITY bull   "&#x2022;"><!-- Bullet -->
  <!ENTITY hellip "&#x2026;"><!-- Horizontal ellipsis -->
  <!ENTITY nbsp     "&#xA0;"><!-- Non-breaking space -->
]>
<music-arc>
`, w, true, false); err != nil {
		return err
	}

	// artist-list
	pat := filepath.Join(srcDir, "artist/*.xml")
	if err := appendFiles("artist-list", pat, w); err != nil {
		return err
	}

	// album-list
	pat = filepath.Join(srcDir, "album/**/*.xml")
	if err := appendFiles(`album-list`, pat, w); err != nil {
		return err
	}
	// playlist-list
	pat = filepath.Join(srcDir, "playlist/*.xml")
	if err := appendFiles(`playlist-list`, pat, w); err != nil {
		return err
	}

	if err := appendString(`
</music-arc>
`, w, false, false); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func appendString(text string, w *bufio.Writer, insertXMLDecl bool, stripXMLDecl bool) error {
	const searchStart = `<?xml`
	const searchEnd = `?>`
	const xmlDecl = `<?xml version="1.0" encoding="UTF-8" ?>`

	a := -1

	// strip old XML declaration
	if stripXMLDecl || insertXMLDecl {
		a = strings.Index(text, searchStart)
		if a >= 0 {
			b := strings.Index(text, searchEnd)
			if b < a {
				return fmt.Errorf("Invalid XML declaration")
			}

			if _, err := w.Write([]byte(text[:a])); err != nil {
				return err
			}
			a = b + len(searchEnd)
		}

	}
	// insert new XML declaration
	if insertXMLDecl {
		if _, err := w.Write([]byte(xmlDecl)); err != nil {
			return err
		}
	}
	// append rest of the text
	_, err := w.Write([]byte(text[a+1:]))
	return err
}

func appendFile(path string, w *bufio.Writer) error {
	// open input file
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	buf, err := ioutil.ReadAll(fi)
	if err != nil {
		return err
	}
	text := string(buf)

	return appendString(text, w, false, true)

}

func appendFiles(tag, pattern string, w *bufio.Writer) error {

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	w.WriteString(fmt.Sprintf("<%s>\n", tag))

	for _, srcPath := range matches {
		if err := appendFile(srcPath, w); err != nil {
			return err
		}

	}

	w.WriteString(fmt.Sprintf("</%s>\n", tag))

	return nil

}

// ****************************************************************
// [Unmarshal an ISO-8859-1 XML input in Go](http://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go)

/*
func charsetNewReader(cs string, input io.Reader) (io.Reader, error) {
	return charset.NewReader(input, cs)
}

func readFileISO88591(path string) (string, error) {
	// http://stackoverflow.com/questions/24555819/golang-persist-using-iso-8859-1-charset
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	r, err := charset.NewReader(file, "latin1")
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func convertFile(srcPath, dstPath string) error {
	const searchStart = "<?xml"
	const searchEnd = "?>"

	text, err := readFileISO88591(srcPath)
	if err != nil {
		return err
	}

	fo, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer fo.Close()

	// make a write buffer
	w := bufio.NewWriter(fo)

	a := strings.Index(text, searchStart)
	if a >= 0 {
		b := strings.Index(text, searchEnd)
		if b < a {
			return fmt.Errorf("Invalid XML declaration")
		}

		if _, err := w.Write([]byte(text[:a])); err != nil {
			return err
		}
		a = b + len(searchEnd)
	}

	if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>`)); err != nil {
		return err
	}

	// write last chunk or all text, id idx==-1
	if _, err := w.Write([]byte(text[a+1:])); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func convertFiles(pattern, destFolder string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for j, srcPath := range matches {
		destPath := filepath.Join(destFolder, filepath.Base(srcPath))
		fmt.Printf("%02d) '%s' -> '%s'\n", j+1, srcPath, destPath)
		convertFile(srcPath, destPath)
	}
	fmt.Println("Done.")
	return nil
}
*/
