package view

import (
	"math"
	"strconv"
	"time"
)

var (
	suffixes [5]string
)

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func convert(sizeB float64) string {
	var getSuffix string
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"
	base := math.Log(sizeB) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	if sizeB == 0 {
		getSuffix = suffixes[0]
		getSize = 0
	} else {
		getSuffix = suffixes[int(math.Floor(base))]
	}
	resp := strconv.FormatFloat(getSize, 'f', -1, 64) + string(getSuffix)
	return resp
}

func timeElapsed(inputDate int64, full bool) string {
	var precise [8]string
	var text string

	formattedDate := time.Unix(inputDate, 0)

	now := time.Now()
	year2, month2, day2 := now.Date()
	hour2, minute2, second2 := now.Clock()

	year1, month1, day1 := formattedDate.Date()
	hour1, minute1, second1 := formattedDate.Clock()

	year := math.Abs(float64(int(year2 - year1)))
	month := math.Abs(float64(int(month2 - month1)))
	day := math.Abs(float64(int(day2 - day1)))
	hour := math.Abs(float64(int(hour2 - hour1)))
	minute := math.Abs(float64(int(minute2 - minute1)))
	second := math.Abs(float64(int(second2 - second1)))

	week := math.Floor(day / 7)

	if year > 0 {
		if int(year) == 1 {
			precise[0] = strconv.Itoa(int(year)) + " year"
		} else {
			precise[0] = strconv.Itoa(int(year)) + " years"
		}
	}

	if month > 0 {
		if int(month) == 1 {
			precise[1] = strconv.Itoa(int(month)) + " month"
		} else {
			precise[1] = strconv.Itoa(int(month)) + " months"
		}
	}

	if week > 0 {
		if int(week) == 1 {
			precise[2] = strconv.Itoa(int(week)) + " week"
		} else {
			precise[2] = strconv.Itoa(int(week)) + " weeks"
		}
	}

	if day > 0 {
		if int(day) == 1 {
			precise[3] = strconv.Itoa(int(day)) + " day"
		} else {
			precise[3] = strconv.Itoa(int(day)) + " days"
		}
	}

	if hour > 0 {
		if int(hour) == 1 {
			precise[4] = strconv.Itoa(int(hour)) + " hour"
		} else {
			precise[4] = strconv.Itoa(int(hour)) + " hours"
		}
	}

	if minute > 0 {
		if int(minute) == 1 {
			precise[5] = strconv.Itoa(int(minute)) + " minute"
		} else {
			precise[5] = strconv.Itoa(int(minute)) + " minutes"
		}
	}

	if second > 0 {
		if int(second) == 1 {
			precise[6] = strconv.Itoa(int(second)) + " second"
		} else {
			precise[6] = strconv.Itoa(int(second)) + " seconds"
		}
	}

	for _, v := range precise {
		if v != "" {
			if v[len(v)-5:len(v)-1] != "cond" {
				precise[7] += v + ", "
			} else {
				precise[7] += v
			}
		}
	}

	if full {
		return precise[7] + text
	} else {
		for k, v := range precise {
			if v != "" {
				return precise[k] + " ago"
			}
		}
	}
	return "invalid date"
}
