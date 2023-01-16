package moneytracker

import "errors"

var ErrNotFound error = errors.New("not found")
var ErrUnauthorized error = errors.New("unauthorized")
