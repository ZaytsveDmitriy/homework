package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(s string) []string {
	// Place your code here.

	if s == "" {
		return nil
	}

	words := regexp.MustCompile(`(\S*[^[:punct:]\s])|[-]{2,}`).FindAllString(s, -1)

	wordsQnt := make(map[string]int)

	for _, w := range words {
		wordsQnt[strings.ToLower(w)]++
	}

	words = words[:0]

	for k := range wordsQnt {
		words = append(words, k)
	}

	sort.Strings(words)
	sort.SliceStable(words, func(i, j int) bool { return wordsQnt[words[i]] > wordsQnt[words[j]] })

	return words[:min(len(words), 10)]
}
