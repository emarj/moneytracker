package datetime

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
