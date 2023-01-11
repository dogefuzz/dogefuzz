package common

import "time"

func Now() time.Time {
	loc, _ := time.LoadLocation("")
	now := time.Now().In(loc)
	return now
}
