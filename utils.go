package moneytracker

import "strconv"

func Atoi64(s string) (int64, error) {
	var res int64
	res, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return res, err
	}
	return res, err
}
