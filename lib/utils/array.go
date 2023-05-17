package utils

import "strconv"

func IntArrToStrArr(data []int) []string {
	result := make([]string, len(data))

	for _, i := range data {
		result = append(result, strconv.Itoa(i))
	}

	return result
}