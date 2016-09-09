package timeparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var unitsre string = "(seconds|second|secs|sec|minutes|minute|min|hours|hour|hrs|hr|days|day|weeks|week|w[0-6]|months|month|mon|quarters|quarter|qtrs|qtr|years|year|yrs|yr|s|h|m|d|w|y|w|q)"
var reltimere string = "(?i)(?P<plusminus>[+-]*)(?P<num>\\d{1,})(?P<unit>" + unitsre + "{1})(([\\@](?P<snapunit>" + unitsre + "{1})((?P<snapplusminus>[+-])(?P<snaprelnum>\\d+)(?P<snaprelunit>" + unitsre + "{1}))*)*)"
var re *regexp.Regexp = regexp.MustCompile(reltimere)

func TimeParser(ts string, now func() time.Time) (time.Time, error) {
	if now == nil {
		now = time.Now
	}

	if ts == "now" {
		return now(), nil
	} else {
		if ts[:1] == "+" || ts[:1] == "-" {
			ret := now()

			match := re.FindStringSubmatch(ts)
			results := make(map[string]string)
			for i, name := range re.SubexpNames() {
				if i != 0 {
					results[name] = match[i]
				}
			}

			// Handle first part of the time string
			if len(results["plusminus"]) != 0 && len(results["num"]) != 0 && len(results["unit"]) != 0 {
				timeParserTimeMath(results["plusminus"], results["num"], results["unit"], &ret)

				return ret, nil
			}
		}
	}
	return now(), fmt.Errorf("Got to the end but didn't return")
}

func timeParserTimeMath(plusminus string, numstr string, unit string, ret *time.Time) {
	num, _ := strconv.Atoi(numstr)
	if plusminus == "-" {
		num *= -1
	}

	secs := map[string]bool{"s": true, "sec": true, "secs": true, "second": true, "seconds": true}
	mins := map[string]bool{"m": true, "min": true, "minute": true, "minutes": true}
	hours := map[string]bool{"h": true, "hr": true, "hrs": true, "hour": true, "hours": true}
	days := map[string]bool{"d": true, "day": true, "days": true}
	weeks := map[string]bool{"w": true, "week": true, "weeks": true}
	months := map[string]bool{"mon": true, "month": true, "months": true}
	quarters := map[string]bool{"q": true, "qtr": true, "qtrs": true, "quarter": true, "quarters": true}
	years := map[string]bool{"y": true, "yr": true, "yrs": true, "year": true, "years": true}

	switch {
	case secs[unit]:
		*ret = ret.Add(time.Duration(num) * time.Second)
	case mins[unit]:
		*ret = ret.Add(time.Duration(num) * time.Minute)
	case hours[unit]:
		*ret = ret.Add(time.Duration(num) * time.Hour)
	case days[unit]:
		*ret = ret.AddDate(0, 0, num)
	case weeks[unit]:
		*ret = ret.AddDate(0, 0, num*7)
	case months[unit]:
		*ret = ret.AddDate(0, num, 0)
	case quarters[unit]:
		*ret = ret.AddDate(0, num*3, 0)
	case years[unit]:
		*ret = ret.AddDate(num, 0, 0)
	}
}
