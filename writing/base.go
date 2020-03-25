package writing

import (
	"strconv"
	"time"
)

func GetBaseTs() string {
	var ts = strconv.FormatInt(time.Now().Unix(), 10)
	for {
		if len(ts) < 16 {
			ts = "0" + ts
		} else {
			break
		}
	}
	return ts
}
