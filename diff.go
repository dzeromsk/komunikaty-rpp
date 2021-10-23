//go:build diff
// +build diff

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	input  = flag.String("input", "pdf/", "Input dir path")
	output = flag.String("output", "content/docs/nbp/", "Output dir path")
)

var frontMatter = template.Must(template.New("frontmatter").Parse(tpl))

func main() {
	flag.Parse()

	matches, err := filepath.Glob(filepath.Join(*input, "*.md"))
	if err != nil {
		log.Panicln(err)
	}

	for i := range matches[1:] {
		a, b := matches[i], matches[i+1]

		var year, month, day int
		n, err := fmt.Sscanf(filepath.Base(b), "rpp_%d_%02d_%02d.md", &year, &month, &day)
		if err != nil || n != 3 {
			log.Panicln(err, n, b)
		}

		name := filepath.Clean(fmt.Sprintf("%s/%d/%d-%02d-%02d.html", *output, year, year, month, day))

		if _, err := os.Stat(name); err == nil {
			println("skip", name)
			continue
		}

		println("diff:", a, b, name)

		var data [2][]byte
		for j := 0; j < 2; j++ {
			data[j], err = ioutil.ReadFile(matches[i+j])
			if err != nil {
				log.Panicln(err)
			}
		}
		patch := diffmatchpatch.New()

		var diffs []diffmatchpatch.Diff
		diffs = patch.DiffMain(string(data[0]), string(data[1]), false)
		diffs = patch.DiffCleanupEfficiency(diffs)
		diffs = patch.DiffCleanupSemantic(diffs)
		diffs = patch.DiffCleanupMerge(diffs)

		err = os.MkdirAll(filepath.Dir(name), 0755)
		if err != nil {
			log.Panicln(err)
		}

		f, err := os.Create(name)
		if err != nil {
			log.Panicln(err)
		}
		defer f.Close()

		values := map[string]string{
			"Title":  fmt.Sprintf("%d-%02d-%02d", year, month, day),
			"URL":    fmt.Sprintf("komunikaty/%d/%02d/%02d", year, month, day),
			"Weight": fmt.Sprintf("%d", 13-month),
			"Header": fmt.Sprintf("%d %s %d r.", day, m[month], year),
		}
		if err := frontMatter.Execute(f, values); err != nil {
			log.Panicln(err)
		}

		html := patch.DiffPrettyHtml(diffs)
		if _, err := f.WriteString(html); err != nil {
			log.Panicln(err)
		}

		if err := f.Close(); err != nil {
			log.Panicln(err)
		}
	}
}

const tpl = `---
title: {{ .Title }}
weight: {{ .Weight }}
url: {{ .URL }}
---
<h1>{{ .Header }}</h1>
`

var m = []string{
	"",
	"stycznia",
	"lutego",
	"marca",
	"kwietnia",
	"maja",
	"czerwca",
	"lipca",
	"sierpnia",
	"września",
	"października",
	"listopada",
	"grudnia",
}
