package utils

import "strconv"

func ParseUintSlice(strings []string) ([]uint, error) {
	res := make([]uint, len(strings), len(strings))

	for i := 0; i < len(strings); i++ {
		elem, err := strconv.ParseUint(strings[i], 10, 64)
		if err != nil {
			return nil, err
		}

		res[i] = uint(elem)
	}

	return res, nil
}
