package workers

import "errors"

var ErrOutOfCapacity = errors.New("out of worker capacity")
var ErrInvalidWeblink = errors.New("invalid weblink")
