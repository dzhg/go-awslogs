package tsp

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"time"
)

var unitMap = map[string]int64{
	"s": int64(time.Second),
	"m": int64(time.Minute),
	"h": int64(time.Hour),
	"d": int64(time.Hour * 24),
	"w": int64(time.Hour * 24 * 7),
}

var agoPatten = regexp.MustCompile(`(\d+\.?\d*)\s?(s|second|seconds|m|min|minute|minutes|h|hr|hour|hours|d|day|days|w|wk|week|weeks)(?: ago)?`)

// ParseString parses the string as millisecond
// it also supports relative time presentation
func ParseString(s string) (int64, error) {

	// shortcut for relative time of now
	t, err := ParseRelative(s, time.Now())
	if err == nil {
		return t.UnixNano() / int64(time.Millisecond), err
	}

	// let's try RFC3339
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return t.UnixNano() / int64(time.Millisecond), err
	}

	// let's try another format
	t, err = time.Parse("01/02/2006 15:04:05 MST", s)

	if err == nil {
		return t.UnixNano() / int64(time.Millisecond), err
	}

	return 0, errors.Wrap(err, "parse string to time")
}

// ParseRelative parses the string and returns the time relative to the input time
func ParseRelative(s string, t time.Time) (time.Time, error) {
	if agoPatten.MatchString(s) {
		matches := agoPatten.FindStringSubmatch(s)
		n := matches[1]
		u := string(matches[2][0])
		unit := unitMap[u]
		i, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return time.Time{}, err
		}
		var d = time.Duration(float64(unit) * i)
		return t.Add(-d), nil
	}

	return time.Time{}, fmt.Errorf("str \"%s\" doesn't match relative time pattern", s)
}
