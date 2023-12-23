package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func getWordsByCnt(words map[string]int, cnt int) []string {
	result := make([]string, 0, 2)
	for k, v := range words {
		if v == cnt {
			result = append(result, k)
			delete(words, k)
		}
	}

	return result
}

func sortAndLimit(w []string, lim int) []string {
	sort.Strings(w)
	if lim > len(w) {
		lim = len(w)
	}
	return w[:lim]
}

func Top10(s string) []string {
	// Place your code here.

	if s == "" {
		return nil
	}

	words := regexp.MustCompile(`(\S*[^[:punct:]\s])|[-]{2,}`).FindAllString(s, -1)
	cntWords := make(map[string]int)
	for _, w := range words {
		cntWords[strings.ToLower(w)]++
	}
	variants := make([]int, 0, len(cntWords))

	fmt.Println(cntWords)
	for _, val := range cntWords {
		variants = append(variants, val)
	}

	sort.Ints(variants)

	result := make([]string, 0, 10)

	for cnt := 0; cnt < 10; {
		fmt.Println(cnt)
		words := getWordsByCnt(cntWords, variants[len(variants)-cnt-1])

		if len(words) > 1 {
			result = append(result, sortAndLimit(words, 10-cnt)...)
		} else {
			result = append(result, words...)
		}
		cnt = len(result)
	}

	return result
}
