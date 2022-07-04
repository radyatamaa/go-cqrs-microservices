package helper

// InlineConditionPointerFloat64 returns val1 if condition, otherwise val2
func InlineConditionPointerFloat64(condition bool, val1, val2 *float64) *float64 {
	if condition {
		return val1
	}
	return val2
}

func InlineConditionFloatAndPointerToFloat(condition bool, val1 float64, val2 *float64) float64 {
	if condition {
		return val1
	}
	return *val2
}

// InlineConditionString returns val1 if condition, otherwise val2
func InlineConditionString(condition bool, val1, val2 string) string {
	if condition {
		return val1
	}
	return val2
}

func InlineConditionPointerInt32(condition bool, val1, val2 *int32) *int32 {
	if condition {
		return val1
	}
	return val2
}