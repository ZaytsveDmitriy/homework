package hw09structvalidator

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrStringTemplateInvalid = errors.New("validation template invalid")
	ErrStringNotContain      = errors.New("value not containt in limit slice")
	ErrStringLenNotEqual     = errors.New("value len not equal with len lim")
	ErrStringNotMatchRE      = errors.New("value not match RE")
)

var (
	lenRe      = regexp.MustCompile(`(?:\||^)len:(\d+)(\||$)`)
	stringInRE = regexp.MustCompile(`(?:\||^)in:((?:[^[:cntrl:]|][,]?)+[^|])(?:\||$)`)
	reRE       = regexp.MustCompile(`(?:\||^)regexp:([^|]+)(?:\||$)`)
)

type StringFilter struct {
	re     *regexp.Regexp
	lenLim int64
	inLim  []string
}

func NewStringFilter(pattern string) (*StringFilter, error) {
	var err error

	filter := StringFilter{
		nil,
		-1,
		nil,
	}

	match := lenRe.FindStringSubmatch(pattern)
	if len(match) > 0 {
		filter.lenLim, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, ErrTemplateInvalid
		}
	}

	match = reRE.FindStringSubmatch(pattern)
	if len(match) > 0 {
		filter.re, err = regexp.Compile(match[1])
		if err != nil {
			return nil, ErrTemplateInvalid
		}
	}

	match = stringInRE.FindStringSubmatch(pattern)
	if len(match) > 0 {
		filter.inLim = strings.Split(match[1], ",")
		if err != nil {
			return nil, ErrTemplateInvalid
		}
	}

	return &filter, nil
}

func (filt StringFilter) Validate(args ...string) error {
	for _, arg := range args {
		if filt.lenLim >= 0 && utf8.RuneCountInString(arg) != int(filt.lenLim) {
			return ErrStringLenNotEqual
		}
		if filt.inLim != nil && !slices.Contains(filt.inLim, arg) {
			return ErrStringNotContain
		}
		if filt.re != nil {
			res := filt.re.FindStringSubmatch(arg)
			if len(res) == 0 {
				return ErrStringNotMatchRE
			}
		}
	}
	return nil
}
