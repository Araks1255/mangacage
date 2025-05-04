package utils

func ConvertToMap(slice []string) map[int]string {
	result := make(map[int]string, len(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = slice[i]
	}

	return result
}
