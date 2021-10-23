//go:build addindex
// +build addindex

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	dir = flag.String("dir", "content/docs/nbp/", "Scan dir path")
)

var frontMatter = template.Must(template.New("frontmatter").Parse(tpl))

func main() {
	flag.Parse()

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Panicln(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		year, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		name := filepath.Join(*dir, file.Name(), "_index.md")

		log.Println("Create:", name)

		f, err := os.Create(name)
		if err != nil {
			log.Panicln(err)
		}
		defer f.Close()

		collapsable := "true"
		if year == 2021 {
			collapsable = "false"
		}
		values := map[string]string{
			"Title":       file.Name(),
			"Weight":      fmt.Sprintf("%d", 2077-year),
			"Collapsable": collapsable,
		}
		if err := frontMatter.Execute(f, values); err != nil {
			log.Panicln(err)
		}

		if err := f.Close(); err != nil {
			log.Panicln(err)
		}
	}
}

const tpl = `---
bookCollapseSection: {{ .Collapsable }}
weight:  {{ .Weight }}
title: {{ .Title }}
---
`
