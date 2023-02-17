package cachetime

import "time"

const (

	Permanent = time.Duration(0)

	Second_1  = time.Duration(1 * time.Second)
	Second_2  = time.Duration(2 * time.Second)
	Second_3  = time.Duration(3 * time.Second)
	Second_5  = time.Duration(5 * time.Second)
	Second_10 = time.Duration(10 * time.Second)
	Second_20 = time.Duration(20 * time.Second)
	Second_30 = time.Duration(30 * time.Second)

	Minute_1  = time.Duration(1 * time.Minute)
	Minute_2  = time.Duration(2 * time.Minute)
	Minute_3  = time.Duration(3 * time.Minute)
	Minute_5  = time.Duration(5 * time.Minute)
	Minute_10 = time.Duration(10 * time.Minute)
	Minute_20 = time.Duration(20 * time.Minute)
	Minute_30 = time.Duration(30 * time.Minute)

	Hour_1  = time.Duration(1 * time.Hour)
	Hour_2  = time.Duration(2 * time.Hour)
	Hour_6  = time.Duration(6 * time.Hour)
	Hour_12 = time.Duration(12 * time.Hour)

	Day_1 = time.Duration(24 * time.Hour)
	Day_3 = Day_1 * 3
	Day_7 = Day_1 * 7

	Month_1 = Day_1 * 30
	Month_3 = Day_1 * 90

	Year_1 = Day_1 * 365
)
