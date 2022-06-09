package utils

import "fmt"

func GetRelTimeFromSecond(sec uint16) string {
	now := int(sec)
	hour := 0
	minute := 0
	second := 0
	if now >= 60*60 {
		hour = now / 3600
		now = now % 3600
	}
	if now >= 60 {
		minute = now / 60
		now = now % 60
	}
	if now >= 1 {
		second = now
	}
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}
