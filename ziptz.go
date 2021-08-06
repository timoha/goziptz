package goziptz

import (
	"bufio"
	"bytes"
	"compress/zlib"
	_ "embed"
	"errors"
	"regexp"
	"time"
)

//go:embed data/tz.data
var tz []byte

var (
	// ErrInvalidZip is returned when zip code is not in zip or zip+4 formats.
	ErrInvalidZip = errors.New("invalid zip code")

	// ErrLocationNotFound is returned when location for passed zip code is not found.
	ErrLocationNotFound = errors.New("error location not found")
)

var zipRegexp = regexp.MustCompile(`^\d{5}(?:[-\s]\d{4})?$`)
var sep = []byte("|")

// ZipToLocation returns IANA timezone of the US zip code.
// The zip code can either be 5 digit format ("94606") or 5+4 digit format ("94606-1234").
func ZipToLocation(zip string) (*time.Location, error) {
	if !zipRegexp.MatchString(zip) {
		return nil, ErrInvalidZip
	}

	zb := []byte(zip[:5])

	r, err := zlib.NewReader(bytes.NewBuffer(tz))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		if bytes.HasPrefix(b, zb) {
			split := bytes.SplitN(b, sep, 3)
			return time.LoadLocation(string(split[1]))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, ErrLocationNotFound
}
