package ptrConversion

func BoolToInt64Ptr(boolPtr *bool) *int64 {
	if boolPtr != nil {
		var result int64

		if *boolPtr {
			result = int64(1)
		} else {
			result = int64(0)
		}

		return &result
	}

	return nil
}

func Int64ToBoolPtr(int64ptr *int64) *bool {
	if int64ptr != nil {
		var result bool

		if *int64ptr == 0 {
			result = false
		} else {
			result = true
		}

		return &result
	}

	return nil
}
