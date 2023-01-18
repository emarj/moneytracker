package sqlite

import (
	"os"
	"testing"

	tt "github.com/emarj/moneytracker/timestamp/testtimes"
)

func TestMain(m *testing.M) {
	tt.Init()
	os.Exit(m.Run())
}
