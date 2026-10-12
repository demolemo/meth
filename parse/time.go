// Package parse turns an exported file into metric values. It knows about file
// formats and nothing about storage - the sqlite package does the inserting.
package parse

import "time"

// weight and steps don't need exact seconds
const TimeLayout = "2006-01-02 15:04"

func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation(TimeLayout, s, time.Local)
}
