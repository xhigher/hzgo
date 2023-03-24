package utils

import (
	"fmt"
	"time"
)

func FormatYmd(ts int64) string {
	return time.Unix(ts, 0).Format(TimeFormatYMD)
}

func FormatYmdhms(ts int64) string {
	return time.Unix(ts, 0).Format(TimeFormatYMDHMS)
}

func ParseYmd(ymd string) time.Time {
	t, _ := time.ParseInLocation(TimeFormatYMD, ymd, time.Local)
	return t
}

func ParseYmdhms(ymdhms string) time.Time {
	t, _ := time.ParseInLocation(TimeFormatYMDHMS, ymdhms, time.Local)
	return t
}

func FormatMoney(money int64) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

func FmtMaskPhoneno(phoneno string) string {
	num := len(phoneno)
	if num <= 6 {
		return phoneno
	}
	if num > 11 {
		phoneno = phoneno[num-11:]
	}
	start := phoneno[:3]
	end := phoneno[11-3:]
	return fmt.Sprintf("%s*****%s", start, end)
}

func FmtMaskIdCardCode(code string) string {
	codeLen := len(code)
	if codeLen == 18 || codeLen == 15 {
		start := code[:6]
		end := code[codeLen-3:]
		mask := "*********"
		if codeLen == 15 {
			mask = "******"
		}
		return fmt.Sprintf("%s%s%s", start, mask, end)
	}
	return code
}

func FmtMaskIdCardName(name string) string {
	nameLen := len(name)
	if nameLen < 2 {
		return name
	} else if nameLen == 2 {
		return fmt.Sprintf("%s*", name[:1])
	} else if nameLen == 3 {
		return fmt.Sprintf("%s*%s", name[:1], name[nameLen-1:])
	}
	return fmt.Sprintf("%s**%s", name[:1], name[nameLen-1:])
}
