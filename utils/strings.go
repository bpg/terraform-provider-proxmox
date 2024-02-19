package utils

// ConvertToStringSlice helps convert interface slice to string slice.
func ConvertToStringSlice(interfaceSlice []interface{}) []string {
	resultSlice := []string{}
	for _, val := range interfaceSlice {
		resultSlice = append(resultSlice, val.(string))
	}

	return resultSlice
}
