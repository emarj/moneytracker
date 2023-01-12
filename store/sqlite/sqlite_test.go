package sqlite

import (
	"os"
	"testing"
	"time"

	"github.com/emarj/moneytracker/datetime"
)

func TestMain(m *testing.M) {
	initTimes()
	os.Exit(m.Run())
}

var timeNow time.Time

var BeginningOfTime datetime.DateTime
var BEFORE datetime.DateTime
var Before datetime.DateTime
var before datetime.DateTime
var now datetime.DateTime
var later datetime.DateTime
var Later datetime.DateTime
var LATER datetime.DateTime
var EndOfTime datetime.DateTime

var times []datetime.DateTime
var past []datetime.DateTime
var future []datetime.DateTime

func initTimes() {
	timeNow = time.Now()

	BeginningOfTime = datetime.FromTime(timeNow.AddDate(-1000, 0, 0))
	BEFORE = datetime.FromTime(timeNow.AddDate(0, -1, 0))
	Before = datetime.FromTime(timeNow.AddDate(0, 0, -1))
	before = datetime.FromTime(timeNow.Add(-time.Minute))

	now = datetime.FromTime(timeNow)

	later = datetime.FromTime(timeNow.Add(time.Minute))
	Later = datetime.FromTime(timeNow.AddDate(0, 0, 1))
	LATER = datetime.FromTime(timeNow.AddDate(0, 1, 0))
	EndOfTime = datetime.FromTime(timeNow.AddDate(1000, 0, 0))

	past = []datetime.DateTime{
		BeginningOfTime, BEFORE, Before, before,
	}

	future = []datetime.DateTime{
		later, Later, LATER, EndOfTime,
	}

	times = append(past, now)
	times = append(times, future...)
}
