package datetime

import (
	"time"

	dt "github.com/emarj/moneytracker/timestamp"
)

const Epsilon = time.Millisecond

var BeginningOfTime dt.Timestamp
var BEFORE dt.Timestamp
var Before dt.Timestamp
var Now dt.Timestamp
var Later dt.Timestamp
var LATER dt.Timestamp
var EndOfTime dt.Timestamp

var Times []dt.Timestamp
var Past []dt.Timestamp
var Future []dt.Timestamp

func Init() {
	timeNow := time.Now()

	BeginningOfTime = dt.FromTime(timeNow.AddDate(-1000, 0, 0))
	BEFORE = dt.FromTime(timeNow.AddDate(0, -1, 0))
	Before = dt.FromTime(timeNow.AddDate(0, 0, -1))

	Now = dt.FromTime(timeNow)

	Later = dt.FromTime(timeNow.AddDate(0, 0, 1))
	LATER = dt.FromTime(timeNow.AddDate(0, 1, 0))
	EndOfTime = dt.FromTime(timeNow.AddDate(1000, 0, 0))

	Past = []dt.Timestamp{
		BeginningOfTime, BEFORE, Before,
	}

	Future = []dt.Timestamp{
		Later, LATER, EndOfTime,
	}

	Times = append(Past, Now)
	Times = append(Times, Future...)
}
