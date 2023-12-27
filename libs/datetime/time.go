package datetime

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Defined constant custom date/time format.
// See https://golang.org/src/time/format.go for additional custom format.

const (
	// DateDMYSlash predefined layouts DD/MM/YY
	DateDMYSlash = "02/01/06"
	// DateMDYSlash predefined layouts MM/Dd/YY
	DateMDYSlash = "01/02/06"
	// DateMDYTimeSlash predefined layouts MM/Dd/YY hh:mm
	DateMDYTimeSlash = "01/02/06 15:04"
	// DateYYYMMDDDash predefined layouts YYYY-MM-DD
	DateYYYMMDDDash = "2006-01-02"
	// DateUTC predefined layouts YYYY-MM-DDThh:mm:ss
	DateUTC = "2006-01-02T15:04:05"
)

// TimeSlice is alias a slice time to implement the sort.Interface
type TimeSlice []time.Time

// Sort sorts a slice of time in increasing order.
func Sort(times []time.Time) { sort.Sort(TimeSlice(times)) }

func (d TimeSlice) Len() int { return len(d) }

func (d TimeSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d TimeSlice) Less(i, j int) bool {
	return d[i].Before(d[j])
}

// NewTime parses a formatted string and returns the time value it represents.
func NewTime(layout, value string) (time.Time, error) {
	var (
		dateVal   = strings.TrimSpace(value)
		layoutVal = strings.TrimSpace(layout)
	)
	if len(dateVal) == 0 || len(layoutVal) == 0 {
		return time.Time{}, fmt.Errorf("requires value [%v] and layout [%v]", value, layout)
	}
	t, err := time.Parse(layoutVal, dateVal)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// NewMilliSecTime parses the given milliseconds and returns the time value.
func NewMilliSecTime(timestamps int64) time.Time {
	return time.Unix(0, timestamps*int64(time.Millisecond))
}

// TimeValue returns a time
func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
