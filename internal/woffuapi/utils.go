package woffuapi

import "time"

// The function `getTimezoneOffsetInMinutes` calculates the time zone offset in minutes between the local
// time and UTC.
func getTimezoneOffsetInMinutes() int {
	_, localOffset := time.Now().Zone()
	_, utcOffset := time.Now().UTC().Zone()
	return (utcOffset - localOffset) / 60
}
