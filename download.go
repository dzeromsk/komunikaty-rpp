//go:build download
// +build download

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
)

var (
	dir = flag.String("dir", "pdf/", "Input dir path")
	all = flag.Bool("all", false, "Try downloading all files")
)

func main() {
	flag.Parse()

	index := indexdir(*dir)

	now := time.Now()
	y, m, d := now.Year(), now.Month(), now.Day()

	for i := y; i > 2002; i-- {
		limit := time.Month(12)
		if i == y {
			limit = m
		}
		for j := limit; j > 0; j-- {
			if _, ok := index[i<<4|int(j)]; ok {
				// skip if we already have this month
				continue
			}
			limit := 31
			if j == m && i == y {
				limit = d
			}
			for k := limit; k > 0; k-- {
				url := fmt.Sprintf("https://www.nbp.pl/polityka_pieniezna/dokumenty/files/rpp_%d_%02d_%02d.pdf", i, j, k)

				log.Println(url)

				resp, err := grab.Get(*dir, url)
				if err != nil {
					continue
				}

				log.Println(resp.Filename)

				if !*all {
					return // fast exit
				}

				break // try next month
			}
		}
	}
}

func indexdir(dir string) map[int]struct{} {
	index := map[int]struct{}{
		// vacations
		2021<<4 | 8: {},
		2020<<4 | 8: {},
		2019<<4 | 8: {},
		2018<<4 | 8: {},
		2017<<4 | 8: {},
		2016<<4 | 8: {},
		2015<<4 | 8: {},
		2014<<4 | 8: {},
		2013<<4 | 8: {},
		2012<<4 | 8: {},
		2011<<4 | 8: {},
		2011<<4 | 2: {},
		2010<<4 | 7: {},
	}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		var year, month, day int
		n, err := fmt.Sscanf(filepath.Base(path), "rpp_%d_%02d_%02d.pdf", &year, &month, &day)
		if err != nil || n != 3 {
			return nil
		}
		index[year<<4|month] = struct{}{}
		return nil
	})
	return index
}
