// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"go/format"

	"github.com/naoina/toml"
)

const (
	defaultConfigFile      = "templates.config.toml"
	defaultPageEnumType    = "PageEnum"
	defaultOutputFile      = "templates.go"
	defaultPackageName     = "templates"
	defaultDefaultPageBase = "base"

	errLoadPrefix  = "Error loading configuration file: "
	errParsePrefix = "Error parsing configuration file: "

	mTemplates      = "mTemplates"
	mTemplateFolder = "mTemplatesFolder"
)

type tomlConfig struct {
	FuncMap         string
	Folder          string
	OutputFile      string
	PackageName     string
	DefaultPageBase string
	PageEnumType    string
	PageEnumPrefix  string
	PageEnumSuffix  string
	Templates       map[string][]string
	Pages           map[string]struct {
		Template string
		Base     string
	}
}

type extConfig struct {
	pages               *OrderedSetString
	templates           *OrderedSetString
	pageIdx2templateIdx []int
}

// package variables
var cfg *tomlConfig

var ext *extConfig

func initConfig() *tomlConfig {
	quit := func(err error) {
		panic(fmt.Errorf("%s%s", errLoadPrefix, err.Error()))
	}

	// command line arguments
	pConfigFile := flag.String("config", defaultConfigFile, "templates configuration file")
	pOutFile := flag.String("output", defaultOutputFile, "output file")
	pPkgName := flag.String("package", defaultPackageName, "package name")
	flag.Parse()

	// open config file
	f, err := os.Open(*pConfigFile)
	if err != nil {
		quit(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		quit(err)
	}

	var config tomlConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		quit(err)
	}

	// defaults

	if *pOutFile != defaultOutputFile || config.OutputFile == "" {
		config.OutputFile = *pOutFile
	}
	if *pPkgName != defaultPackageName || config.PackageName == "" {
		config.PackageName = *pPkgName
	}
	if config.PageEnumType == "" {
		config.PageEnumType = defaultPageEnumType
	}
	if config.DefaultPageBase == "" {
		config.DefaultPageBase = defaultDefaultPageBase
	}

	fmt.Printf("config.FuncMap = %s\n", config.FuncMap)

	return &config
}

func initExt() *extConfig {

	// pages
	pages := NewOrderedSetString()
	for pagename := range cfg.Pages {
		pages.Add(pagename)
	}
	pages.Sort()

	// templates
	templates := NewOrderedSetString()
	for _, pagename := range pages.items {
		templates.Add(cfg.Pages[pagename].Template)
	}
	templates.Sort()

	// page-index -> template-idx
	p2t := make([]int, pages.Len())
	for pageIdx, pageName := range pages.ToSlice() {
		templateIdx, ok := templates.Index(cfg.Pages[pageName].Template)
		if !ok {
			panic("initExt error: template not found")
		}
		p2t[pageIdx] = templateIdx
	}

	extcfg := extConfig{
		pages:               pages,
		templates:           templates,
		pageIdx2templateIdx: p2t,
	}

	return &extcfg
}

// ----------------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------------

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

func astr2str(items []string) string {
	b := new(bytes.Buffer)

	for j, s := range items {
		if j > 0 {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprintf(b, "%q", s)
	}
	return b.String()
}

func aint2str(items []int) string {
	b := new(bytes.Buffer)

	for j, item := range items {
		if j > 0 {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprintf(b, "%d", item)
	}
	return b.String()
}

// -------------------------------------------------------------------------
// OrderedSetString
// -------------------------------------------------------------------------

// OrderedSetString is an ordered set of string
type OrderedSetString struct {
	items      []string
	item2index map[string]int
}

// NewOrderedSetString return a new empty OrderedSetString
func NewOrderedSetString() *OrderedSetString {
	set := OrderedSetString{
		item2index: make(map[string]int),
	}
	return &set
}

// Add adds the element. Returns the index of the element
func (set *OrderedSetString) Add(s string) int {
	idx, ok := set.item2index[s]
	if !ok {
		idx = len(set.items)
		set.item2index[s] = idx
		set.items = append(set.items, s)
	}
	return idx
}

// AddSlice adds all element of the slice.
func (set *OrderedSetString) AddSlice(as []string) {
	for _, s := range as {
		set.Add(s)
	}
}

// ToSlice returns the members of the set as a slice.
func (set *OrderedSetString) ToSlice() []string {
	return set.items
}

// Contains returns whether the given item is in the set.
func (set *OrderedSetString) Contains(s string) bool {
	_, ok := set.item2index[s]
	return ok
}

// Index returns the position of the given element in the set
func (set *OrderedSetString) Index(s string) (int, bool) {
	j, ok := set.item2index[s]
	return j, ok
}

// Len returns the number of elements of the set
func (set *OrderedSetString) Len() int {
	return len(set.items)
}

// Sort order the array of string
func (set *OrderedSetString) Sort() {
	sort.Strings(set.items)
	for idx, item := range set.items {
		set.item2index[item] = idx
	}
}

// -------------------------------------------------------------------------
// resolveIncludes
// -------------------------------------------------------------------------

func resolveIncludes() map[string][]string {

	type set map[string]struct{}

	m := make(map[string][]string)

	var resolve func(string, set)

	resolve = func(templateName string, visited set) {

		if _, ok := m[templateName]; ok {
			// already resolved
			return
		}

		if _, ok := visited[templateName]; ok {
			panic(fmt.Errorf("found invalid cyclic template (%s)", templateName))
		}

		// add name to the set of already included templates
		visited[templateName] = struct{}{}

		// iter over each template files
		var files []string

		for _, item := range cfg.Templates[templateName] {
			// check if it's an include item
			if _, ok := cfg.Templates[item]; ok {
				// it's an include
				resolve(item, visited)
				files = append(files, m[item]...)
			} else {
				// append the file
				files = append(files, item)
			}

		}

		m[templateName] = files
	}

	res := make(map[string][]string)

	for _, tmplName := range ext.templates.ToSlice() {
		resolve(tmplName, set{})
		res[tmplName] = m[tmplName]
	}

	return res
}

// -------------------------------------------------------------------------

func printHeader(w io.Writer) {
	fmt.Fprintf(w, `
// Generated by %v; DO NOT EDIT
// Creation date: %v

package %s

import (
  "html/template"
  "io"
  "path/filepath"
)

`, strings.Join(os.Args, " "), time.Now(), cfg.PackageName)
}

// printPageEnum stampa la definizione del tipo PageEnum e delle relative costanti
func printPageEnum(w io.Writer) {

	fmt.Fprintf(w, "// %s type definition\n", cfg.PageEnumType)
	fmt.Fprintf(w, "type %s uint%d\n\n", cfg.PageEnumType, usize(ext.pages.Len()))

	fmt.Fprintf(w, "// %s constants\n", cfg.PageEnumType)
	fmt.Fprintf(w, "const (\n")

	first := true
	s := fmt.Sprintf(" %s = iota", cfg.PageEnumType)
	for _, pagename := range ext.pages.ToSlice() {

		fmt.Fprintf(w, "  %s%s%s%s\n", cfg.PageEnumPrefix, pagename, cfg.PageEnumSuffix, s)
		if first {
			s = ""
			first = false
		}

	}
	fmt.Fprintf(w, ")\n\n")
}

func printVars(w io.Writer) {
	fmt.Fprintf(w, "const %s = %q\n", mTemplateFolder, cfg.Folder)
	fmt.Fprintf(w, "var %s [%d]*template.Template\n\n", mTemplates, ext.templates.Len())
}

func printGetTemplate(w io.Writer) {
	const varname = "idx"

	fmt.Fprintf(w, "// Template returns the template.Template of the page\n")
	fmt.Fprintf(w, "func (page %s) Template() *template.Template {\n", cfg.PageEnumType)
	fmt.Fprintf(w, "  var %s = [...]uint%d{%s}\n", varname, usize(ext.templates.Len()), aint2str(ext.pageIdx2templateIdx))
	fmt.Fprintf(w, "  return %s[%s[page]]\n", mTemplates, varname)
	fmt.Fprint(w, "}\n\n")
}

func printGetBase(w io.Writer) {
	const (
		vArrPageIdx2BaseIdx = "pi2bi"
		vArrBases           = "bases"
	)

	names := NewOrderedSetString()
	p2n := make([]int, ext.pages.Len())

	for j, pagename := range ext.pages.ToSlice() {
		basename := cfg.Pages[pagename].Base
		p2n[j] = names.Add(basename)
	}

	fmt.Fprintf(w, "// Base returns the template name of the page\n")
	fmt.Fprintf(w, "func (page %s) Base() string {\n", cfg.PageEnumType)
	fmt.Fprintf(w, "  var %s = [...]string{%s}\n", vArrBases, astr2str(names.ToSlice()))
	if names.Len() == len(p2n) {
		// each page has a different base
		fmt.Fprintf(w, "  return %s[page]\n", vArrBases)
	} else {
		// some pages have the same base -> remap needed
		fmt.Fprintf(w, "  var %s = [...]uint%d{%s}\n", vArrPageIdx2BaseIdx, usize(ext.pages.Len()), aint2str(p2n))
		fmt.Fprintf(w, "  return %s[%s[page]]\n", vArrBases, vArrPageIdx2BaseIdx)
	}
	fmt.Fprint(w, "}\n\n")
}

func printHelpers(w io.Writer) {

	fmt.Fprintf(w, `
func files2paths(files []string) []string {
	var path string
	paths := make([]string, len(files))
	for i, file := range files {
		switch {
		case len(file) == 0, file[0] == '.', file[0] == filepath.Separator:
			path = file
		default:
			path = filepath.Join(%s, file)
		}
		paths[i] = path
	}
	return paths
}

`, mTemplateFolder)

}

func printExecute(w io.Writer) {
	fmt.Fprintf(w, `
// Execute applies a parsed page template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output, execution
// stops, but partial results may already have been written to the output writer.
// A template may be executed safely in parallel.
func (page %s) Execute(wr io.Writer, data interface{}) error {
	tmpl := page.Template()
	name := page.Base()
	if name == "" {
		return tmpl.Execute(wr, data)
	}
	return tmpl.ExecuteTemplate(wr, name, data)
}

`, cfg.PageEnumType)
}

func printMain(w io.Writer) {
	fmt.Fprintf(w, "/*\n")
	fmt.Fprintf(w, "func main() {\n")
	fmt.Fprintf(w, "  var page %s\n", cfg.PageEnumType)
	fmt.Fprintf(w, "  page = %s%s%s\n", cfg.PageEnumPrefix, ext.pages.ToSlice()[0], cfg.PageEnumSuffix)
	fmt.Fprintf(w, `
  wr := os.Stdout

  if err := page.Execute(wr, nil); err != nil {
	fmt.Print(err)
  }

}
*/
`)
}

/*
func printInitOld(w io.Writer) {
	const sFiles = "templatesFiles"

	tmpl := resolveIncludes()

	fmt.Fprintf(w, "func init() {\n")

	// s_files declaration
	fmt.Fprintf(w, "  var %s = [...][]string{\n", sFiles)
	for _, name := range ext.templates.ToSlice() {
		fmt.Fprintf(w, "  {%s}, // %s\n", astr2str(tmpl[name]), name)
	}
	fmt.Fprintf(w, "  }\n\n")

	fmt.Fprintf(w, `
// init base templates
for i, files := range %s {
	%s[i] = template.Must(template.ParseFiles(files2paths(files)...))
}


`, sFiles, mTemplates)

	fmt.Fprintf(w, "}\n\n")
}
*/

func printInit(w io.Writer) {
	const sFiles = "templatesFiles"
	const sIdxs = "templatesIdxs"

	tmpl := resolveIncludes()
	files := NewOrderedSetString()
	t2f := make([][]int, ext.templates.Len())

	for tmplIdx, tmplName := range ext.templates.ToSlice() {
		a := make([]int, len(tmpl[tmplName]))
		for j, fileName := range tmpl[tmplName] {
			a[j] = files.Add(fileName)
		}
		t2f[tmplIdx] = a
	}

	fmt.Fprintf(w, "func init() {\n")

	// s_files declaration
	fmt.Fprintf(w, "  var %s = [...]string{%s}\n", sFiles, astr2str(files.ToSlice()))

	fmt.Fprintf(w, "  var %s = [...][]uint%d{\n", sIdxs, usize(files.Len()))
	for tmplIdx, fileIdxs := range t2f {
		fmt.Fprintf(w, "  {%s}, // %s\n", aint2str(fileIdxs), ext.templates.ToSlice()[tmplIdx])
	}
	fmt.Fprintf(w, "  }\n\n")

	var templMust string
	if len(cfg.FuncMap) > 0 {
		// with funcMap
		templMust = fmt.Sprintf("template.Must(template.New(filepath.Base(files[0])).Funcs(%s).ParseFiles(files2paths(files)...))", cfg.FuncMap)
	} else {
		// without funcMap
		templMust = "template.Must(template.ParseFiles(files2paths(files)...))"
	}

	fmt.Fprintf(w, `
// init base templates
for i, idxs := range %[1]s {
	files := make([]string, len(idxs))
	for j, idx := range idxs {
		files[j] = %[2]s[idx]
	}
	%[3]s[i] = %[4]s
}

`, sIdxs, sFiles, mTemplates, templMust)

	/*
	   	for _, name := range ext.templates.ToSlice() {
	   		fmt.Fprintf(w, "  {%s}, // %s\n", astr2str(tmpl[name]), name)
	   	}
	   	fmt.Fprintf(w, "  }\n\n")
	   	fmt.Fprintf(w, `
	   // init base templates
	   for i, files := range %s {
	   	%s[i] = template.Must(template.ParseFiles(files2paths(files)...))
	   }


	   `, sFiles, mTemplates)

	*/
	fmt.Fprintf(w, "}\n\n")
}

func printAll(w io.Writer) {
	printHeader(w)
	printPageEnum(w)
	printVars(w)
	printInit(w)
	printHelpers(w)
	printGetTemplate(w)
	printGetBase(w)
	printExecute(w)
	printMain(w)
}

func doJob() {

	cfg = initConfig()
	ext = initExt()

	//w := os.Stdout
	var buf bytes.Buffer // A Buffer needs no initialization.
	w := &buf

	printAll(w)

	// write buffer to output file
	file, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		panic(fmt.Errorf("Error writing output file: %s", err.Error()))
	}
	defer file.Close()

	//buf.WriteTo(file)
	byt, err := format.Source(buf.Bytes())
	file.Write(byt)

}

func main() {
	doJob()
}
