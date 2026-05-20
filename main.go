package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"go.abhg.dev/goldmark/wikilink"

	"html/template"
)

type HalamanData struct {
	Files []string
}

func main() {
	textTemplate := `
<!DOCTYPE html>
<html>
<head>
	<title>Uji Coba</title>
</head>
<body>

<ul>
	{{ range .Files }}
		<li>
			<a href="{{ . }}">{{ . }}</a>
		</li>
	{{ end }}
</ul>

</body>
</html>
`
	r := http.NewServeMux()

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		files, err := getEntries("./content")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("web").Parse(textTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := HalamanData{
			Files: files,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	r.HandleFunc("GET /content/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		path := filepath.Join("content", id)

		data, err := os.ReadFile(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		buf, err := toHtml(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, buf.String())
	})

	http.ListenAndServe(":8080", r)
}

func getEntries(root string) ([]string, error) {
	var entries []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root {
				return filepath.SkipDir
			}
		}

		if !d.IsDir() {
			entries = append(entries, path)
		}

		return nil
	})

	return entries, err
}

func toHtml(raw []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			&wikilink.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	err := md.Convert(raw, &buf)

	return buf, err
}
