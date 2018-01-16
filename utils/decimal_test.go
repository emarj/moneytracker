package utils

import "testing"

func TestParseDecimal(t *testing.T) {

	var n int
	n, err := ParseDecimal("123")
	if err != nil {
		t.Error(err)
	}
	expecting := 12300
	if n != expecting {
		t.Errorf("got %d instead of %d", n, expecting)
	}

	n, err = ParseDecimal("123.23")
	if err != nil {
		t.Error(err)
	}
	expecting = 12323
	if n != expecting {
		t.Errorf("got %d instead of %d", n, expecting)
	}

	n, err = ParseDecimal("123.9")
	if err != nil {
		t.Error(err)
	}
	expecting = 12390
	if n != expecting {
		t.Errorf("got %d instead of %d", n, expecting)
	}

	n, err = ParseDecimal("123,89")
	if err != nil {
		t.Error(err)
	}
	expecting = 12389
	if n != expecting {
		t.Errorf("got %d instead of %d", n, expecting)
	}

	n, err = ParseDecimal("123.890")
	if err == nil {
		t.Error("expecting error")
	}

	n, err = ParseDecimal("123.56.9")
	if err == nil {
		t.Error("expecting error")
	}

}

func TestFormatDecimal(t *testing.T) {
	str := FormatDecimal(12389)
	if str != "123.89" {
		t.Errorf("got %s instead", str)
	}

	str = FormatDecimal(0)
	if str != "0.00" {
		t.Errorf("got %s instead", str)
	}

	str = FormatDecimal(9)
	if str != "0.09" {
		t.Errorf("got %s instead", str)
	}

	str = FormatDecimal(50)
	if str != "0.50" {
		t.Errorf("got %s instead", str)
	}

}
