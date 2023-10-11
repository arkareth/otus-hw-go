package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

var wr = regexp.MustCompile(`[\p{L}\d]+`)

var pm = ",.!?-:;'\" /\\()[]{}<>"

func Top10(s string) []string {
	if len(s) == 0 {
		return nil
	}

	f := strings.Fields(s)
	m := make(map[string]int)
	for _, i := range f {
		if wr.MatchString(i) {
			c := strings.ToLower(strings.Trim(i, pm))
			m[c]++
			continue
		}
		if len(i) > 1 {
			m[i]++
			continue
		}
	}

	wc := make([]wordCount, 0, len(m))
	for k, v := range m {
		wc = append(wc, wordCount{k, v})
	}

	sort.Slice(wc, func(i, j int) bool {
		if wc[i].count == wc[j].count {
			return wc[i].word < wc[j].word
		}
		return wc[i].count > wc[j].count
	})

	l := 10
	if len(wc) < l {
		l = len(wc)
	}

	ret := make([]string, 0, l)
	for _, v := range wc[:l] {
		ret = append(ret, v.word)
	}

	return ret
}
