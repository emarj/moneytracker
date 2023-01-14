package datetime

import (
	"time"

	dt "github.com/emarj/moneytracker/datetime"
)

const Epsilon = time.Millisecond

var BeginningOfTime dt.DateTime
var BEFORE dt.DateTime
var Before dt.DateTime
var Now dt.DateTime
var Later dt.DateTime
var LATER dt.DateTime
var EndOfTime dt.DateTime

var Times []dt.DateTime
var Past []dt.DateTime
var Future []dt.DateTime

func Init() {
	timeNow := time.Now()

	BeginningOfTime = dt.FromTime(timeNow.AddDate(-1000, 0, 0))
	BEFORE = dt.FromTime(timeNow.AddDate(0, -1, 0))
	Before = dt.FromTime(timeNow.AddDate(0, 0, -1))

	Now = dt.FromTime(timeNow)

	Later = dt.FromTime(timeNow.AddDate(0, 0, 1))
	LATER = dt.FromTime(timeNow.AddDate(0, 1, 0))
	EndOfTime = dt.FromTime(timeNow.AddDate(1000, 0, 0))

	Past = []dt.DateTime{
		BeginningOfTime, BEFORE, Before,
	}

	Future = []dt.DateTime{
		Later, LATER, EndOfTime,
	}

	Times = append(Past, Now)
	Times = append(Times, Future...)
}
