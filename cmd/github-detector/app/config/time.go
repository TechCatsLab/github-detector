/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	hour  = int64(time.Hour)
	day   = 24 * hour
	mouth = 30 * day
	year  = 365 * day
)

// StringToTime converts string to time.
func StringToTime(t, sep string) (time.Duration, error) {
	ts := strings.Split(t, sep)

	var duration int64
	for index := range ts {
		if yi := strings.Index(ts[index], "y"); yi != -1 {
			yn, err := strconv.Atoi(ts[index][:yi])
			if err != nil {
				return 0, err
			}
			duration += int64(yn) * year
		}
		if mi := strings.Index(ts[index], "m"); mi != -1 {
			mn, err := strconv.Atoi(ts[index][:mi])
			if err != nil {
				return 0, err
			}
			duration += int64(mn) * mouth
		}
		if di := strings.Index(ts[index], "d"); di != -1 {
			mn, err := strconv.Atoi(ts[index][:di])
			if err != nil {
				return 0, err
			}
			duration += int64(mn) * day
		}
		if hi := strings.Index(ts[index], "h"); hi != -1 {
			hn, err := strconv.Atoi(ts[index][:hi])
			if err != nil {
				return 0, err
			}
			duration += int64(hn) * hour
		}
	}

	return time.Duration(duration), nil
}

// TimeToString converts time to string.
func TimeToString(t time.Duration, sep string) (string, error) {
	ns := int64(t)

	if ns < 0 {
		return "", fmt.Errorf("%v is invalid time duration", t)
	}

	y := ns / year
	ns = ns % year

	m := ns / mouth
	ns = ns % mouth

	d := ns / day
	ns = ns % day

	h := ns / hour

	var str string
	if y != 0 {
		str = fmt.Sprintf("%dy", y)
	} else if m != 0 {
		str = fmt.Sprintf("%s%s%dm", str, sep, m)
	} else if d != 0 {
		str = fmt.Sprintf("%s%s%dd", str, sep, d)
	} else if h != 0 {
		str = fmt.Sprintf("%s%s%dh", str, sep, h)
	}

	if str == "" {
		return "", fmt.Errorf("%v is too short to be parsed", t)
	}

	return str, nil
}
