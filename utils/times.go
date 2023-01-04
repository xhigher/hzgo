package utils

import (
	"fmt"
	"strings"
	"time"

)

const (
	TimeNone          = "0000-00-00"
	TimeFormatYM      = "2006-01"
	TimeFormatYM2     = "200601"
	TimeFormatYMD     = "2006-01-02"
	TimeFormatYMD2    = "20060102"
	TimeFormatYMDHM   = "2006-01-02 15:04"
	TimeFormatYMDHMS  = "2006-01-02 15:04:05"
	TimeFormatYMDHMS2 = "20060102150405"
	TimeFormatHM      = "15:04"
	TimeFormatHMS     = "15:04:05"
)

func NowTimeNano() int64 {
	return time.Now().UnixNano()
}

func NowTime() int64 {
	return time.Now().Unix()
}

func NowTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func NowYmd() string {
	return time.Now().Format(TimeFormatYMD)
}

func YesterdayYmd() string {
	return time.Now().Add(time.Hour * -24).Format(TimeFormatYMD)
}

func NowYmdhms() string {
	return time.Now().Format(TimeFormatYMDHMS)
}

func NowYmdhms2() string {
	return time.Now().Format(TimeFormatYMDHMS2)
}

func NowWeekDate(now time.Time) string {
	if now.IsZero() {
		now = time.Now()
	}
	year, week := now.ISOWeek()
	return GetWeekDate(year, week)
}

func GetWeekDate(year, week int) string {
	format := "%d-00-%d"
	if week < 10 {
		format = "%d-00-0%d"
	}
	return fmt.Sprintf(format, year, week)
}

func ParseDayTime(hms string) (int64, error) {
	layout := TimeFormatHMS
	if len(hms) == 5 {
		layout = TimeFormatHM
	}
	p, err := time.Parse(layout, hms)
	if err != nil {
		return 0, err
	}
	return int64(p.Hour())*3600 + int64(p.Minute())*60 + int64(p.Second()), nil
}

func GetTodayTimestamps() (now, start, end int64) {
	nt := time.Now()
	now = nt.Unix()
	st := time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, time.Local)
	start = st.Unix()
	nt = nt.Add(time.Hour * 24)
	et := time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, time.Local)
	end = et.Unix()
	return
}

func ParseDateToTimestamp(ymd string) (ts int64, err error) {
	layout := TimeFormatYMD
	if strings.Count(ymd, "-") == 0 {
		layout = TimeFormatYMD2
	}
	date, err := time.ParseInLocation(layout, ymd, time.Local)
	if err != nil {
		return
	}
	ts = date.Unix()
	return
}

func ParseDateTimeToTimestamp(ymdhms string) (ts int64, err error) {
	layout := TimeFormatYMDHMS
	if strings.Count(ymdhms, "-") == 0 {
		layout = TimeFormatYMDHMS2
	}
	date, err := time.ParseInLocation(layout, ymdhms, time.Local)
	if err != nil {
		return
	}
	ts = date.Unix()
	return
}

func GetDatesFromNow(days int) (dates []string) {
	now := time.Now()
	dates = make([]string, days)
	for i := 0; i < days; i++ {
		dates[i] = now.AddDate(0, 0, -(i + 1)).Format(TimeFormatYMD)
	}
	return
}

func GetStartTimeOfWeek() time.Time {
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	nowDate = nowDate.AddDate(0, 0, -int(now.Weekday()))
	return nowDate
}
