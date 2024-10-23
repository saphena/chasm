package main

import "time"

const StoredDateWithTZ = "2006-01-02T15:04:05-07:00"
const StoredDateNoTZ = "2006-01-02T15:04"

func parseStoredDate(ts string) time.Time {

	var res time.Time
	var err error
	res, err = time.Parse(StoredDateWithTZ, ts)
	if err != nil {
		res, _ = time.ParseInLocation(StoredDateNoTZ, ts, RallyTimezone)
	}
	return res
}
