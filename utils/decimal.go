package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseDecimal(numstr string) (int, error) {
	numstr = strings.Replace(numstr, ",", ".", -1)
	num := strings.Split(numstr, ".")
	if len(num) > 2 || (len(num) == 2 && len(num[1]) > 2) {
		return 0, fmt.Errorf("invalid Amount format\"%s\", has to be \"xxx.xx\"", numstr)
	}
	euro := num[0]
	cent := "00"
	if len(num) != 1 {
		cent = num[1]
		if len(cent) == 1 {
			cent = cent + "0"
		}

	}
	intgr, err := strconv.Atoi(euro + cent)
	if err != nil {
		return 0, nil
	}

	return intgr, nil
}

func FormatDecimal(n int) string {
	var str string

	sign := ""
	if n < 0 {
		sign = "-"
		n = -1 * n
	}

	if n < 100 {
		str = "0."
		if n < 10 {
			str = str + "0"
		}
		return sign + str + strconv.Itoa(n)
	}

	str = strconv.Itoa(n)
	return sign + str[:len(str)-2] + "." + str[len(str)-2:]
}
