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
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
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

		fmt.Fprint(w, buf.String())
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8080", r)
}

func toHtml(raw []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
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
