package hw09structvalidator

import (
	"errors"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrTemplateInvalid = errors.New("validation template invalid")
	ErrIntLess         = errors.New("value less then min limit")
	ErrIntBiger        = errors.New("value biger then max limit")
	ErrIntNotContain   = errors.New("value not containt in limit slice")
)

var (
	minRE = regexp.MustCompile(`(?:\||^)min:(\d+)(\||$)`)
	maxRE = regexp.MustCompile(`(?:\||^)max:(\d+)(\||$)`)
	inRE  = regexp.MustCompile(`(?:\||^)in:((\d+[,]?)+)(\||$)`)
)

type IntFilter struct {
	minLim, maxLim int64
	inLim          []int64
}

func NewIntFilter(pattern string) (*IntFilter, error) {
	var (
		err          error
		hasValidator bool
	)

	filter := IntFilter{
		minLim: math.MinInt,
		maxLim: math.MaxInt,
		inLim:  nil,
	}

	minMatch := minRE.FindStringSubmatch(pattern)
	if len(minMatch) > 0 {
		filter.minLim, err = strconv.ParseInt(minMatch[1], 10, 64)
		if err != nil {
			return nil, ErrTemplateInvalid
		}
		hasValidator := true
	}

	maxMatch := maxRE.FindStringSubmatch(pattern)
	if len(maxMatch) > 0 {
		filter.maxLim, err = strconv.ParseInt(minMatch[1], 10, 64)
		if err != nil {
			return nil, ErrTemplateInvalid
		}
		hasValidator := true
	}

	inMatch := inRE.FindStringSubmatch(pattern)
	if len(inMatch) > 0 {
		vals := strings.Split(inMatch[1], ",")
		for _, val := range vals {
			digit, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, ErrTemplateInvalid
			}
			filter.inLim = append(filter.inLim, digit)
		}
		hasValidator := true
	}

	if !hasValidator {
		return nil, ErrTemplateInvalid
	}

	return &filter, nil
}

func (filt IntFilter) Validate(args ...int64) error {
	for _, arg := range args {
		if arg < filt.minLim {
			return ErrIntLess
		}
		if arg > filt.maxLim {
			return ErrIntBiger
		}
		if filt.inLim != nil && !slices.Contains(filt.inLim, arg) {
			return ErrIntNotContain
		}
	}
	return nil
}
