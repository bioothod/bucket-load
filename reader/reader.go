package reader

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Reader struct {
	Re	*regexp.Regexp
}


type Entry struct {
	Prefix		string
	Date		time.Time
	Method		string
	Path		string
	URI		string
	Status		int
	Size		uint64
	Time		float64
	Error		string
}

func NewReader() (re *Reader, err error) {
	re = &Reader {
	}

	re.Re, err = regexp.Compile("([[:graph:]]+): (\\d+/\\d+/\\d+ \\d+:\\d+:\\d+)\\.(\\d+) access_log: method: '(\\w+)', path: '([[:graph:]]+)', encoded-uri: '([[:graph:]]+)', status: (\\d+), size: (\\d+), time: ([0-9\\.]+) ms, err: '(\\w+)'")
	return
}

func (re *Reader) Run(str string) (*Entry, error) {
	res := re.Re.FindAllStringSubmatch(str, -1)
	if res == nil {
		return nil, fmt.Errorf("nothing matched: string: '%s'", str)
	}

	r := res[0]

	status, _ := strconv.Atoi(r[7])
	size, _ := strconv.ParseUint(r[8], 0, 64)
	ts, _ := strconv.ParseFloat(r[9], 64)

	format := "2006/01/02 15:04:05"
	tm, _ := time.Parse(format, r[2])
	usecs, _ := strconv.ParseUint(r[3], 0, 64)
	tm.Add(time.Duration(usecs) * time.Microsecond)

	e := &Entry {
		Prefix:		r[1],
		Date:		tm,
		Method:		r[4],
		Path:		r[5],
		URI:		r[6],
		Status:		status,
		Size:		size,
		Time:		ts,
		Error:		r[10],
	}

	return e, nil
}
