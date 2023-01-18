package timestamp

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlusMinus(t *testing.T) {
	d := Now()

	d2 := d.Minus(time.Second)

	assert.False(t, d.Time.Equal(d2.Time), fmt.Sprintf("%q and %q must not be equal", d, d2))

}

func TestFormattingTrailingZeros(t *testing.T) {

	N := 1_000_000
	times := make([]string, N)
	for i := 0; i < N; i++ {
		times[i] = Now().String()
	}

	length := len(times[0])
	for i := 0; i < N; i++ {
		if len(times[i]) != length {
			t.Fatalf("[%d]: %s != %s", i, times[i], times[0])
		}
	}

}
