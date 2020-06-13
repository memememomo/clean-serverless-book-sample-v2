package utils

import "strconv"

func ParseUint(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 64)
}
